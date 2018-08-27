package main

//
//
//		Second test with go-fsm -> part of a structure
//
//
import (
	"fmt"
	"github.com/theckman/go-fsm"
	"log"
	"time"
)

type door struct {
	To   string
	FSM  *fsm.Machine
	done chan int

	// For use in any states.  All are canceled on entry to new state
	timers []*time.Timer
}

func NewFsmDoor(to string) (*door, error) {
	d := &door{
		To:     to,
		FSM:    &fsm.Machine{},
		done:   make(chan int, 2),
		timers: make([]*time.Timer, 0),
	}
	log.Println("Setting up state transition rules")
	err := d.FSM.AddStateTransitionRules("", "opened", "closed", "finished")
	if err != nil {
		return nil, err
	}
	err = d.FSM.AddStateTransitionRules("opened", "closed", "finished")
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

	log.Println("Setting up state transition callbacks")

	synchronous := true
	err = d.FSM.SetStateTransitionCallback(d, synchronous)
	if err != nil {
		return nil, err
	}
	log.Printf("Initial State: %v", d.FSM.CurrentState())
	return d, nil
}

func (d *door) cancelTimers() {
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
func (d *door) StateTransitionCallback(state fsm.State) error {
	log.Printf("Entered callback for state: %v\n", state)
	// Stop any timers
	d.cancelTimers()

	switch state {
	case "opened":
		log.Println("Callback: In the opened state now")
		log.Println("Starting timer to fire in 3 seconds to close the entry")
		timer := time.NewTimer(3 * time.Second)
		go func() {
			<-timer.C
			fmt.Println("Timer expired")
			err := d.FSM.StateTransition("closed")
			if err != nil {
				log.Fatalf("Failed to request that the entry be closed: %v", err)
			}
		}()

	case "closed":
		log.Println("Callback: In the opened closed now, launching exit event")
		go func() {
			err := d.FSM.StateTransition("finished")
			if err != nil {
				log.Fatalf("Failed to request an exit from this thing: %v", err)
			}
		}()

	case "finished":
		log.Println("Callback: In the opened finished now")
		d.done <- 1
	}
	return nil
}

func main() {
	door, err := NewFsmDoor("heaven")
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("Create FSM Door, initial state: %v\n", door.FSM.CurrentState())

	// Start things off by trying to enter our initial opened state
	err = door.FSM.StateTransition("opened")
	if err != nil {
		log.Printf("Failed to open an opened entry: %v\n", err)
	}
	log.Println("Starting timer to fire in 3 seconds to close the entry")

	timer := time.NewTimer(3 * time.Second)
	door.timers = append(door.timers, timer)

	go func(t *time.Timer) {
		<-t.C
		fmt.Println("Timer expired, requesting transition to closed")
		err = door.FSM.StateTransition("closed")
		if err != nil {
			fmt.Println(err)
		}
	}(timer)
	log.Println("Waiting for done to be done")
	<-door.done

	fmt.Println("Done")
}
