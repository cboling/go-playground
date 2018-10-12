package GoParseYourself

import (
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

// TODO: Will focus on creation of an OMCI Frame decoder  (separate one for encode).

// TODO: Look into the go 'encoding' interfaces  -> import "encoding" and derive from BinaryMarshaller
//       and BinaryUnmarshaler

// TODO: Leverage design patterns from the 'gopacket:  https://godoc.org/github.com/google/gopacket#pkg-examples' library

// TODO: Take a look also at: https://github.com/ghedo/go.pkt

// TODO: Read over Brad Israel's blog article:  http://www.bisrael8191.com/Go-Packet-Sniffer/

func main() {
	packet := gopacket.NewPacket(myPacketData, layers.LayerTypeEthernet, gopacket.Default)
	// Get the TCP layer from this packet
	if tcpLayer := packet.Layer(layers.LayerTypeTCP); tcpLayer != nil {
		fmt.Println("This is a TCP packet!")
		// Get actual TCP data from this layer
		tcp, _ := tcpLayer.(*layers.TCP)
		fmt.Printf("From src port %d to dst port %d\n", tcp.SrcPort, tcp.DstPort)
	}
	// Iterate over all layers, printing out each layer type
	for _, layer := range packet.Layers() {
		fmt.Println("PACKET LAYER:", layer.LayerType())
	}
}
