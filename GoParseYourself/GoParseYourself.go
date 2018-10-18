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
	"encoding/hex"
	"fmt"
	"github.com/cboling/go-playground/GoParseYourself/omci"
	"github.com/google/gopacket"
)

// TODO: Will focus on creation of an OMCI Frame decoder  (separate one for encode).

// TODO: Take a look also at: https://github.com/ghedo/go.pkt

// TODO: Read over Brad Israel's blog article:  http://www.bisrael8191.com/Go-Packet-Sniffer/

func stringToPacket(input string) ([]byte, error) {
	var p []byte

	p, err := hex.DecodeString(input)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return p, nil
}


func main() {

	var MibResetRequest = "00014F0A000200000000000000000000" +
		"00000000000000000000000000000000" +
		"000000000000000000000028"

	data, err := stringToPacket(MibResetRequest)
	if err != nil {
		fmt.Println(err)
	} else {
		packet := gopacket.NewPacket(data, omci.LayerTypeOMCI, gopacket.Lazy)
		fmt.Println(packet)
	}
}
