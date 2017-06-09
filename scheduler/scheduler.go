package scheduler

import "fmt"

type CallActuatorFunc func(e *Event) error

type Scheduler struct {
	stateMachine map[int][]int
	stateChan    map[int]chan *Event
	stateCall    map[int]CallActuatorFunc
}

func NewScheduler() *Scheduler {
	scheuler := &Scheduler{}

	scheuler.stateMachine = make(map[int][]int)
	scheuler.stateChan = make(map[int]chan *Event)
	scheuler.stateCall = make(map[int]CallActuatorFunc)

	return scheuler
}

func (sch*Scheduler) Start() error {
	for id, _ := range sch.stateChan {
		go sch.act(id)
	}

	return nil
}

func (sch*Scheduler) Stop() error {
	return nil
}

func (sch*Scheduler) Tell(event *Event) {
	fmt.Printf("Tell id:%v\n", event.StateId)

	stateCh, ok := sch.stateChan[event.StateId]
	if ok {
		stateCh <- event
	}
}

func (sch*Scheduler)  Register(stateId int, nextStateId int, act CallActuatorFunc) {
	_, ok := sch.stateMachine[stateId]
	if !ok {
		sch.stateMachine[stateId] = make([]int, 0, statesCapacity)
	}

	if len(sch.stateMachine[stateId]) == 0 {
		sch.stateMachine[stateId] = append(sch.stateMachine[stateId], nextStateId)
	} else {
		var index int
		for index = 0; index < len(sch.stateMachine[stateId]); index++ {
			if sch.stateMachine[stateId][index] == nextStateId {
				break
			}
		}
		if index >= len(sch.stateMachine[stateId]) {
			sch.stateMachine[stateId] = append(sch.stateMachine[stateId], nextStateId)
		}
	}

	_, ok = sch.stateChan[nextStateId]
	if !ok {
		sch.stateChan[nextStateId] = make(chan *Event, stateChanBuf)
	}

	_, ok = sch.stateCall[nextStateId]
	if !ok {
		sch.stateCall[nextStateId] = act
	}

}

func (sch*Scheduler) act(stateId int) {
	chanBuf, ok := sch.stateChan[stateId]
	if !ok {
		return
	}

	for {
		data, ok := <-chanBuf
		if !ok {
			continue
		}

		nextStates, ok := sch.stateMachine[stateId]
		if !ok {
			continue
		}

		for _, id := range nextStates {
			fmt.Printf("stateId:%v, call next state id:%v\n", stateId, id)
			callFunc, ok := sch.stateCall[id]
			if !ok {
				continue
			}

			callFunc(data)
		}
	}
}


