package locator

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	policyv1 "k8s.io/api/policy/v1"
)

func (p *PDBLocator) PDBsForPod(ctx context.Context, pod *corev1.Pod) ([]*policyv1.PodDisruptionBudget, error) {
	// locate by pod
	return p.pdbLister.GetPodPodDisruptionBudgets(pod)
}

func (p *PDBLocator) Stop() {
	p.infStop <- struct{}{}
}
