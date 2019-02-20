/*
 * Copyright (c) 2019 - present.  Boling Consulting Solutions (bcsw.net)
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
	"encoding/binary"
	"errors"
	"fmt"
	"log"
	"reflect"
)

// Given an array of interfaces as an interface, return octet buffer (in network
// byte order) with proper sizes.
//
// Note that this is for protocol oriented decoding and only handle basic types and
// all values are unsigned

func interfaceSliceToBytes(slice interface{}) ([]byte, error) {
	// First interface to []interface
	sValue := reflect.ValueOf(slice)
	if sValue.Kind() != reflect.Slice {
		return nil, errors.New("interfaceSliceToBytes() given a non-slice type")
	}
	// Now to a buffer
	results := make([]byte, 0)
	tmpBuffer := make([]byte, 8) // Large enough to hold anything below

	for i := 0; i < sValue.Len(); i++ {
		iface := sValue.Index(i).Interface()
		switch t := iface.(type) {
		default:
			return nil, errors.New(fmt.Sprintf("unsupported type: %v: (%v)", t, iface))

		//case octets:
		//	fallthrough
		//case byte:
		//	fallthrough
		case uint8: // Should handle 'byte' and 'octet' types
			results = append(results, t)

		case uint16:
			binary.BigEndian.PutUint16(tmpBuffer, t)
			results = append(results, tmpBuffer[0:2]...)

		case uint32:
			binary.BigEndian.PutUint32(tmpBuffer, t)
			results = append(results, tmpBuffer[0:4]...)

		case uint64:
			binary.BigEndian.PutUint64(tmpBuffer, t)
			results = append(results, tmpBuffer[0:8]...)

		case string:
			results = append(results, []byte(t)...)
		}
	}
	return results, nil
}

type messageBuffer struct {
	payload interface{}
}

type octets = uint8

func main() {
	octetsArray := []octets{0, 1, 2, 4, 5, 6, 7, 8, 9}
	iOctets := messageBuffer{octetsArray}
	results, err := interfaceSliceToBytes(iOctets.payload)
	if err != nil {
		log.Fatal(results)
	}
	log.Printf("Octets results: %v", results)

	byteArray := []byte{0, 1, 2, 4, 5, 6, 7, 8, 9}
	iBytes := messageBuffer{byteArray}
	results, err = interfaceSliceToBytes(iBytes.payload)
	if err != nil {
		log.Fatal(results)
	}
	log.Printf("Bytes results: %v", results)

	uArray := []uint8{0, 1, 2, 4, 5, 6, 7, 8, 9}
	uBytes := messageBuffer{uArray}
	results, err = interfaceSliceToBytes(uBytes.payload)
	if err != nil {
		log.Fatal(results)
	}
	log.Printf("Uint8s results: %v", results)

	wordArray := []uint16{0, 1, 2, 4, 5, 6, 7, 8, 9}
	iWords := messageBuffer{wordArray}
	results, err = interfaceSliceToBytes(iWords.payload)
	if err != nil {
		log.Fatal(results)
	}
	log.Printf("Words results: %v", results)

	dWordArray := []uint32{0, 1, 2, 4, 5, 6, 7, 8, 9}
	iDoubles := messageBuffer{dWordArray}
	results, err = interfaceSliceToBytes(iDoubles.payload)
	if err != nil {
		log.Fatal(results)
	}
	log.Printf("DWords results: %v", results)

	quadWordArray := []uint64{0, 1, 2, 4, 5, 6, 7, 8, 9}
	iQuads := messageBuffer{quadWordArray}
	results, err = interfaceSliceToBytes(iQuads.payload)
	if err != nil {
		log.Fatal(results)
	}
	log.Printf("QuadWords results: %v", results)

	stringArray := []string{"abc", "def", "xyz"}
	iStrings := messageBuffer{stringArray}
	results, err = interfaceSliceToBytes(iStrings.payload)
	if err != nil {
		log.Fatal(results)
	}
	log.Printf("Strings results: %v", results)

	log.Printf("Done")
}
