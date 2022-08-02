package checker

import (
	"bytes"
	"encoding/json"

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
	tbl.SetHeader([]string{"namespace", "pod", "reason"})
	tbl.SetAutoWrapText(false)

	// Load table with data
	for _, res := range r {
		tbl.Append([]string{res.Pod.Namespace, res.Pod.Name, res.Reason.Error()})
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
