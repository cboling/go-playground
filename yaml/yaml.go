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

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

const (
	defaultConfigFilePath = "./config.yaml"
	defaultCPUCount       = 0     // Maximum CPUs to span across. 0 = all CPUs on system
	defaultPonCount       = 16    // PONs are numbered 0..n
	defaultMaxOnuPerPon   = 256   // Do not allow more than these many active on a single PON
	defaultRestPort       = 55055 // Listen for REST requests on this port, 0 to disable

)

// Config provides for the MOCK ONU runtime settings
type Config struct {
	ConfigPath   string `yaml:"_"`
	MaxCpus      int    `yaml:"maxCpus"`
	PonCount     int    `yaml:"ponCount"`
	MaxOnuPerPon int    `yaml:"maxOnuPerPon"`
	RestPort     int    `yaml:"restPort"`

	PonConfig PonConfig `yaml:"ponConfig"`
}

const (
	defaultPonEnablePattern = "Enabled"
	defaultPonDelayTime     = 10.0
	defaultPonEnableCount   = 1
	defaultPonReverseOrder  = false
)

type PonConfig struct {
	// "Enabled"    All at once (step function)
	// "Disabled"   Leave all disabled
	// "Sequential  Enable PONs sequentially
	// "Random"	   Enable PONs Random
	EnablePattern string  `yaml:"pattern"`
	Delay         float32 `yaml:"delay"`   // Seconds between enables if sequential/random
	Count         int     `yaml:"count"`   // Number of pons per enable if sequential/random
	ReverseOrder  bool    `yaml:"reverse"` // If sequential, start at highest PON Id
}

func main() {
	path := "config.yaml"
	configuration, err := ReadConfig(path)

	if err != nil {
		fmt.Printf("Invalid input file: '%v': %v\n", path, err)
		os.Exit(1)
	}
	fmt.Printf("Read config file: '%v'", configuration)
}

func ReadConfig(path string) (Config, error) {
	// Load in a configuration YAML file if available
	c := Config{defaultConfigFilePath,
		defaultCPUCount,
		defaultPonCount,
		defaultMaxOnuPerPon,
		defaultRestPort,
		PonConfig{
			defaultPonEnablePattern,
			defaultPonDelayTime,
			defaultPonEnableCount,
			defaultPonReverseOrder,
		},
	}
	yamlData, err := ioutil.ReadFile(path)
	if err != nil {
		return c, err
	}
	err = yaml.Unmarshal(yamlData, &c)
	return c, err
}
