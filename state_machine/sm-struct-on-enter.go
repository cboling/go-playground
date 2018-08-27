package main

//
//
//		Fourth test with go-fsm -> Use introspection to support an on_enter
//		functionality in the callback.
//
//  Looks like have to declare as interface and dig into reflect since it
//  needs to create a string.  The basic 'switch' statement in the callback
//  looks to be best.
//
import (
	"fmt"
	"github.com/theckman/go-fsm"
	"log"
	"sync"
	"time"
)

var waitGrp sync.WaitGroup

type entry struct {
	Name        string
	Synchronous bool
	FSM         *fsm.Machine
	// For use in any states.  All are canceled on entry to new state while SM mutex held
	timers []*time.Timer
}

func NewFsmEntry(to string, synchronous bool) (*entry, error) {
	d := &entry{
		Name:        to,
		Synchronous: synchronous,
		FSM:         &fsm.Machine{},
		timers:      make([]*time.Timer, 0),
	}
	log.Println("Setting up state transition rules")
	err := d.FSM.AddStateTransitionRules("", "initial")
	if err != nil {
		return nil, err
	}
	err = d.FSM.AddStateTransitionRules("disabled", "opened")
	if err != nil {
		return nil, err
	}
	err = d.FSM.AddStateTransitionRules("opened", "closing")
	if err != nil {
		return nil, err
	}
	err = d.FSM.AddStateTransitionRules("closing", "closed")
	if err != nil {
		return nil, err
	}
	err = d.FSM.AddStateTransitionRules("closed", "finished")
	if err != nil {
		return nil, err
	}
	err = d.FSM.AddStateTransitionRules("finished")
	if err != nil {
		return nil, err
	}
	////////////////////////////////////////////////////////////////////
	// For this implementation, the callback is called after transition
	// to the new state and when the FSM mutex is locked. If synchronous,
	// it will call directly into method. If async, uses a go routine to
	// call the callback
	//log.Println("Setting up state transition callbacks")

	err = d.FSM.SetStateTransitionCallback(d, d.Synchronous)
	if err != nil {
		return nil, err
	}
	log.Printf("Initial State before 'any' transition: %v", d.FSM.CurrentState())
	return d, nil
}

func (d *entry) cancelTimers() {
	var timer *time.Timer
	for _, timer = range d.timers {
		timer.Stop()
	}
	d.timers = make([]*time.Timer, 0)
}

////////////////////////////////////////////////////////////////////
// StateTransitionCallback is called after transition
// to the new state and when the FSM mutex is locked. If synchronous,
// it will call directly into method. If async, uses a go routine to
// call the callback
func (d *entry) StateTransitionCallback(state fsm.State) error {
	// Stop any timers
	d.cancelTimers()

	log.Printf("Entered callback for state: %v\n", state)

	switch state {
	case "opened":
		log.Println("Callback: In the opened state now")
		log.Println("Starting timer to fire in 1 second to start closing the entry")

		timer := time.NewTimer(1 * time.Second)
		d.timers = append(d.timers, timer)

		go func(t *time.Timer) {
			<-t.C
			err := d.FSM.StateTransition("closing")
			if err != nil {
				log.Fatalf("Failed to request that the entry start closing: %v", err)
			}
		}(timer)

	case "closing":
		log.Println("Callback: In the closing state now, scheduling the closed event (5 seconds)")

		timer := time.NewTimer(5 * time.Second)
		d.timers = append(d.timers, timer)

		go func(t *time.Timer) {
			<-t.C
			log.Println("Timer fired")
			err := d.FSM.StateTransition("closed")
			if err != nil {
				log.Fatalf("Failed to request that the entry be closed: %v", err)
			}
		}(timer)
		// Block here for 10 seconds to see how things work out in both sync/async
		// If async callbacks, the go routine above will fire and transition us to
		// finished and we never go to the last callback below. Probably not what you
		// would normally want

		log.Println("Callback: Sleeping 10 seconds to let timer fire")
		time.Sleep(10 * time.Second)

		// Next statement may not run if async callbacks are used.  Be aware!
		log.Println("Callback: done sleeping for 10 seconds")

	case "closed":
		log.Println("Callback: In the opened closed now, going straight to finished")

		// If callbacks are marked as synchronous, must perform the transition in
		// a go routine
		if d.Synchronous {
			go func() {
				err := d.FSM.StateTransition("finished")
				if err != nil {
					log.Fatalf("Failed to request an exit from this thing: %v", err)
				}
			}()
		} else {
			err := d.FSM.StateTransition("finished")
			if err != nil {
				log.Fatalf("Failed to request an exit from this thing: %v", err)
			}
		}
	case "finished":
		log.Println("Callback: In the finished state now")
		waitGrp.Done()
	}
	return nil
}

func main() {
	numDoors := 1
	synchronous := false
	var err error

	waitGrp.Add(numDoors)

	doors := make([]*entry, numDoors)
	for i := 0; i < numDoors; i++ {
		doors[i], err = NewFsmEntry(fmt.Sprintf("Door: %v", i), synchronous)
		if err != nil {
			log.Fatalln(err)
		}
	}
	log.Printf("Create %v FSM Door(s), synchronous: %v\n", numDoors, synchronous)

	for i := 0; i < numDoors; i++ {
		// Make sure the initial state is known (does not call a callback since
		// this sets default state
		err = doors[i].FSM.StateTransition("disabled")

		// Can use a timer here to have the completion routine to do the transition
		// to opened, or you can call the State Transition directly
		timer := time.NewTimer(1 * time.Second)
		doors[i].timers = append(doors[i].timers, timer)
		go func(d *entry, t *time.Timer) {
			<-t.C
			err := d.FSM.StateTransition("opened")
			if err != nil {
				log.Fatalf("Failed to request that the entry be closed the first time: %v", err)
			}
		}(doors[i], timer)
	}
	// Wait for completion
	log.Println("Waiting for all iterations and doors to complete")

	waitGrp.Wait()
	fmt.Println("Done....")
}
