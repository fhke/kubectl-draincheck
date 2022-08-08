package test

import (
	"os"
	"testing"

	kindTesting "github.com/fhke/kube-test-utils/pkg/kind/testing"

	"k8s.io/client-go/kubernetes"
)

// shared clientset for all tests to use
var clientset = &kubernetes.Clientset{}

// configure test environment
func TestMain(m *testing.M) {
	os.Exit(
		kindTesting.TestMain(
			m,
			kindTesting.WithClientset(clientset),
			kindTesting.WithDeleteKubeconfig(),
		),
	)
}
