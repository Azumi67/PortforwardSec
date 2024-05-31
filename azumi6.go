package main

import (
        "flag"
        "log"
        "os"
        "fmt"
        "os/exec"
        "github.com/Azumi67/PortforwardSec/udp6"
        "github.com/klauspost/reedsolomon"

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
	var iranPort, remoteIP, remotePort, command string
	var bufferSize int

	flag.StringVar(&iranPort, "iranPort", "", "Iran port")
	flag.StringVar(&remoteIP, "remoteIP", "", "Remote IP")
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

	if iranPort == "" || remoteIP == "" || remotePort == "" {
		log.Fatal("Plz provide correct inputs for iranPort, remoteKharej, and remotePort")
	}

	enc, err := reedsolomon.New(2, 2)
	if err != nil {
		log.Println("Error creating ReedSolomon encoder:", err)
		return
	}

	udp6.PortFwdUDP(iranPort, remoteIP, remotePort, command, bufferSize, enc)
}
