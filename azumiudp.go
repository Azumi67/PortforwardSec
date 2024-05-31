package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"

	"github.com/klauspost/reedsolomon"
  "github.com/Azumi67/PortforwardSec/udp2"
	flag "github.com/spf13/pflag"
)

func installSct() error {
	_, err := exec.LookPath("socat")
	if err == nil {
		log.Println("socatudp is already installed")
		return nil
	}

	log.Println("Installing udp..")

	cmd := exec.Command("sudo", "apt-get", "install", "socat", "-y")

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("couldn't install socatudp: %v", err)
	}

	log.Println("It installed successfully")

	return nil
}

func main() {
	install := flag.Bool("install", false, "Install socat")
	var localPt, remoteKharej, remotePort, command string
	var bufferSize int

	flag.StringVar(&localPt, "iranPort", "", "Local port")
	flag.StringVar(&remoteKharej, "remoteIP", "", "Remote host")
	flag.StringVar(&remotePort, "remotePort", "", "Remote port")
	flag.StringVar(&command, "command", "", "Custom command to use instead of socat")
	flag.IntVar(&bufferSize, "bufferSize", 0, "Buffer size in bytes")

	flag.Parse()

	if *install {
		err := installSct()
		if err != nil {
			log.Fatal(err)
		}
	}

	if localPt == "" || remoteKharej == "" || remotePort == "" {
		log.Fatal("Please provide correct inputs for iranPort, remoteKharej, and remotePort")
	}

	enc, err := reedsolomon.New(2, 2)
	if err != nil {
		log.Println("Error creating ReedSolomon encoder:", err)
		return
	}

	if remoteKharej == "::1" || remoteKharej == "localhost" {
		udp4.PortFwdUDP(localPt, remoteKharej, remotePort, command, bufferSize, enc)
	} else {
		udp6.PortFwdUDP(localPt, remoteKharej, remotePort, command, bufferSize, enc)
	}
}
