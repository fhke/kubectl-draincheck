package evictor

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
)

type (
	Evictor interface {
		DryRun(ctx context.Context, pod corev1.Pod) error // test a pod for
	}

	// implementation of the Evictor interface
	evictorImpl struct {
		k kubernetes.Interface
	}
)

var _ Evictor = &evictorImpl{}
