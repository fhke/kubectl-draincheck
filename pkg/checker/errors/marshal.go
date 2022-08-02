package errors

import "encoding/json"

func (e *errorImpl) MarshalJSON() ([]byte, error) {
	return json.Marshal(e.error.Error())
}
