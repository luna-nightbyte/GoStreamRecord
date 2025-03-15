package connection

import (
	"errors"
	"sync"
)

// Use make new-notify-token
var authToken = "21b3fc1def893f4dc1d619acbe6635b12c8605b1b93f81b165a4bfa745bebc73"

// Notifier manages listeners and message dispatch.
type Notifier struct {
	listeners map[chan<- []byte]struct{}
	mutex     sync.Mutex
	token     string
}

// Source defines an interface for subscribing to notifications.
type Source interface {
	Subscribe(listener chan<- []byte) (UnsubscribeFunc, error)
}

// UnsubscribeFunc defines the function to remove a listener.
type UnsubscribeFunc func() error

// NewNotifier initializes a new Notifier with the required auth token.
func NewNotifier() *Notifier {
	return &Notifier{
		listeners: make(map[chan<- []byte]struct{}),
		token:     authToken,
	}
}

// Subscribe adds a new listener if the provided token is valid.
// It returns an unsubscribe function that removes the listener.
func (n *Notifier) Connect(listener chan<- []byte) (UnsubscribeFunc, error) {
	if n.token != authToken {
		return nil, errors.New("unauthorized: invalid auth token")
	}

	n.mutex.Lock()
	n.listeners[listener] = struct{}{}
	n.mutex.Unlock()

	// Return an unsubscribe function to remove the listener and close the channel.
	unsubscribe := func() error {
		n.mutex.Lock()
		delete(n.listeners, listener)
		close(listener)
		n.mutex.Unlock()
		return nil
	}
	return unsubscribe, nil
}

// Notify sends a message to all listeners using non-blocking sends.
func (n *Notifier) Notify(msg []byte) error {
	if n.token != authToken {
		return errors.New("unauthorized: invalid auth token")
	}

	n.mutex.Lock()
	defer n.mutex.Unlock()

	for ch := range n.listeners {
		select {
		case ch <- msg:
			// Message sent successfully.
		default:
			// Skip sending if the channel is full.
		}
	}
	return nil
}
