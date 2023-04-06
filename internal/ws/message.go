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
	AddOrder    = "add-order"
	AddMatch    = "add-match"
	UpdateOrder = "update-order"
	UpdateMatch = "update-match"
)

func (m *Message) encode() ([]byte, error) {
	data, err := json.Marshal(m)
	if err != nil {
		return nil, errors.Wrap(err, "failed to encode websocket message")
	}
	return data, nil
}
