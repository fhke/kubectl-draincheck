package evictor_test

import (
	"context"
	"testing"

	"github.com/fhke/kubectl-draincheck/pkg/evictor"
	"github.com/fhke/kube-test-utils/pkg/kuberesources/factory"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Try to evict a pod that does not exist. This should raise an error.
func TestEvictionOnNonexistentPod(t *testing.T) {
	t.Parallel()

	// try to evict the pod
	err := evictor.NewEvictor(clientset).DryRun(
		context.TODO(),
		*factory.NewBasicPod("invalid-pod", "default", "nginx:mainline", nil),
	)
	require.Error(t, err, "Eviction attempt should raise error")
	assert.Equal(t, evictor.ErrNotFound, err, "Eviction error should not be unevictable")
}
