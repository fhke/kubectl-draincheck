package checker

import (
	"github.com/fhke/kubectl-draincheck/pkg/evictor"
	"github.com/fhke/kubectl-draincheck/pkg/locator"
	corev1 "k8s.io/api/core/v1"
	policyv1 "k8s.io/api/policy/v1"
	"k8s.io/client-go/kubernetes"
)

type (
	Checker struct {
		k          kubernetes.Interface // kubernetes clientset interface
		e          evictor.Evictor      // pod dry run evictor
		pdbLocator *locator.PDBLocator
	}
	Result struct {
		Reason               error      `json:"reason"`
		Pod                  corev1.Pod `json:"pod"`
		PodDisruptionBudgets []*policyv1.PodDisruptionBudget
	}
	Results []Result
)
