package main

import (
	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
)

func source(iface string, port int) *gopacket.PacketSource {
	handle, err := pcap.OpenLive(iface, 1600, true, pcap.BlockForever)
	if err != nil {
		panic(err)
	}

	// bpf generated with sudo tcpdump -i lo -dd "ip and tcp and port 8300"
	bpfInstructions := []pcap.BPFInstruction{
		{0x28, 0, 0, 0x0000000c},
		{0x15, 0, 10, 0x00000800},
		{0x30, 0, 0, 0x00000017},
		{0x15, 0, 8, 0x00000006},
		{0x28, 0, 0, 0x00000014},
		{0x45, 6, 0, 0x00001fff},
		{0xb1, 0, 0, 0x0000000e},
		{0x48, 0, 0, 0x0000000e},
		{0x15, 2, 0, 0x0000206c},
		{0x48, 0, 0, 0x00000010},
		{0x15, 0, 1, uint32(port)},
		{0x6, 0, 0, 0x00040000},
		{0x6, 0, 0, 0x00000000},
	}

	if err := handle.SetBPFInstructionFilter(bpfInstructions); err != nil {
		panic(err)
	}

	return gopacket.NewPacketSource(handle, handle.LinkType())
}
