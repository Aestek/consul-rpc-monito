package main

import (
	"flag"
	"log"
	"os"

	"github.com/hashicorp/go-bexpr"
)

func main() {
	iface := flag.String("iface", "lo", "Interface to capture")
	port := flag.Int("port", 8300, "Port to capture")
	debug := flag.Bool("debug", false, "Debug")
	filter := flag.String("filter", "", "Filter using bexpr")
	flag.Parse()

	log.SetOutput(os.Stdout)

	var filterFn *bexpr.Evaluator
	if *filter != "" {
		eval, err := bexpr.CreateEvaluatorForType(*filter, nil, Msg{})
		if err != nil {
			log.Fatal(err)
		}
		filterFn = eval
	}

	packetSource := source(*iface, *port)
	for packet := range packetSource.Packets() {
		msg, err := decode(packet)
		if err != nil {
			if *debug {
				log.Printf("error: %s", err)
			}
			continue
		}

		if filterFn != nil {
			ok, err := filterFn.Evaluate(msg)
			if err != nil {
				log.Println(err)
			}

			if !ok {
				continue
			}
		}

		log.Printf("%+v", msg)
	}
}
