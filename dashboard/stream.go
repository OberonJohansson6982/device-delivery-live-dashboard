package dashboard

import "sync"

type Stream struct {
	mu          sync.Mutex
	nextID      int
	subscribers map[int]chan DeviceEvent
}

func NewStream() *Stream {
	return &Stream{subscribers: make(map[int]chan DeviceEvent)}
}

func (s *Stream) Subscribe() (<-chan DeviceEvent, func()) {
	s.mu.Lock()
	id := s.nextID
	s.nextID++
	ch := make(chan DeviceEvent, 8)
	s.subscribers[id] = ch
	s.mu.Unlock()

	return ch, func() {
		s.mu.Lock()
		delete(s.subscribers, id)
		s.mu.Unlock()
	}
}

func (s *Stream) Send(event DeviceEvent) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, ch := range s.subscribers {
		select {
		case ch <- event:
		default:
		}
	}
}
