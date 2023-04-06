package ContactHub

import (
	"encoding/json"
	"github.com/openline-ai/openline-customer-os/packages/server/comms-api/model"
	"log"
)

type Time struct {
	Seconds string `json:"seconds"`
	Nanos   string `json:"nanos"`
}

type ContactEvent struct {
	ContactId        string                       `json:"contactId"`
	InteractionEvent model.InteractionEventCreate `json:"event"`
}

// ContactHub Hub maintains the set of active clients and broadcasts messages to the
// clients.
type ContactHub struct {
	// Registered clients.
	Clients map[*ContactClient]bool

	// Inbound messages from the clients.
	Broadcast chan ContactEvent

	// Register requests from the clients.
	Register chan *ContactClient

	// Unregister requests from clients.
	unregister chan *ContactClient

	Quit chan bool
}

func NewContactHub() *ContactHub {
	return &ContactHub{
		Broadcast:  make(chan ContactEvent),
		Register:   make(chan *ContactClient),
		unregister: make(chan *ContactClient),
		Clients:    make(map[*ContactClient]bool),
		Quit:       make(chan bool),
	}
}

func (h *ContactHub) Run() {
	for {
		select {
		case quit := <-h.Quit:
			if quit {
				log.Printf("Kill request received, shutting down")
				return
			}
		case client := <-h.Register:
			log.Printf("Registered: " + client.contactId)
			h.Clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.Clients[client]; ok {
				delete(h.Clients, client)
				close(client.send)
			}
		case message := <-h.Broadcast:
			for client := range h.Clients {
				if client.contactId == message.ContactId {
					byteMsg, err := json.Marshal(message.InteractionEvent)
					if err != nil {
						log.Printf("Unable to marchal event for contact: %s, reason: %s", message.ContactId, err)
						continue
					}
					select {
					case client.send <- byteMsg:
					default:
						close(client.send)
						delete(h.Clients, client)
					}
				}
			}
		}
	}
}
