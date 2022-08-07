package checker

import (
	"context"
	"fmt"

	"github.com/fhke/kubectl-draincheck/pkg/evictor"
	"github.com/fhke/kubectl-draincheck/pkg/locator"
	"k8s.io/client-go/kubernetes"
)

// instantiate a new Checker{} from a clientset
func NewChecker(ctx context.Context, clientset kubernetes.Interface) (*Checker, error) {
	// create an Evictor
	e := evictor.NewEvictor(clientset)

	return NewCheckerForEvictor(ctx, clientset, e)
}

// instantiate a new Checker{} from a clientset & Evictor
func NewCheckerForEvictor(ctx context.Context, clientset kubernetes.Interface, e evictor.Evictor) (*Checker, error) {
	// create pdb locator
	l, err := locator.NewPDBLocator(ctx, clientset)
	if err != nil {
		return nil, fmt.Errorf("error creating pod disruption budget locator: %w", err)
	}

	return &Checker{
		k:          clientset,
		e:          e,
		pdbLocator: l,
	}, nil
}
