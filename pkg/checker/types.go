package checker

import (
	"github.com/fhke/kubectl-draincheck/pkg/evictor"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
)

type (
	Checker struct {
		k kubernetes.Interface // kubernetes clientset interface
		e evictor.Evictor      // pod dry run evictor
	}
	Result struct {
		Reason error      `json:"reason"`
		Pod    corev1.Pod `json:"pod"`
	}
	Results []Result
)
