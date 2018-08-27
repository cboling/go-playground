package main

//
//
//		Third test with go-fsm -> Lots of repetitive iterations
//
//  Works well with both sync and async/go-routine callback support. The
//  synchronous calls the 'callback' while holding the state machine lock
//
import (
	"fmt"
	"github.com/theckman/go-fsm"
	"log"
	"math/rand"
	"sync"
	"time"
)

var wg sync.WaitGroup

type entryWay struct {
	Name       string
	FSM        *fsm.Machine
	swingsLeft int

	// For use in any states.  All are canceled on entry to new state
	timers []*time.Timer
}

func NewFSMDoor(to string, swings int, synchronous bool) (*entryWay, error) {
	d := &entryWay{
		Name:       to,
		FSM:        &fsm.Machine{},
		timers:     make([]*time.Timer, 0),
		swingsLeft: swings,
	}
	// log.Println("Setting up state transition rules")
	err := d.FSM.AddStateTransitionRules("", "opened", "closed", "finished")
	if err != nil {
		return nil, err
	}
	err = d.FSM.AddStateTransitionRules("opened", "closed", "finished")
	if err != nil {
		return nil, err
	}
	err = d.FSM.AddStateTransitionRules("closed", "opened", "finished")
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

	err = d.FSM.SetStateTransitionCallback(d, synchronous)
	if err != nil {
		return nil, err
	}
	//log.Printf("Initial State: %v", d.FSM.CurrentState())
	return d, nil
}

func (d *entryWay) cancelTimers() {
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
func (d *entryWay) StateTransitionCallback(state fsm.State) error {
	// log.Printf("Entered callback for state: %v\n", state)
	// Stop any timers
	d.cancelTimers()

	switch state {
	case "opened":
		// log.Println("Callback: In the opened state now")
		// log.Println("Starting timer to fire in 3 seconds to close the entry")
		variance := time.Duration(rand.Int63n(4)) * time.Millisecond
		timer := time.NewTimer(20*time.Millisecond + variance)
		d.timers = append(d.timers, timer)

		go func(t *time.Timer) {
			<-t.C
			err := d.FSM.StateTransition("closed")
			if err != nil {
				log.Fatalf("Failed to request that the entry be closed: %v", err)
			}
		}(timer)

	case "closed":
		// log.Println("Callback: In the opened closed now, launching exit event")
		d.swingsLeft--

		if d.swingsLeft <= 0 {
			go func() {
				err := d.FSM.StateTransition("finished")
				if err != nil {
					log.Fatalf("Failed to request an exit from this thing: %v", err)
				}
			}()
		} else {
			if d.swingsLeft%100 == 0 {
				log.Printf("Door %v has %v left", d.Name, d.swingsLeft)
			}
			variance := time.Duration(rand.Int63n(4)) * time.Millisecond
			timer := time.NewTimer(10*time.Millisecond + variance)
			d.timers = append(d.timers, timer)

			go func(t *time.Timer) {
				<-t.C
				err := d.FSM.StateTransition("opened")
				if err != nil {
					log.Fatalf("Failed to request that the entry be opened: %v", err)
				}
			}(timer)
		}

	case "finished":
		log.Println("Callback: In the opened finished now")
		wg.Done()
	}
	return nil
}

func main() {
	numDoors := 1000
	numIterations := 10000
	//numDoors := 1
	//numIterations := 10
	//synchronous := true
	synchronous := false
	var err error

	wg.Add(numDoors)

	doors := make([]*entryWay, numDoors)
	for i := 0; i < numDoors; i++ {
		doors[i], err = NewFSMDoor(fmt.Sprintf("Door: %v", i), numIterations, synchronous)
		if err != nil {
			log.Fatalln(err)
		}
	}
	log.Printf("Create %v FSM Doors, synchronous: %v\n", numDoors, synchronous)
	// Make sure the state is known (does not call a callback since this sets default state
	for i := 0; i < numDoors; i++ {
		err = doors[i].FSM.StateTransition("opened")

		timer := time.NewTimer(100 * time.Millisecond)
		doors[i].timers = append(doors[i].timers, timer)
		go func(d *entryWay, t *time.Timer) {
			<-t.C
			err := d.FSM.StateTransition("closed")
			if err != nil {
				log.Fatalf("Failed to request that the entry be closed the first time: %v", err)
			}
		}(doors[i], timer)
	}
	// Wait for completion
	log.Println("Waiting for all iterations and doors to complete")

	wg.Wait()
	fmt.Println("Done....")
}
