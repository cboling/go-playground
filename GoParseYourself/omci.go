package GoParseYourself

import "github.com/google/gopacket/layers"

package layers

import (
"encoding/binary"
"fmt"
"github.com/google/gopacket"
"layers"
)

// EAPOL defines an EAP over LAN (802.1x) layer.
type EAPOL struct {
	layers.BaseLayer
	stuff uint8
	Type    EAPOLType
	Length  uint16
}

