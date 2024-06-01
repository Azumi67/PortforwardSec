package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"
        "github.com/Azumi67/PortforwardSec/udp4"
	"github.com/klauspost/reedsolomon"
)

func screenclean() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func anime() {
	screenclean()

	boxWidth := 10
	fullbar := "█"
	emptybox := "-"

	for i := 0; i <= boxWidth; i++ {
		bar := ""
		for j := 0; j < i; j++ {
			bar += fullbar
		}
		for j := i; j < boxWidth; j++ {
			bar += emptybox
		}

		screenclean()
		fmt.Printf("[%s] %d%%\n", bar, i*10)
		time.Sleep(500 * time.Millisecond)
	}
}

func installSct() error {
	_, err := exec.LookPath("socat")
	if err == nil {
		log.Println("socatudp is already installed")
		return nil
	}

	log.Println("Installing udp requirements...")
	go anime()

	cmd := exec.Command("sudo", "DEBIAN_FRONTEND=noninteractive", "apt-get", "install", "-y", "socat")

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("couldn't create installations for stdout: %v", err)
	}
	err = cmd.Start()
	if err != nil {
		return fmt.Errorf("couldn't start da command: %v", err)
	}

	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Println(line)

	}

	err = cmd.Wait()
	if err != nil {
		return fmt.Errorf("couldn't install udp requirements: %v", err)
	}

	log.Println("socatudp installed successfully")

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
	} else {
		if iranPort == "" || remoteIP == "" || remotePort == "" {
			log.Fatal("Plz provide correct inputs for iranPort, remoteIP, and remotePort")
		}

		enc, err := reedsolomon.New(2, 2)
		if err != nil {
			log.Println("Error creating ReedSolomon encoding:", err)
			return
		}

		maxGoro(remoteIP, remotePort)
		PortFwdUDP(iranPort, remoteIP, remotePort, command, bufferSize, enc)
	}
}
