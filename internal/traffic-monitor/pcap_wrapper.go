package trafficmonitor

import (
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/aaltgod/gokyrie/internal/config"
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
	config        *config.Config
	wg            *sync.WaitGroup
	mu            sync.Mutex
	SendersByIPCh chan map[string]*Sender
	ErrorCh       chan error
}

func NewPcapWrapper(cfg *config.Config) *PcapWrapper {
	return &PcapWrapper{
		config:        cfg,
		wg:            &sync.WaitGroup{},
		mu:            sync.Mutex{},
		SendersByIPCh: make(chan map[string]*Sender, len(cfg.Teams)),
		ErrorCh:       make(chan error),
	}
}

func (p *PcapWrapper) StartListeners() {

	go p.statistic()

	p.wg.Add(len(p.config.Services))
	for _, service := range p.config.Services {
		go p.listener(service)
	}
	p.wg.Wait()
}

func (p *PcapWrapper) listener(service config.Service) {
	p.capturePackets(service)
}

func (p *PcapWrapper) statistic() {
	for {
		select {
		case err := <-p.ErrorCh:
			// TODO: return error to tui maybe
			log.Fatal("Listener error: ", err)
		case sender := <-p.SendersByIPCh:
			log.Println(sender)
		}
	}
}

func (p *PcapWrapper) capturePackets(service config.Service) {
	defer p.wg.Done()

	if handle, err := pcap.OpenLive(
		p.config.InterfaceName, snapshotLen,
		promiscuous, pcap.BlockForever,
	); err != nil {
		p.ErrorCh <- err
		return
	} else if err := handle.SetBPFFilter(
		fmt.Sprintf("tcp and port %d", service.Port),
	); err != nil {
		p.ErrorCh <- err
		return
	} else {
		packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
		for packet := range packetSource.Packets() {
			p.handlePacket(packet)
		}
	}
	defer handle.Close()
}

func (p *PcapWrapper) handlePacket(packet gopacket.Packet) {

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

			ip := ipv4Packet.SrcIP.String()
			if !p.existTeamIP(ip) {
				continue
			}

			sender := make(map[string]*Sender)
			sender[ip] = &Sender{
				IP:         ip,
				PageAmount: 1,
			}
			p.SendersByIPCh <- sender

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

func (p *PcapWrapper) existTeamIP(ip string) bool {
	for _, team := range p.config.Teams {
		if team.IP == ip {
			return true
		}
	}
	return false
}
