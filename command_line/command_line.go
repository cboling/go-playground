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

import (
	"flag"
	"fmt"
	"os"
)

// Run this with command line argumentL  -config=config.yaml
func main() {
	fmt.Printf("Before call:  Original OS Args: %v\n", os.Args)

	// Change it for a short period of time
	changeIt()

	// Check if it reverted
	fmt.Printf("After return: Original OS Args: %v\n", os.Args)
}

func parseArgs(configPath *string) {
	flag.StringVar(configPath, "config", "", "Configuration file")
	flag.Parse()
}

func changeIt() {
	oldArgs, restoreIt := ModifyCommandLine([]string{`-config="no-file.yaml"`})
	defer restoreIt()

	var configFile string
	parseArgs(&configFile)

	fmt.Println("+----------- inside  ---------------------------")
	fmt.Printf("|      Original os Args: %v\n", oldArgs)
	fmt.Printf("|           New os Args: %v\n", os.Args)
	fmt.Printf("|     config file value: %v\n", configFile)
	fmt.Println("+----------- inside  ---------------------------")
}

// ModifyCommandLine takes in a new array of command line arguments and
// replaces all OS command line arguments (starting after argument 0 which
// is the command name).
// It returns a copy of the original command line arguments and a function
// that can be called to restore the original values
func ModifyCommandLine(args []string) ([]string, func()) {
	oldArgs := make([]string, len(os.Args))
	copy(oldArgs, os.Args)
	os.Args = make([]string, 0)
	os.Args = append(os.Args, oldArgs[0])
	os.Args = append(os.Args, args...)
	return oldArgs, func() { copy(os.Args, oldArgs) }
}
