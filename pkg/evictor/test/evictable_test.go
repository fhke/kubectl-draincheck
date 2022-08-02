package evictor_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/fhke/kubectl-draincheck/pkg/evictor"
	pdbFactory "github.com/fhke/kubectl-draincheck/pkg/testutils/factory"
	"github.com/fhke/kube-test-utils/pkg/kuberesources/factory"
	"github.com/fhke/kube-test-utils/pkg/wait"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Check pod with one pod disruption budget that allows evictions
func TestEvictablePod(t *testing.T) {
	t.Parallel()

	var (
		podLabels = map[string]string{
			"application": "nginx",
			"test-case":   "single-working-pdb",
		}
		podName   = fmt.Sprintf("test-%d", time.Now().UnixNano())
		namespace = "default"
	)

	// create a basic nginx pod
	pod, err := clientset.CoreV1().Pods("default").Create(
		context.TODO(),
		factory.NewBasicPod(podName, namespace, "nginx:mainline", podLabels),
		metav1.CreateOptions{},
	)
	require.NoError(t, err, "Creating pod should not return an error")

	// create one valid PDB affecting the pod
	_, err = clientset.PolicyV1().PodDisruptionBudgets("default").Create(
		context.TODO(),
		pdbFactory.NewBasicPodDisruptionBudget("nginx-single-pdb", namespace, 0, podLabels),
		metav1.CreateOptions{},
	)
	require.NoError(t, err, "Creating pod disruption budget should not return an error")

	// wait for pod to become healthy
	ctx, can := context.WithTimeout(context.Background(), time.Minute)
	defer can()
	err = wait.PodReady(ctx, clientset, pod.Name, pod.Namespace)
	require.NoError(t, err, "Pod should become healthy within a minute")

	// try to evict the pod
	err = evictor.NewEvictor(clientset).DryRun(ctx, *pod)
	assert.NoError(t, err, "Eviction attempt should not create an error")
}

// Check pod with no pod disruption budgets
func TestEvictablePod2(t *testing.T) {
	var (
		podLabels = map[string]string{
			"application": "nginx",
			"test-case":   "no-pdb",
		}
		podName   = fmt.Sprintf("test-%d", time.Now().UnixNano())
		namespace = "default"
	)

	// create a basic nginx pod
	pod, err := clientset.CoreV1().Pods("default").Create(
		context.TODO(),
		factory.NewBasicPod(podName, namespace, "nginx:mainline", podLabels),
		metav1.CreateOptions{},
	)
	require.NoError(t, err, "Creating pod should not return an error")

	// wait for pod to become healthy
	ctx, can := context.WithTimeout(context.Background(), time.Minute)
	defer can()
	err = wait.PodReady(ctx, clientset, pod.Name, pod.Namespace)
	require.NoError(t, err, "Pod should become healthy within a minute")

	// try to evict the pod
	err = evictor.NewEvictor(clientset).DryRun(ctx, *pod)
	assert.NoError(t, err, "Eviction attempt should not create an error")
}
