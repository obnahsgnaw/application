package event

import (
	"github.com/asaskevich/EventBus"
)

// Manger event manager
type Manger struct {
	bus           EventBus.Bus
	eventHandlers []*Handler
}

// Handler event handler
type Handler struct {
	Topic  string
	Handle func(event *Event)
}

// Event the event target
type Event struct {
	Topic   string
	Data    []interface{}
	manager *Manger
}

// Init managed event
func (m *Manger) init() error {
	for _, handler := range m.eventHandlers {
		if err := m.bus.SubscribeAsync(handler.Topic, handler.Handle, false); err != nil {
			return err
		}
	}

	return nil
}

func Init(m *Manger) error {
	return m.init()
}

// Register a event
func (m *Manger) Register(topic string, handle func(e *Event)) {
	m.eventHandlers = append(m.eventHandlers, &Handler{
		Topic:  topic,
		Handle: handle,
	})
}

// Fire event
func (m *Manger) Fire(e *Event) {
	m.bus.Publish(e.Topic, e)
	m.bus.WaitAsync()
}

// NewEvent return a new event
func (m *Manger) NewEvent(topic string, data []interface{}) *Event {
	return &Event{
		Topic:   topic,
		Data:    data,
		manager: m,
	}
}

// New return a new event manager
func New() *Manger {
	return &Manger{bus: EventBus.New()}
}

// Fire to fire event
func (e *Event) Fire() {
	e.manager.Fire(e)
}
