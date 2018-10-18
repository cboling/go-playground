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
	"fmt"
	"github.com/google/gopacket"
)

// OMCI defines the Baseline (not extended) protocol. Extended will be added once
// I can get basic working (and layered properly).  See ITU-T G.988 11/2017 section
// A.3 for more information
type OMCIPayload struct {
	Data             []byte		 // Octets 8:39
}


func (omci *OMCIPayload) String() string {
	return fmt.Sprintf("OMCIPayload")
}

// LayerType returns LayerTypeOMCI
func (omci *OMCIPayload) LayerType() gopacket.LayerType {
	return LayerTypeOMCI
}

func (omci *OMCIPayload) CanDecode() gopacket.LayerClass {
	return LayerTypeOMCI
}

// NextLayerType returns the layer type contained by this DecodingLayer.
func (omci *OMCIPayload) NextLayerType() gopacket.LayerType {
	return LayerTypeOMCITrailer
}

func decodeOMCIPayload(data []byte, p gopacket.PacketBuilder) error {
	omci := &OMCIPayload{}
	return omci.DecodeFromBytes(data, p)	// TODO: Where is our exposed access function ?
}

func (omci *OMCIPayload) DecodeFromBytes(data []byte, p gopacket.PacketBuilder) error {

	//omci.TransactionID 	  = binary.BigEndian.Uint16(data[0:2])
	//omci.MessageType      = OMCIMsgType(data[3])
	//omci.DeviceIdentifier = OMCIDeviceIdent(data[4])
	//omci.EntityClass      = binary.BigEndian.Uint16(data[4:6])
	//omci.EntityInstance   = binary.BigEndian.Uint16(data[6:8])
	return p.NextDecoder(LayerTypeOMCIPayload)
}

// SerializeTo writes the serialized form of this layer into the
// SerializationBuffer, implementing gopacket.SerializableLayer.
// See the docs for gopacket.SerializableLayer for more info.
func (omci *OMCIPayload) SerializeTo(b gopacket.SerializeBuffer, opts gopacket.SerializeOptions) error {
	// Basic (common) OMCI Header is 8 octets, 10
	//bytes, err := b.PrependBytes(8)
	//if err != nil {
	//	return err
	//}
	//binary.BigEndian.PutUint16(bytes, omci.TransactionID)
	//bytes[2] = byte(omci.MessageType)
	//bytes[3] = byte(omci.DeviceIdentifier)
	//binary.BigEndian.PutUint16(bytes[4:], omci.EntityClass)
	//binary.BigEndian.PutUint16(bytes[6:], omci.EntityInstance)
	//
	//length := 48	// TODO: Only Baseline Messages currently supported
	//padding, err := b.AppendBytes(length - 8)
	//if err != nil {
	//	return err
	//}
	//copy(padding, lotsOfZeros[:])
	return nil
}
