package checker

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/fhke/kubectl-draincheck/pkg/evictor"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var ErrNoOwnerRefs = errors.New("pod has no owner references")

// Check eligibility of all pods to be evicted
func (c *Checker) AllPods(ctx context.Context, namespace string, timeout time.Duration, workers uint) (Results, error) {
	// list pods on cluster
	ctx2, can := context.WithTimeout(ctx, timeout)
	defer can()
	podList, err := c.k.CoreV1().Pods(namespace).List(ctx2, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("error listing pods: %w", err)
	}

	return c.checkPods(ctx, timeout, workers, podList.Items...)
}

// Check eligibility of pods by name
func (c *Checker) PodsByName(ctx context.Context, timeout time.Duration, namespace string, workers uint, podNames ...string) (Results, error) {
	var pods []corev1.Pod

	// get pods
	for _, podName := range podNames {
		ctx2, can := context.WithTimeout(ctx, timeout)
		defer can()
		pod, err := c.k.CoreV1().Pods(namespace).Get(ctx2, podName, metav1.GetOptions{})
		if err != nil {
			return nil, fmt.Errorf("error getting pod %s/%s: %w", namespace, podName, err)
		}
		pods = append(pods, *pod)
	}

	return c.checkPods(ctx, timeout, workers, pods...)
}

// Check eligibility of specified pods
func (c *Checker) checkPods(ctx context.Context, timeout time.Duration, workers uint, pods ...corev1.Pod) (Results, error) {
	// create channel for worker goroutines to read pods
	podCh := make(chan corev1.Pod, len(pods))

	// create channels for returned errors & results
	errCh := make(chan error, len(pods))
	resCh := make(chan Result, len(pods))

	// write each pod to the pod channel
	for _, pod := range pods {
		podCh <- pod
	}
	close(podCh)

	// start worker goroutines
	wg := sync.WaitGroup{}
	for i := uint(0); i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			// read pod from podCh
			for pod := range podCh {
				res, err := c.checkPod(ctx, pod, timeout)
				if err != nil && err != evictor.ErrNotFound {
					// Unexpected error that is not a 404.
					// We swallow 404 errors as we do a get/list before calling this function,
					// so the pod was most likely deleted between initial get/list & checking.
					errCh <- fmt.Errorf("error checking eligibility of pod %s/%s for eviction: %w", pod.Namespace, pod.Name, err)
				} else if res != nil {
					// No unexpected errors but pod cannot be evicted, return error
					resCh <- *res
				}
			}
		}()
	}

	// wait for workers to complete
	wg.Wait()
	close(errCh)
	close(resCh)

	// if there were any errors, return the first one
	for err := range errCh {
		return nil, err
	}

	// read results into slice
	var results Results
	for result := range resCh {
		results = append(results, result)
	}

	return results, nil
}

// Check eligibility of a single pod to be evicted
func (c *Checker) checkPod(ctx context.Context, pod corev1.Pod, timeout time.Duration) (*Result, error) {
	// check if pod has owner references
	if len(pod.OwnerReferences) == 0 {
		// get the PDBs affecting pod
		ctx2, can := context.WithTimeout(ctx, timeout)
		defer can()
		pdbs, err := c.pdbLocator.PDBsForPod(ctx2, &pod)
		if err != nil {
			if !strings.HasPrefix(err.Error(), fmt.Sprintf("could not find PodDisruptionBudget for pod %s in namespace %s with labels: ", pod.Name, pod.Namespace)) {
				// Don't return error if there are no PDBs affecting pod. The lister doesn't use error wrapping so we need to use string matching here
				return nil, fmt.Errorf("error locating pod disruption budgets for pod: %w", err)
			}
		}

		return &Result{
			Reason:               ErrNoOwnerRefs,
			Pod:                  pod,
			PodDisruptionBudgets: pdbs,
		}, nil
	}

	// create child context
	ctx2, can := context.WithTimeout(ctx, timeout)
	defer can()

	// Run eviction dry-run
	evictErr := c.e.DryRun(ctx2, pod)

	if evictErr == nil {
		// no error, return nil
		return nil, nil
	} else if !evictor.IsUnevictableError(evictErr) {
		// unexpected error
		return nil, evictErr
	}

	// if we got here, the pod is unevictable

	// get the PDBs affecting pod
	ctx2, can = context.WithTimeout(ctx, timeout)
	defer can()
	pdbs, err := c.pdbLocator.PDBsForPod(ctx2, &pod)
	if err != nil {
		return nil, fmt.Errorf("error locating pod disruption budgets for pod: %w", err)
	}

	return &Result{
		Reason:               evictErr,
		Pod:                  pod,
		PodDisruptionBudgets: pdbs,
	}, nil
}
