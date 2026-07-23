package core

type Event struct {
	Type    string
	Payload string
}

type EventInterface interface {
	GetType() string
	GetPayload() string
}

func (e Event) GetType() string {
	return e.Type
}

func (e Event) GetPayload() string {
	return e.Payload
}
