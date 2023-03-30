package ws

import (
	"encoding/json"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

type Message struct {
	Action string `json:"action"`
	Data   []byte `json:"data"`
}

const (
	AddOrder = "add-order"
	AddMatch = "add-match"
)

func (m *Message) encode() ([]byte, error) {
	json, err := json.Marshal(m)
	if err != nil {
		return nil, errors.Wrap(err, "failed to encode websocket message")
	}
	return json, nil
}
