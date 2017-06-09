package exam

import (
	"fmt"
	"testing"
	"time"

	"github.com/karts/scheduler"
)

type HelloMessage struct {
	Message string
}

type HelloActor struct {
	stateId int
	sch     *scheduler.Scheduler
}

func NewHelloActor(stateId int, scheduler *scheduler.Scheduler) *HelloActor {
	h := &HelloActor{stateId: stateId, sch: scheduler}

	return h
}

func (h *HelloActor) RegisterNextId(stateId int) {
	h.sch.Register(h.stateId, stateId, h.Receive)

}

func (h *HelloActor) Tell(e *scheduler.Event) {
	h.sch.Tell(e)
}

// self drive
func (h *HelloActor) Timer() {

}

// call func
func (h *HelloActor) Receive(e *scheduler.Event) error {
	fmt.Printf("%v receive: %v\n", h.stateId, e.Message)
	return nil
}

func TestRun(t *testing.T) {
	sch := scheduler.NewScheduler()
	actor01 := NewHelloActor(1, sch)
	actor01.RegisterNextId(2)

	actor02 := NewHelloActor(2, sch)
	actor02.RegisterNextId(3)

	actor03 := NewHelloActor(3, sch)
	actor03.RegisterNextId(0)

	sch.Start()

	e1 := &scheduler.Event{StateId: 2, Message: HelloMessage{Message: "I'm id 01"}}
	actor01.Tell(e1)

	e2 := &scheduler.Event{StateId: 3, Message: HelloMessage{Message: "I'm id 02"}}
	actor02.Tell(e2)
	actor02.Tell(e2)

	time.Sleep(time.Second * 1)
}
