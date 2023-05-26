package rpc

import (
	"fmt"
	"github.com/pkg/errors"
	"sync"
)

type Router struct {
	lock     sync.Mutex
	handlers map[string]Handler
}

func NewRouter() *Router {
	return &Router{handlers: map[string]Handler{}}
}

func (r *Router) Handle(method string, body []byte) error {
	r.lock.Lock()
	handler, ok := r.handlers[method]
	r.lock.Unlock()
	if !ok {
		return errors.New("unknown method")
	}
	return handler.Handle(method, body)
}

func (r *Router) Register(handlers map[string]Handler) error {
	r.lock.Lock()
	defer r.lock.Unlock()

	for method, _ := range handlers {
		if _, ok := r.handlers[method]; ok {
			return fmt.Errorf("method already registered: %s", method)
		}
	}

	for method, handler := range handlers {
		r.handlers[method] = handler
	}

	return nil
}
