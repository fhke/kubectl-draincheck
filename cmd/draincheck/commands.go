package draincheck

import (
	"context"
	"fmt"
	"time"

	"github.com/fhke/kubectl-draincheck/pkg/checker"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var log *zap.SugaredLogger = mustNewLogger()

const (
	OutputYAML = "yaml"
	OutputJSON = "json"
	OutputText = "text"
)

func NewCmd() *cobra.Command {
	var (
		// flags
		namespace, kubeconfig, output *string
		allNamespaces                 *bool
		timeout                       *time.Duration
	)

	cmd := &cobra.Command{
		Use:   "kubectl draincheck [POD ...]",
		Short: "Check whether pods can be evicted by kubectl drain",

		// Validate args
		PreRun: func(cmd *cobra.Command, args []string) {
			if len(args) > 0 && *allNamespaces {
				log.Panic("cannot specify --all-namespaces and specific pods")
			}
			if *output != OutputYAML && *output != OutputJSON && *output != OutputText {
				log.Panicf("Unexpected output format %s. Valid values are %s, %s or %s", *output, OutputJSON, OutputYAML, OutputText)
			}
		},

		Run: func(cmd *cobra.Command, pods []string) {
			defer log.Sync()

			// create clientset
			cs, err := newClientset(getKubeconfigPath(*kubeconfig))
			if err != nil {
				log.Panicw("Error getting clientset", "error", err)
			}

			// create eviction checker
			ch := checker.NewChecker(cs)

			// create parent context
			ctx := context.Background()

			var res checker.Results

			if len(pods) > 0 {
				// check by name
				res, err = ch.PodsByName(ctx, *timeout, *namespace, pods...)
			} else {
				// check all in namespace/cluster
				var ns string
				if !*allNamespaces {
					ns = *namespace
				}
				res, err = ch.AllPods(ctx, ns, *timeout)
			}

			if err != nil {
				log.Panicw("Error checking eligibility of pods for eviction", "error", err)
			}

			// Write data in preferred format
			switch *output {
			case OutputText:
				fmt.Print(string(res.Table()))
			case OutputYAML:
				mustMarshalWrite(log, res.YAML)
			case OutputJSON:
				mustMarshalWrite(log, res.JSON)
			default:
				// We should never get here, as invalid options should be
				// picked up in PreRun
				log.Panicw("Internal error: no formatter found", "output", *output)
			}
		},
	}

	kubeconfig = cmd.PersistentFlags().String("kubeconfig", defaultKubeConfig(), "Path to kubeconfig")
	namespace = cmd.Flags().StringP("namespace", "n", "default", "Namespace")
	allNamespaces = cmd.Flags().BoolP("all-namespaces", "A", false, "Check pods in all namespaces")
	timeout = cmd.Flags().DurationP("api-timeout", "T", time.Second*30, "Timeout for calls to Kubernetes API server")
	output = cmd.Flags().StringP("output", "o", OutputText, "Output format - yaml, json or text")

	return cmd
}
