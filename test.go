package main

import (
	"flag"
	"log"

	"github.com/Azumi67/PortforwardSec/udp.lite"
	"github.com/klauspost/reedsolomon"
)

func main() {
	install := flag.Bool("install", false, "Install socat")
	var iranPort, remoteIP, remotePort, command string
	var bufferSize int

	flag.StringVar(&iranPort, "iranPort", "", "Iran port")
	flag.StringVar(&remoteIP, "remoteIP", "", "Remote IP")
	flag.StringVar(&remotePort, "remotePort", "", "Remote port")
	flag.StringVar(&command, "command", "", "Custom command to use instead of socat")
	flag.IntVar(&bufferSize, "bufferSize", 0, "Buffer size in bytes")

	flag.Parse()

	if *install {
		err := udplite.installSct()
		if err != nil {
			log.Fatal(err)
		}
	}

	if iranPort == "" || remoteIP == "" || remotePort == "" {
		log.Fatal("Plz provide correct inputs for iranPort, remoteKharej, and remotePort")
	}

	enc, err := reedsolomon.New(2, 2)
	if err != nil {
		log.Println("Error creating ReedSolomon encoder:", err)
		return
	}

	udplite.PortFwdUDP(iranPort, remoteIP, remotePort, command, bufferSize, enc)
}
