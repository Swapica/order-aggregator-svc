package responses

import (
	"encoding/json"

	"gitlab.com/distributed_lab/logan/v3/errors"
)

func toRawMetaField(count int64) json.RawMessage {
	c := struct {
		Count int64 `json:"count"`
	}{Count: count}
	bytes, err := json.Marshal(c)
	if err != nil {
		panic(errors.Wrap(err, "unexpected error on marshalling count metadata"))
	}
	return bytes
}
