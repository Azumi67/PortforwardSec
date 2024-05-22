package main

import (
	"fmt"
	"os"
	"github.com/Azumi67/PortforwardSEC/tcp"
)

func main() {
	if len(os.Args) != 6 {
		fmt.Printf("Usage: %s local_host local_port remote_host remote_port protocol\n", os.Args[0])
		return
	}

	localHost := os.Args[1]
	localPort := os.Args[2]
	remoteHost := os.Args[3]
	remotePort := os.Args[4]
	protocol := os.Args[5]

	switch protocol {
	case "tcp":
		tcp.PortForwardTCP(localHost, localPort, remoteHost, remotePort)
	default:
		fmt.Println("Invalid protocol. Supported protocols are tcp and udp. More methods coming!")
	}
}
