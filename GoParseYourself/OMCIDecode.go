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

package GoParseYourself

import (
	"encoding/hex"
	"fmt"
	"github.com/google/gopacket"
)

// TODO: Will focus on creation of an OMCI Frame decoder  (separate one for encode).

// TODO: Look into the go 'encoding' interfaces  -> import "encoding" and derive from BinaryMarshaller
//       and BinaryUnmarshaler

// TODO: Leverage design patterns from the 'gopacket:  https://godoc.org/github.com/google/gopacket#pkg-examples' library

// TODO: Take a look also at: https://github.com/ghedo/go.pkt

// TODO: Read over Brad Israel's blog article:  http://www.bisrael8191.com/Go-Packet-Sniffer/

func string_to_packet(input string) ([]byte, error) {
	var p []byte
	size := len(input)/2

	p, err := hex.DecodeString(input)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	fmt.Println(p)

	p2 := make([]byte, size, size)
	fmt.Println(p2)
	return p, nil
}


func main() {

	var MibResetRequest = "00014F0A000200000000000000000000" +
		"00000000000000000000000000000000" +
		"000000000000000000000028"

	data, err := string_to_packet(MibResetRequest)
	if err != nil {
		fmt.Println(err)
	} else {
		packet := gopacket.NewPacket(data, LayerTypeOMCI, gopacket.Lazy)
		fmt.Println(packet)
	}
}
