package main

//
//
//		First test with go-fsm -> stand alone
//
//
import (
	"fmt"
	"github.com/theckman/go-fsm"
)

type T struct {
	M *fsm.Machine
}

func main() {
	t := &T{M: &fsm.Machine{}}

	// add initial rule
	err := t.M.AddStateTransitionRules("started", "finished", "aborted", "exited")

	if err != nil {
		// handle
	}
	fmt.Printf("Initial State: %v\n", t.M.CurrentState())

	// add rest of rules
	t.M.AddStateTransitionRules("aborted", "started")
	t.M.AddStateTransitionRules("finished", "started")
	t.M.AddStateTransitionRules("exited") // final state

	fmt.Printf("State now: %v\n", t.M.CurrentState())

	// set initial state
	err = t.M.StateTransition("aborted") // nil

	// get the current state
	fmt.Printf("State should be aborted: %v\n", t.M.CurrentState()) // "aborted"

	// try to transition to an non-whitelisted state
	err = t.M.StateTransition("finished") // ErrTransitionNotPermitted

	// try to transition to a permitted state
	err = t.M.StateTransition("started") // nil

	fmt.Printf("State finally: %v\n", t.M.CurrentState())
}
