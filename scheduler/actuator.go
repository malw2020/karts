package scheduler

type Actuator interface {
	Receive(e *Event) error
}

