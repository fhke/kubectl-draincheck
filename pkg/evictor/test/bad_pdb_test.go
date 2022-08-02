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

// Check pod with a pod disruption budget that does not allow evictions. This should be an error.
func TestBlockingPDB(t *testing.T) {
	t.Parallel()

	var (
		podLabels = map[string]string{
			"application": "nginx",
			"test-case":   "blocking-pdb",
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

	// create a PDB affecting the pod that does not allow disruptions
	_, err = clientset.PolicyV1().PodDisruptionBudgets("default").Create(
		context.TODO(),
		pdbFactory.NewBasicPodDisruptionBudget("nginx-blocking-pdb", namespace, 1, podLabels),
		metav1.CreateOptions{},
	)
	require.NoError(t, err, "Creating pod disruption budget should not return an error")

	// wait for pod to become healthy
	ctx, can := context.WithTimeout(context.Background(), time.Minute)
	defer can()
	err = wait.PodReady(ctx, clientset, pod.Name, pod.Namespace)
	require.NoError(t, err, "Pod should become healthy within a minute")

	// try to evict the pod
	err = evictor.NewEvictor(clientset).DryRun(context.TODO(), *pod)
	require.Error(t, err, "Eviction attempt should raise error")
	assert.Equal(t, evictor.ErrNoDisruptions, err, "Eviction error should be unevictable")
}
