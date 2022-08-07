package locator

import (
	"k8s.io/client-go/informers"
	policyv1Listers "k8s.io/client-go/listers/policy/v1"
)

type PDBLocator struct {
	infFactory informers.SharedInformerFactory
	infStop    chan struct{}
	pdbLister  policyv1Listers.PodDisruptionBudgetLister
}
