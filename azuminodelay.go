package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/Azumi67/PortforwardSec/nodelay"
)

func main() {
	if len(os.Args) != 8 {
		fmt.Printf("Usage: %s iranIP iranPort remoteIP remotePort protocol [tcpnodelay] true or false buffersize\n", os.Args[0])
		return
	}

	iranIP := os.Args[1]
	localPort := os.Args[2]
	remoteIP := os.Args[3]
	remotePort := os.Args[4]
	protocol := os.Args[5]
	noNagle, err := strconv.ParseBool(os.Args[6])
	if err != nil {
		fmt.Println("Invalid input for tcpnodelay. Enter true or false.")
		return
	}
	bufferSize, err := strconv.Atoi(os.Args[7])
	if err != nil {
		fmt.Println("Invalid input for buffersize. example: 65535.")
		return
	}

	switch protocol {
	case "tcp":
		nodelay.PortForwardTCP(iranIP, localPort, remoteIP, remotePort, noNagle, bufferSize)
	default:
		fmt.Println("Invalid protocol. Supported protocols are tcp and udp. More methods coming!")
	}
}
