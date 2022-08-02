package errors

import "encoding/json"

// An error type that is handled correctly by json.Marshal()
type errorImpl struct {
	error
}

var _ error = &errorImpl{}
var _ json.Marshaler = &errorImpl{}
