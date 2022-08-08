package test

import (
	"context"
	"testing"
	"time"

	"github.com/fhke/kubectl-draincheck/pkg/checker"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIntegration(t *testing.T) {
	t.Parallel()

	// create checker
	ch, err := checker.NewChecker(context.TODO(), clientset)
	require.NoError(t, err, "Creating checker should not return error")
	require.NotNil(t, ch, "Checker should not be nil")
	defer ch.Stop()

	// create k8s resources
	expectedResources := createResources(t)
	require.GreaterOrEqual(t, len(expectedResources), 2, "There should be at least two pods")

	// wait for pods to start
	waitPods(t, expectedResources, time.Minute/2, 10)

	// run checker on all pods
	res, err := ch.AllPods(context.TODO(), "", time.Second*30, 5)
	require.NoError(t, err, "Running checker on all pods should not return error")

	// validate that the number of results == the number of expected results
	assert.Equal(t, len(expectedResources), len(res), "Number of results should equal expected")

	// validate each expected pod against the results
	for _, ex := range expectedResources {
		validate(t, ex, res)
	}
}

func validate(t *testing.T, ex expected, res checker.Results) {
	// find pod in ex
	for _, re := range res {
		if re.Pod.Name == ex.podName && re.Pod.Namespace == ex.namespace {
			assert.Equalf(t, len(ex.pdbNames), len(re.PodDisruptionBudgets), "Number of PDBs should match expected for pod %s/%s", ex.namespace, ex.podName)
			assert.Equalf(t, ex.reason, re.Reason.Error(), "Reason should be correct for pod %s/%s", ex.namespace, ex.podName)
			for _, pdbName := range ex.pdbNames {
				assert.Truef(
					t,
					hasPDB(ex.namespace, pdbName, re),
					"Pod %s/%s should have pod disruption budget %s in results",
					ex.namespace,
					ex.podName,
					pdbName,
				)
			}
			return
		}
	}

	// if we got here, we didn't find any matching pods
	t.Errorf("Could not find expected pod %s/%s in results", ex.namespace, ex.podName)
}

func hasPDB(namespace, pdbName string, re checker.Result) bool {
	for _, resPDB := range re.PodDisruptionBudgets {
		if resPDB.Name == pdbName && resPDB.Namespace == namespace {
			return true
		}
	}
	return false
}
