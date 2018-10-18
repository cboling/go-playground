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
package omci

import (
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

// Frame defines the Baseline (not extended) protocol. Extended will be added once
// I can get basic working (and layered properly).  See ITU-T G.988 11/2017 section
// A.3 for more information
type Trailer struct {
	layers.BaseLayer
	OMCITrailer uint64
}

func (omci *Trailer) String() string {
	return fmt.Sprintf("Trailer: (%v)", omci.OMCITrailer)
}

// LayerType returns LayerTypeOMCI
func (omci *Trailer) LayerType() gopacket.LayerType {
	return LayerTypeOMCITrailer
}

func (omci *Trailer) CanDecode() gopacket.LayerClass {
	return LayerTypeOMCITrailer
}

// NextLayerType returns the layer type contained by this DecodingLayer.
func (omci *Trailer) NextLayerType() gopacket.LayerType {
	return gopacket.LayerTypePayload
}

func decodeOMCITrailer(data []byte, p gopacket.PacketBuilder) error {
	omci := &Trailer{}
	return omci.DecodeFromBytes(data, p) // TODO: Where is our exposed access function ?
}

func (omci *Trailer) DecodeFromBytes(data []byte, p gopacket.PacketBuilder) error {
	if len(data) < 8 {
		return errors.New("frame trailer too small")
	}
	omci.OMCITrailer = binary.BigEndian.Uint64(data[0:2])
	return p.NextDecoder(LayerTypeOMCIPayload)
}

// SerializeTo writes the serialized form of this layer into the
// SerializationBuffer, implementing gopacket.SerializableLayer.
// See the docs for gopacket.SerializableLayer for more info.
func (omci *Trailer) SerializeTo(b gopacket.SerializeBuffer, opts gopacket.SerializeOptions) error {
	// Basic (common) Frame Header is 8 octets, 10
	bytes, err := b.PrependBytes(8)
	if err != nil {
		return err
	}
	binary.BigEndian.PutUint64(bytes, omci.OMCITrailer)
	return nil
}
