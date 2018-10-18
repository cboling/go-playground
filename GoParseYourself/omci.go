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
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

type OMCIDeviceIdent byte

const (
	// Device Identifiers
	_                            = iota
	OMCIBaseline OMCIDeviceIdent = 0x0A // All G-PON OLTs and ONUs support the baseline message set
	OMCIExtended                 = 0x0B
)

// OMCIMsgType represents a DHCP operation
type OMCIMsgType byte

const (
	// Message Types
	_                                 = iota
	Create                OMCIMsgType = 4
	Delete                            = 6
	Set                               = 8
	Get                               = 9
	GetAllAlarms                      = 11
	GetAllAlarmsNext                  = 12
	MibUpload                         = 13
	MibUploadNext                     = 14
	MibReset                          = 15
	AlarmNotification                 = 16
	AttributeValueChange              = 17
	Test                              = 18
	StartSoftwareDownload             = 19
	DownloadSection                   = 20
	EndSoftwareDownload               = 21
	ActivateSoftware                  = 22
	CommitSoftware                    = 23
	SynchronizeTime                   = 24
	Reboot                            = 25
	GetNext                           = 26
	TestResult                        = 27
	GetCurrentData                    = 28
	SetTable                          = 29 // Defined in Extended Message Set Only
)

// OMCI defines the Baseline (not extended) protocol. Extended will be added once
// I can get basic working (and layered properly).  See ITU-T G.988 11/2017 section
// A.3 for more information
type OMCI struct {
	layers.BaseLayer
	TransactionID    uint16
	MessageType      OMCIMsgType
	DeviceIdentifier OMCIDeviceIdent
	EntityClass		 uint16
	EntityInstance	 uint16
	//Data             []byte		 // Octets 8:39
	//OMCITrailer	     []byte      // Octets 40:47
}


func (omci *OMCI) String() string {
	return fmt.Sprintf("OMCI %v: (%v/%v)", omci.MessageType,
		omci.EntityClass, omci.EntityInstance)
}

// LayerType returns LayerTypeOMCI
func (omci *OMCI) LayerType() gopacket.LayerType {
	return LayerTypeOMCI
}

func (omci *OMCI) CanDecode() gopacket.LayerClass {
	return LayerTypeOMCI
}

// NextLayerType returns the layer type contained by this DecodingLayer.
func (omci *OMCI) NextLayerType() gopacket.LayerType {
	return LayerTypeOMCIPayload
}

func decodeOMCI(data []byte, p gopacket.PacketBuilder) error {
	omci := &OMCI{}
	return omci.DecodeFromBytes(data, p)
}

func (omci *OMCI) DecodeFromBytes(data []byte, p gopacket.PacketBuilder) error {
	if len(data) < 8 {
		return errors.New("OMCI packet header too small")
	}
	omci.TransactionID 	  = binary.BigEndian.Uint16(data[0:2])
	omci.MessageType      = OMCIMsgType(data[3])
	omci.DeviceIdentifier = OMCIDeviceIdent(data[4])
	omci.EntityClass      = binary.BigEndian.Uint16(data[4:6])
	omci.EntityInstance   = binary.BigEndian.Uint16(data[6:8])
	return p.NextDecoder(LayerTypeOMCIPayload)
}

// SerializeTo writes the serialized form of this layer into the
// SerializationBuffer, implementing gopacket.SerializableLayer.
// See the docs for gopacket.SerializableLayer for more info.
func (omci *OMCI) SerializeTo(b gopacket.SerializeBuffer, opts gopacket.SerializeOptions) error {
	// Basic (common) OMCI Header is 8 octets, 10
	bytes, err := b.PrependBytes(8)
	if err != nil {
		return err
	}
	binary.BigEndian.PutUint16(bytes, omci.TransactionID)
	bytes[2] = byte(omci.MessageType)
	bytes[3] = byte(omci.DeviceIdentifier)
	binary.BigEndian.PutUint16(bytes[4:], omci.EntityClass)
	binary.BigEndian.PutUint16(bytes[6:], omci.EntityInstance)

	length := 48	// TODO: Only Baseline Messages currently supported
	padding, err := b.AppendBytes(length - 8)
	if err != nil {
		return err
	}
	copy(padding, lotsOfZeros[:])
	return nil
}

// hacky way to zero out memory... there must be a better way?
var lotsOfZeros [48]byte