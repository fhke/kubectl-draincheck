package main

import "github.com/fhke/kubectl-draincheck/cmd/draincheck"

func main() {
	if err := draincheck.NewCmd().Execute(); err != nil {
		panic(err)
	}
}
