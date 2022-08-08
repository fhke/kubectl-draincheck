package locator

import (
	"context"

	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
)

func NewPDBLocator(ctx context.Context, k kubernetes.Interface) (*PDBLocator, error) {
	// create shared informer factory & pdb lister
	infFactory := informers.NewSharedInformerFactory(k, 0)
	pdbLister := infFactory.Policy().V1().PodDisruptionBudgets().Lister()

	// start informers
	stopCh := make(chan struct{}, 1)
	infFactory.Start(stopCh)

	// create locator
	pdbl := &PDBLocator{
		infFactory: infFactory,
		infStop:    stopCh,
		pdbLister:  pdbLister,
	}

	// wait for cache sync, then return locator
	return pdbl, pdbl.waitForCacheSync(ctx)
}
