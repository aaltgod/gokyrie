package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

var (
	device      string        = "eth0"
	snapshotLen int32         = 1024
	promiscuous bool          = false
	timeout     time.Duration = 10 * time.Second
	handle      *pcap.Handle
)

func main() {

	handle, err := pcap.OpenLive(device, snapshotLen, promiscuous, timeout)
	if err != nil {
		log.Fatal(err)
	}
	defer handle.Close()

	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	for packet := range packetSource.Packets() {
		printPacketInfo(packet)
	}
}

func printPacketInfo(packet gopacket.Packet) {

	for _, layer := range packet.Layers() {
		var output string

		layerType := layer.LayerType()
		switch layerType {
		case layers.LayerTypeEthernet:
			ethernetPacket, _ := packet.Layer(layers.LayerTypeEthernet).(*layers.Ethernet)
			output = fmt.Sprintf(
				"Ethernet layer detected\nSource MAC: %s Destination MAC: %s\nEthernet type: %s\n\n",
				ethernetPacket.SrcMAC,
				ethernetPacket.DstMAC,
				ethernetPacket.EthernetType,
			)

		case layers.LayerTypeIPv4:
			ipv4Packet, _ := packet.Layer(layers.LayerTypeIPv4).(*layers.IPv4)
			output = fmt.Sprintf(
				"IPv4 layer detected\nFrom %s to %s\tProtocol: %s\n",
				ipv4Packet.SrcIP,
				ipv4Packet.DstIP,
				ipv4Packet.Protocol,
			)

		case layers.LayerTypeTCP:
			tcpPacket, _ := packet.Layer(layers.LayerTypeTCP).(*layers.TCP)
			output = fmt.Sprintf(
				"TCP layer detected\nFrom %d to %d\t Sequence number: %d\n\n",
				tcpPacket.SrcPort,
				tcpPacket.DstPort,
				tcpPacket.Seq,
			)

		default:
			appLayer := packet.ApplicationLayer()
			if appLayer != nil {
				var isHTTPExists bool

				if strings.Contains(string(appLayer.Payload()), "HTTP") {
					isHTTPExists = true
				}

				fmt.Printf(
					"Application layer/Payload detected\n%s\n%t\n\n",
					appLayer.Payload(),
					isHTTPExists,
				)
			}
		}

		fmt.Println(output)
	}
}
