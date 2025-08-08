package sse

import (
	"log"
	"sync"
)

type Hub struct {
	clients map[string]chan string
	mu      sync.Mutex
}

func NewHub() *Hub {
	return &Hub{
		clients: make(map[string]chan string),
	}
}

func (h *Hub) Register(searchID string) chan string {
	h.mu.Lock()
	defer h.mu.Unlock()

	ch := make(chan string, 2)
	h.clients[searchID] = ch
	log.Printf("Client registered for search_id: %s", searchID)
	return ch
}

func (h *Hub) Unregister(searchID string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if ch, ok := h.clients[searchID]; ok {
		close(ch)
		delete(h.clients, searchID)
		log.Printf("Client unregistered for search_id: %s", searchID)
	}
}

func (h *Hub) Broadcast(searchID string, message string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if ch, ok := h.clients[searchID]; ok {
		select {
		case ch <- message:
			log.Printf("Broadcasted message to search_id: %s", searchID)
		default:
			log.Printf("Channel for search_id %s is full. Message dropped.", searchID)
		}
	}
}
