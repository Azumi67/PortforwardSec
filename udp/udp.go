package udp

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"

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

func PortFwdUDP(iranPort string, remoteIP string, remotePort string, command string, buffer int, enc reedsolomon.Encoder) {
	log.Printf("Azumichan is starting port forwarding (UDP4): %s -> %s:%s\n", iranPort, remoteIP, remotePort)

	cmd := exec.Command("socat",
		"UDP4-LISTEN:"+iranPort+",reuseaddr,fork",
		"UDP4:"+remoteIP+":"+remotePort)

	if command != "" {
		cmd = exec.Command("sh", "-c", command)
	}

	if buffer > 0 {
		cmd.Env = append(os.Environ(), "SOCAT_SNDBUF="+strconv.Itoa(buffer))
	}

	startFwding(cmd)
}

func PortFwdUDP6(localPt string, remoteKharej string, remotePort string, command string, buffer int, enc reedsolomon.Encoder) {
	log.Printf("Azumichan is starting port forwarding (UDP6): %s -> %s:%s\n", localPt, remoteKharej, remotePort)

	cmd := exec.Command("socat",
		"UDP6-LISTEN:"+localPt+",reuseaddr,fork",
		"UDP6:["+remoteKharej+"]:"+remotePort)

	if command != "" {
		cmd = exec.Command("sh", "-c", command)
	}

	if buffer > 0 {
		cmd.Env = append(os.Environ(), "SOCAT_SNDBUF="+strconv.Itoa(buffer))
	}

	startFwding(cmd)
}

func startFwding(cmd *exec.Cmd) {
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
		for scanner.Scan() {
			log.Println(scanner.Text())
		}
	}()

	errScanner := bufio.NewScanner(stderrPipe)
	go func() {
		for errScanner.Scan() {
			log.Println(errScanner.Text())
		}
	}()

	err = cmd.Wait()
	if err != nil {
		log.Println("Got hit by a wall while running port forwarding:", err)
	}
}
