package main

import (
	"fmt"
	"log"

	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
)

func main() {

	devices, err := pcap.FindAllDevs()
	if err != nil {
		log.Fatal("DEVICE", err)
	}

	for _, device := range devices {
		fmt.Println(device.Name, device.Description, device.Description)
	}
	if handle, err := pcap.OpenLive("eth0", 1600, true, pcap.BlockForever); err != nil {
		log.Fatal(err)
	} else {
		packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
		for packet := range packetSource.Packets() {
			log.Println(packet)
		}
	}
}
