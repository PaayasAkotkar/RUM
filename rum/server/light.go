package rum

// Flow:
//                 Track the internal Subscriber
//                       & Pass the publish result
import (
	"log"
	"sync"
)

// Light is a generic pub/sub bus keyed by ISequence[In].
// T is the message payload type.
type Light[In, T any] struct {
	mu          sync.Mutex
	subscribers map[ISequence[In]]map[chan *T]struct{}
}

func NewLight[In, T any]() *Light[In, T] {
	return &Light[In, T]{
		subscribers: make(map[ISequence[In]]map[chan *T]struct{}, 1),
	}
}

// Subscribe returns a buffered channel that receives published values for key
func (l *Light[In, T]) Subscribe(key ISequence[In]) chan *T {
	ch := make(chan *T, 1)
	l.mu.Lock()
	defer l.mu.Unlock()
	if _, ok := l.subscribers[key]; !ok {
		l.subscribers[key] = make(map[chan *T]struct{}, 1)
	}
	l.subscribers[key][ch] = struct{}{}
	return ch
}

// Publish sends val to all subscribers of key. Never blocks.
func (l *Light[In, T]) Publish(key ISequence[In], val *T) {
	log.Println("in publish")
	l.mu.Lock()
	defer l.mu.Unlock()
	for ch := range l.subscribers[key] {
		select {
		case ch <- val:
		default:
		}
	}
}

// Unsub removes and closes the channel for key
func (l *Light[In, T]) Unsub(key ISequence[In], ch chan *T) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if sub, ok := l.subscribers[key]; ok {
		delete(sub, ch)
		close(ch)
		if len(sub) == 0 {
			delete(l.subscribers, key)
		}
	}
}
