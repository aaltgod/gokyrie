package trafficmonitor

import (
	"fmt"
	"log"
	"strings"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

const (
	// tcp length
	defaultSnapLen = 262144
)

var (
	snapshotLen int32 = 1024
	promiscuous bool  = true
	handle      *pcap.Handle
)

type Sender struct {
	IP         string
	PageAmount int
}

type PcapWrapper struct {
	SenderByIP map[string]*Sender
}

func NewPcapWrapper() *PcapWrapper {
	return &PcapWrapper{
		SenderByIP: make(map[string]*Sender),
	}
}

func (p *PcapWrapper) CapturePackets(interfaceName string) error {

	if handle, err := pcap.OpenLive(interfaceName, snapshotLen, promiscuous, pcap.BlockForever); err != nil {
		log.Fatal(err)
	} else if err := handle.SetBPFFilter("tcp and port 8081"); err != nil {
		log.Fatal(err)
	} else {
		packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
		for packet := range packetSource.Packets() {
			p.HandlePacket(packet)
			printSenders(p.SenderByIP)
		}
	}
	defer handle.Close()

	return nil
}

func (p *PcapWrapper) HandlePacket(packet gopacket.Packet) {
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

			sender, ok := p.SenderByIP[ipv4Packet.SrcIP.String()]
			if !ok {
				p.SenderByIP[ipv4Packet.SrcIP.String()] = &Sender{
					IP:         ipv4Packet.SrcIP.String(),
					PageAmount: 1,
				}
			} else {
				sender.PageAmount++
			}

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

func printSenders(senders map[string]*Sender) {
	for ip, sender := range senders {
		fmt.Println(ip, sender.IP, sender.PageAmount)
	}
}
