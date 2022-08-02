package evictor

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	policyv1 "k8s.io/api/policy/v1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (e *evictorImpl) DryRun(ctx context.Context, pod corev1.Pod) error {
	// create an eviction object
	err := e.k.PolicyV1().Evictions(pod.Namespace).Evict(
		ctx,
		&policyv1.Eviction{
			// ObjectMeta represents the pod that we want to check eviction for
			ObjectMeta: pod.ObjectMeta,
			DeleteOptions: &metav1.DeleteOptions{
				// DryRun All tests the pod eligibility for eviction rather than actually evicting
				DryRun: []string{
					"All",
				},
				PropagationPolicy: deletionPropagationPtr(metav1.DeletePropagationForeground),
			},
		},
	)

	return errorFor(err)
}

func deletionPropagationPtr(in metav1.DeletionPropagation) *metav1.DeletionPropagation {
	return &in
}
