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

import "github.com/google/gopacket"

var (
	LayerTypeOMCI = gopacket.RegisterLayerType(1000,
		gopacket.LayerTypeMetadata{
			Name: "Frame",
			Decoder: gopacket.DecodeFunc(decodeOMCI),
	})
	LayerTypeOMCIPayload = gopacket.RegisterLayerType(1001,
		gopacket.LayerTypeMetadata{
			Name: "Payload",
			Decoder: gopacket.DecodeFunc(decodeOMCIPayload),
		})
	LayerTypeOMCITrailer = gopacket.RegisterLayerType(1002,
		gopacket.LayerTypeMetadata{
			Name: "Trailer",
			Decoder: gopacket.DecodeFunc(decodeOMCITrailer),
		})
)
