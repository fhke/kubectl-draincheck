package evictor

import (
	"k8s.io/client-go/kubernetes"
)

func NewEvictor(k kubernetes.Interface) Evictor {
	return &evictorImpl{
		k: k,
	}
}
