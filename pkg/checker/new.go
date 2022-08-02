package checker

import (
	"github.com/fhke/kubectl-draincheck/pkg/evictor"
	"k8s.io/client-go/kubernetes"
)

// instantiate a new Checker{} from a clientset
func NewChecker(clientset kubernetes.Interface) *Checker {
	// create an Evictor
	e := evictor.NewEvictor(clientset)

	return NewCheckerForEvictor(clientset, e)
}

// instantiate a new Checker{} from a clientset & Evictor
func NewCheckerForEvictor(clientset kubernetes.Interface, e evictor.Evictor) *Checker {
	return &Checker{
		k: clientset,
		e: e,
	}
}
