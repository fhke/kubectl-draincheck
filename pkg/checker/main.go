package checker

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/fhke/kubectl-draincheck/pkg/evictor"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var ErrNoOwnerRefs = errors.New("pod has no owner references")

// Check eligibility of all pods to be evicted
func (c *Checker) AllPods(ctx context.Context, namespace string, timeout time.Duration) (Results, error) {
	// list pods on cluster
	ctx2, can := context.WithTimeout(ctx, timeout)
	defer can()
	podList, err := c.k.CoreV1().Pods(namespace).List(ctx2, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("error listing pods: %w", err)
	}

	// Check all pods, ignoring not found errors as pods may be
	// deleted between initial list and check.
	return c.checkPods(ctx, timeout, podList.Items...)
}

// Check eligibility of pods by name
func (c *Checker) PodsByName(ctx context.Context, timeout time.Duration, namespace string, podNames ...string) (Results, error) {
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

	return c.checkPods(ctx, timeout, pods...)
}

// Check eligibility of specified pods
func (c *Checker) checkPods(ctx context.Context, timeout time.Duration, pods ...corev1.Pod) (Results, error) {
	var results Results

	// check each pod
	for _, pod := range pods {
		res, err := c.checkPod(ctx, pod, timeout)
		if err != nil && err != evictor.ErrNotFound {
			// Unexpected error that is not a 404.
			// We swallow 404 errors as pod did exist prior to calling function,
			// so it was most likely deleted between initial get/list & checking.
			return nil, fmt.Errorf("error checking eligibility of pod %s/%s for eviction: %w", pod.Namespace, pod.Name, err)
		} else if res != nil {
			// No unexpected errors but pod cannot be evicted, so
			// add to return list.
			results = append(results, *res)
		}
	}

	return results, nil
}

// Check eligibility of a single pod to be evicted
func (c *Checker) checkPod(ctx context.Context, pod corev1.Pod, timeout time.Duration) (*Result, error) {
	// check if pod has owner references
	if len(pod.OwnerReferences) == 0 {
		return &Result{
			Reason: ErrNoOwnerRefs,
			Pod:    pod,
		}, nil
	}

	// create child context
	ctx2, can := context.WithTimeout(ctx, timeout)
	defer can()

	// Run dry-run
	err := c.e.DryRun(ctx2, pod)

	if err == nil {
		// no error, return nil
		return nil, nil
	}

	if evictor.IsUnevictableError(err) {
		return &Result{
			Reason: err,
			Pod:    pod,
		}, nil
	} else {
		return nil, err
	}
}
