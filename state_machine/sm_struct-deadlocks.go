/*
 * Copyright (c) 2018 - present.  Boling Consulting Solutions (bcsw.net)
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

package main

///////////////////////////////////////////////////////////////////////////
//
// NOTE: This FSM library module works nicely but does deadlocks if
//       you try an use go routines to send events...  Trying another one...
//
///////////////////////////////////////////////////////////////////////////

import (
	"fmt"
	"github.com/looplab/fsm"
	"time"
)

type Door struct {
	To   string
	FSM  *fsm.FSM
	done chan int
}

func NewDoor(to string) *Door {
	d := &Door{
		To:   to,
		FSM:  nil,
		done: make(chan int, 2),
	}
	d.FSM = fsm.NewFSM(
		"closed",
		fsm.Events{
			{Name: "open", Src: []string{"closed"}, Dst: "open"},
			{Name: "close", Src: []string{"open"}, Dst: "closed"},
			{Name: "exit", Src: []string{"open", "closed"}, Dst: "finished"},
		},
		fsm.Callbacks{
			"enter_state": func(e *fsm.Event) { d.enterState(e) },
			"enter_finished": func(e *fsm.Event) {
				fmt.Printf("Entering finished state")
				d.done <- 1
			},
			"finished": func(e *fsm.Event) {
				fmt.Printf("now in finished state")
				d.done <- 1
			},
			"open": func(e *fsm.Event) {
				fmt.Println("After open event")
			},
			"close": func(e *fsm.Event) {
				fmt.Println("After close event")
				err := d.FSM.Event("exit")
				if err != nil {
					fmt.Println(err)
				}
			},
			"exit": func(e *fsm.Event) {
				fmt.Println("After exit event")
				d.done <- 1
			},
		},
	)
	return d
}

func (d *Door) enterState(e *fsm.Event) {
	fmt.Printf("The door to %s is %s\n", d.To, e.Dst)
}

func main() {
	door := NewDoor("heaven")

	err := door.FSM.Event("open")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Starting timer to fire in 3 seconds to close the door")
	timer := time.NewTimer(3 * time.Second)
	go func() {
		<-timer.C
		fmt.Println("Timer expired")
		err = door.FSM.Event("close")
		if err != nil {
			fmt.Println(err)
		}
	}()

	fmt.Println("Waiting for done to be done")
	<-door.done

	fmt.Println("Done")
}
