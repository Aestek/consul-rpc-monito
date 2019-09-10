package main

import (
	"bytes"
	"errors"
	"net/rpc"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/hashicorp/go-msgpack/codec"
)

type Msg struct {
	SrcIP  string
	DstIP  string
	Method string
	Body   map[string]interface{}
}

func decode(packet gopacket.Packet) (Msg, error) {
	msg := Msg{}

	alayer := packet.ApplicationLayer()
	if alayer == nil {
		return msg, errors.New("no application layer")
	}

	payload := alayer.Payload()

	buff := bytes.NewBuffer(payload)
	dec := codec.NewDecoder(buff, &codec.MsgpackHandle{
		RawToString: true,
	})

	req := &rpc.Request{}
	err := dec.Decode(req)
	if err != nil {
		return msg, err
	}

	msg.Method = req.ServiceMethod

	body := map[string]interface{}{}
	err = dec.Decode(&body)
	if err != nil {
		return msg, err
	}

	msg.Body = body

	if ipLayer := packet.Layer(layers.LayerTypeIPv4); ipLayer != nil {
		ip, _ := ipLayer.(*layers.IPv4)
		msg.SrcIP = ip.SrcIP.String()
		msg.DstIP = ip.DstIP.String()
	}
	if ipLayer := packet.Layer(layers.LayerTypeIPv6); ipLayer != nil {
		ip, _ := ipLayer.(*layers.IPv6)
		msg.SrcIP = ip.SrcIP.String()
		msg.DstIP = ip.DstIP.String()
	}

	return msg, nil
}
