/*
 * Copyright 2018 - present.  Boling Consulting Solutions (bcsw.net)
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
 */
package main

///////////////////////////////////////////////////////////////////////////
//
// NOTE: This library module works nicely but does deadlocks if you try an
//       use go routines to send events...
//
///////////////////////////////////////////////////////////////////////////

import (
	"fmt"

	"github.com/looplab/fsm"
)

func init() {

}

// Test out the FSM state machine package
func main() {

	libraryExample()
	anotherExample()
}

func libraryExample() {
	fmt.Println()
	fmt.Println("Starting Library Example")
	fsm1 := fsm.NewFSM(
		"closed",
		fsm.Events{
			{Name: "open", Src: []string{"closed"}, Dst: "open"},
			{Name: "close", Src: []string{"open"}, Dst: "closed"},
		},
		fsm.Callbacks{},
	)
	fmt.Println(fsm1.Current())
	err := fsm1.Event("open")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(fsm1.Current())
	err = fsm1.Event("close")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(fsm1.Current())

	fmt.Println()
	fmt.Println("End of Library Example")
	fmt.Println("========================================================================")
	fmt.Println()
}

func anotherExample() {
	fmt.Println()
	fmt.Println("Starting Another Example")
	fsm2 := fsm.NewFSM(
		"locked",
		fsm.Events{
			{Name: "unlock", Src: []string{"locked"}, Dst: "unlocked"},
			{Name: "open", Src: []string{"unlocked"}, Dst: "open"},
			{Name: "close", Src: []string{"open"}, Dst: "unlocked"},
			{Name: "lock", Src: []string{"unlocked"}, Dst: "locked"},
		},
		fsm.Callbacks{},
	)

	fsm2.Event("unlock")
	fsm2.Event("open")
	fsm2.Event("close")
	fsm2.Event("lock")
	fmt.Println(fsm2.Current())

	fmt.Println()
	fmt.Println("End of Another Example")
	fmt.Println("========================================================================")
	fmt.Println()
}
