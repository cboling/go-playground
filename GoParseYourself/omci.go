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
	"github.com/google/gopacket"
	layers "github.com/google/gopacket/layers"
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
	// OMCI Message Type specific data begins here.  It occupies bytes 9-40
	Data             []byte
	OMCITrailer		 uint32
}


// CanDecode returns the set of layer types that this DecodingLayer can decode.
func (omci *OMCI) CanDecode() gopacket.LayerClass {
	return LayerTypeOMCI
}

// NextLayerType returns the layer type contained by this DecodingLayer.
func (omci *OMCI) NextLayerType() gopacket.LayerType {
	return gopacket.LayerTypePayload
}

func decodeOMCI(data []byte, p gopacket.PacketBuilder) error {

	omci := &OMCI{}
	return layers.decodingLayerDecoder(omci, data, p)	// TODO: Where is our exposed access function ?
}
