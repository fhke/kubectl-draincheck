package factory

import (
	"github.com/fhke/kube-test-utils/pkg/kuberesources/factory"
	policyv1 "k8s.io/api/policy/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func NewBasicPodDisruptionBudget(name, namespace string, minAvailable int, selectorLabels map[string]string) *policyv1.PodDisruptionBudget {
	minU := intstr.FromInt(minAvailable)

	return &policyv1.PodDisruptionBudget{
		ObjectMeta: factory.NewObjectMeta(name, namespace, nil, nil),
		Spec: policyv1.PodDisruptionBudgetSpec{
			MinAvailable: &minU,
			Selector: &metav1.LabelSelector{
				MatchLabels: selectorLabels,
			},
		},
	}
}
