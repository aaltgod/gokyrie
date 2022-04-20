package trafficmonitor

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/aaltgod/gokyrie/internal/config"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"github.com/google/gopacket/pcapgo"
)

const (
	// tcp length
	defaultSnapLen = 262144
)

var (
	snapshotLen int32 = 1024
	promiscuous bool  = true
)

type Data struct {
	IP   string
	Text string
}

type Sender struct {
	IP         string
	PageAmount int
}

type PcapWrapper struct {
	config        *config.Config
	wg            *sync.WaitGroup
	mu            *sync.Mutex
	dataCh        chan Data
	SendersByIPCh chan map[string]*Sender
	ErrorCh       chan error
}

func NewPcapWrapper(cfg *config.Config, dataCh chan Data) *PcapWrapper {
	return &PcapWrapper{
		config:        cfg,
		wg:            &sync.WaitGroup{},
		mu:            &sync.Mutex{},
		dataCh:        dataCh,
		SendersByIPCh: make(chan map[string]*Sender, len(cfg.Teams)),
		ErrorCh:       make(chan error),
	}
}

func (p *PcapWrapper) StartListeners(ctx context.Context) error {

	go p.statistic(ctx)

	p.wg.Add(len(p.config.Services))
	for _, service := range p.config.Services {
		go p.listener(ctx, service)
	}
	p.wg.Wait()

	for {
		select {
		case err := <-p.ErrorCh:
			return err
		default:
			return nil
		}
	}
}

func (p *PcapWrapper) listener(ctx context.Context, service config.Service) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			p.capturePackets(ctx, service)
		}
	}
}

func (p *PcapWrapper) statistic(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case err := <-p.ErrorCh:
			// TODO: return error to tui maybe
			log.Fatal("Listener error: ", err)
		case <-p.SendersByIPCh:
			// log.Println(sender)
		}
	}
}

func (p *PcapWrapper) capturePackets(ctx context.Context, service config.Service) {
	defer p.wg.Done()

	file, err := os.Create(service.Name + ".pcap")
	if err != nil {
		p.ErrorCh <- err
		return
	}
	defer file.Close()

	w := pcapgo.NewWriter(file)
	w.WriteFileHeader(defaultSnapLen, layers.LinkTypeEthernet)

	handle, err := pcap.OpenLive(
		p.config.InterfaceName, snapshotLen,
		promiscuous, pcap.BlockForever,
	)
	if err != nil {
		p.ErrorCh <- err
		return
	}
	// FIXME: Close() doesn't return
	//defer handle.Close()

	if err = handle.SetBPFFilter(
		fmt.Sprintf("tcp and port %d", service.Port),
	); err != nil {
		p.ErrorCh <- err
		return
	}

	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	for {
		packetAmount := len(packetSource.Packets())
		select {
		case <-ctx.Done():
			return
		default:
			if packetAmount > 0 {
				for packet := range packetSource.Packets() {
					select {
					case <-ctx.Done():
						return
					default:
						p.handlePacket(w, packet)
					}
				}
			}
		}
	}
}

func (p *PcapWrapper) handlePacket(w *pcapgo.Writer, packet gopacket.Packet) {

	for _, layer := range packet.Layers() {
		layerType := layer.LayerType()
		switch layerType {
		case layers.LayerTypeEthernet:
		case layers.LayerTypeIPv4:
			ipv4Packet, _ := packet.Layer(layers.LayerTypeIPv4).(*layers.IPv4)

			ip := ipv4Packet.SrcIP.String()
			if p.existTeamIP(ip) {
				continue
			}

			sender := make(map[string]*Sender)
			sender[ip] = &Sender{
				IP:         ip,
				PageAmount: 1,
			}
			p.SendersByIPCh <- sender
		case layers.LayerTypeTCP:
			p.mu.Lock()
			w.WritePacket(packet.Metadata().CaptureInfo, packet.Data())
			p.mu.Unlock()
		default:
			appLayer := packet.ApplicationLayer()
			if appLayer != nil {
				var isHTTPExists bool

				if strings.Contains(string(appLayer.Payload()), "HTTP") {
					isHTTPExists = true
				}

				var data Data
				data.IP = appLayer.LayerType().String()
				output := fmt.Sprintf(
					"Application layer/Payload detected\n%s\n%t\n",
					appLayer.Payload(),
					isHTTPExists,
				)
				data.Text = output
				p.dataCh <- data
			}
		}
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
