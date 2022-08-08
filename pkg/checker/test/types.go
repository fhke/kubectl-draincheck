package test

// expected represents an expected outcome
type expected struct {
	podName   string
	namespace string
	pdbNames  []string
	reason    string
}
