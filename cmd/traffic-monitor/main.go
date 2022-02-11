package main

import (
	"log"

	trafficmonitor "github.com/aaltgod/gokyrie/internal/traffic-monitor"
)

func main() {

	pcapWrapper := trafficmonitor.NewPcapWrapper()
	if err := pcapWrapper.CapturePackets("eth0"); err != nil {
		log.Fatal(err)
	}
}
