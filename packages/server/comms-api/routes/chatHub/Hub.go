// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package chatHub

import (
	"encoding/json"
	"log"
)

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	Clients map[*Client]bool

	// Inbound messages from the clients.
	Broadcast chan MessageItem

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client

	// quit
	Quit chan bool
}

type MessageItem struct {
	Username string `json:"username"`
	Message  string `json:"message"`
}

func NewHub() *Hub {
	return &Hub{
		Broadcast:  make(chan MessageItem),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		Clients:    make(map[*Client]bool),
		Quit:       make(chan bool),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case quit := <-h.Quit:
			if quit {
				log.Printf("Kill request received, shutting down")
				return
			}
		case client := <-h.register:
			h.Clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.Clients[client]; ok {
				delete(h.Clients, client)
				close(client.send)
			}
		case message := <-h.Broadcast:
			for client := range h.Clients {
				if client.username == message.Username {
					marshal, err := json.Marshal(message)
					if err != nil {
						log.Printf("Unable to marchal message for user: %s", message.Username)
					}
					select {
					case client.send <- marshal:
					default:
						close(client.send)
						delete(h.Clients, client)
					}
				}
			}
		}
	}
}
