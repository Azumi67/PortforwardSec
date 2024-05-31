package udp

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"strconv"
	"sync"

	"github.com/klauspost/reedsolomon"
)

func installSct() error {
	_, err := exec.LookPath("socat")
	if err == nil {
		log.Println("socat is already installed")
		return nil
	}

	log.Println("Installing socat...")

	cmd := exec.Command("sudo", "apt-get", "install", "socat", "-y")

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("couldn't install socat: %v", err)
	}

	log.Println("Installation completed successfully")

	return nil
}

func PortFwdUDP(localPort string, remoteIP string, remotePort string, command string, buffer int, enc reedsolomon.Encoder) {
	log.Printf("Azumichan is starting port forwarding (UDP4): %s -> %s:%s\n", localPort, remoteIP, remotePort)

	cmd := exec.Command("socat",
		"UDP4-LISTEN:"+localPort+",reuseaddr,fork",
		"UDP4:"+remoteIP+":"+remotePort)

	if command != "" {
		cmd = exec.Command("sh", "-c", command)
	}

	if buffer > 0 {
		cmd.Env = append(os.Environ(), "SOCAT_SNDBUF="+strconv.Itoa(buffer))
	}

	startFwding(cmd)
}

func PortFwdUDP6(localPort string, remoteIP string, remotePort string, command string, buffer int, enc reedsolomon.Encoder) {
	log.Printf("Azumichan is starting port forwarding (UDP6): %s -> %s:%s\n", localPort, remoteIP, remotePort)

	cmd := exec.Command("socat",
		"UDP6-LISTEN:"+localPort+",reuseaddr,fork",
		"UDP6:["+remoteIP+"]:"+remotePort)

	if command != "" {
		cmd = exec.Command("sh", "-c", command)
	}

	if buffer > 0 {
		cmd.Env = append(os.Environ(), "SOCAT_SNDBUF="+strconv.Itoa(buffer))
	}

	startFwding(cmd)
}

func startFwding(cmd *exec.Cmd) {
	var wg sync.WaitGroup
	wg.Add(2) 

	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		log.Println("Some weird error occurred while constructing water pipe:", err)
		return
	}
	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		log.Println("Some weird error occurred with water pipe:", err)
		return
	}

	err = cmd.Start()
	if err != nil {
		log.Println("Something weird happened while starting port forwarding:", err)
		return
	}

	scanner := bufio.NewScanner(stdoutPipe)
	go func() {
		defer wg.Done() 
		for scanner.Scan() {
			log.Println(scanner.Text())
		}
	}()

	errScanner := bufio.NewScanner(stderrPipe)
	go func() {
		defer wg.Done() 
		for errScanner.Scan() {
			log.Println(errScanner.Text())
		}
	}()

	wg.Wait() 

	err = cmd.Wait()
	if err != nil {
		log.Println("Got hit by a wall while running port forwarding:", err)
	}
}
