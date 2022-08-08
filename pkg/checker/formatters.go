package checker

import (
	"bytes"
	"encoding/json"
	"strings"

	"github.com/fhke/kubectl-draincheck/pkg/checker/errors"
	"github.com/olekukonko/tablewriter"
	"k8s.io/apimachinery/pkg/api/meta"
	"sigs.k8s.io/yaml"
)

// Convert results to JSON
func (r Results) JSON() ([]byte, error) {
	return json.MarshalIndent(r.marshalPrepare(), "", "    ")
}

// Convert results to YAML
func (r Results) YAML() ([]byte, error) {
	return yaml.Marshal(r.marshalPrepare())
}

// Convert results to a human-readable table
func (r Results) Table() []byte {
	// Buffer to store table data
	var buf = &bytes.Buffer{}

	// Prepare table
	tbl := tablewriter.NewWriter(buf)
	tbl.SetHeader([]string{"namespace", "pod", "reason", "pod disruption budgets"})
	tbl.SetAutoWrapText(false)

	// Load table with data
	for _, res := range r {
		tbl.Append([]string{
			res.Pod.Namespace,
			res.Pod.Name,
			res.Reason.Error(),
			res.pdbNames(),
		})
	}

	// render table
	tbl.Render()

	return buf.Bytes()
}

// Set the Reason field to an error type that can be marshalled
// to text, and remove managed fields from resources
func (r Results) marshalPrepare() Results {
	out := make(Results, len(r))

	for i := range r {
		// copy struct
		out[i] = r[i]
		// set error type to internal error
		out[i].Reason = errors.For(out[i].Reason)
		// strip managed fields
		ac, err := meta.Accessor(&out[i].Pod)
		if err == nil {
			ac.SetManagedFields(nil)
		}
	}

	return out
}

// Get comma-separated names of pod disruption budgets affecting pod
func (r Result) pdbNames() string {
	names := make([]string, len(r.PodDisruptionBudgets))

	for i, pdb := range r.PodDisruptionBudgets {
		names[i] = pdb.Name
	}

	return strings.Join(names, ", ")
}
