package ws

import (
	"encoding/json"

	"github.com/Swapica/order-aggregator-svc/internal/config"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan []byte

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client

	cfg config.Config
}

func NewHub(cfg config.Config) *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
		cfg:        cfg,
	}
}

func (h *Hub) addClient(client *Client) {
	h.clients[client] = true
}

func (h *Hub) removeClient(client *Client) {
	if _, ok := h.clients[client]; ok {
		delete(h.clients, client)
	}
}

func (h *Hub) BroadcastToClients(action string, order interface{}) error {
	if order == nil {
		return errors.New("an empty order is given")
	}

	orderData, err := json.Marshal(order)
	if err != nil {
		h.cfg.Log().WithError(err).WithFields(logan.F{"action": action, "order": order})
		return errors.Wrap(err, "failed to marshal order")
	}

	message := Message{
		Action: action,
		Data:   orderData,
	}

	jsonMessage, err := message.encode()
	if err != nil {
		h.cfg.Log().WithError(err).WithFields(logan.F{"action": action, "order": order})
		return errors.Wrap(err, "failed to marshal websocket message")
	}

	h.cfg.Log().Debug("try send websocket message to clients")
	for client := range h.clients {
		client.send <- jsonMessage
	}

	return nil
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.addClient(client)

		case client := <-h.unregister:
			h.removeClient(client)
		}
	}
}
