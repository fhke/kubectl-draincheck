package evictor

import (
	"errors"

	kerrors "k8s.io/apimachinery/pkg/api/errors"
)

var (
	ErrTooManyPDBs   error = errors.New("multiple pod disruption budgets are acting on the same pod")
	ErrNoDisruptions error = errors.New("pod disruption budget allows no disruptions")
	ErrNotFound      error = errors.New("pod not found")
)

func IsUnevictableError(err error) bool {
	if err == nil {
		return false
	}

	return err == ErrTooManyPDBs || err == ErrNoDisruptions
}

// See https://kubernetes.io/docs/concepts/scheduling-eviction/api-eviction/#how-api-initiated-eviction-works
// If a pod is unevictable, the API reason will either be "too many requests" or "internal server error"
func errorFor(e error) error {
	if e == nil {
		return nil
	}

	switch codeForKError(e) {
	case 404:
		return ErrNotFound
	case 500:
		return ErrTooManyPDBs
	case 429:
		return ErrNoDisruptions
	default:
		return e
	}
}

// get status code for an error from Kubernetes
func codeForKError(err error) int32 {
	if status := kerrors.APIStatus(nil); errors.As(err, &status) {
		return status.Status().Code
	}
	return 0
}
