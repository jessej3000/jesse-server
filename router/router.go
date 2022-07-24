package router

import (
	"fmt"
)

// Handler type
type Handler func(*Client, interface{}, *Manager)

// Router struct type
type Router struct {
	Rules map[string]Handler
}

// NewMessageRouter creates a new router
func NewMessageRouter() *Router {
	return &Router{
		Rules: make(map[string]Handler),
	}
}

// Handle adds a handler into the rules map
func (r *Router) Handle(msgName string, handler Handler) {
	r.Rules[msgName] = handler
	fmt.Println(msgName + " handler added...")
}

// FindHandler finds handler in the map
func (r *Router) FindHandler(msgName string) (Handler, bool) {
	handler, found := r.Rules[msgName]
	return handler, found
}
