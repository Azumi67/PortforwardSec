package test

import (
	"bufio"
	"log"
	"net"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"sync"

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

func PortFwdUDP(localPt string, remoteKharej string, remotePort string, command string, buffer int, enc reedsolomon.Encoder) {
	log.Printf("Azumichan is starting port forwarding: %s -> %s:%s\n", localPt, remoteKharej, remotePort)

	cmd := exec.Command("socat",
		"UDP6-LISTEN:"+localPt+",reuseaddr,fork",
		"UDP6:["+remoteKharej+"]:"+remotePort)

	if command != "" {
		cmd = exec.Command("sh", "-c", command)
	}

	if buffer > 0 {
		cmd.Env = append(os.Environ(), "SOCAT_SNDBUF="+strconv.Itoa(buffer))
	}

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
func PortFwdUDP(localPt string, remoteKharej string, remotePort string, command string, buffer int, enc reedsolomon.Encoder) {
	log.Printf("Azumichan is starting port forwarding: %s -> %s:%s\n", localPt, remoteKharej, remotePort)

	maxGoro(remoteIP, remotePort)

	cmd := exec.Command("socat",
		"UDP6-LISTEN:"+localPt+",reuseaddr,fork",
		"UDP6:["+remoteKharej+"]:"+remotePort)

	if command != "" {
		cmd = exec.Command("sh", "-c", command)
	}

	if buffer > 0 {
		cmd.Env = append(os.Environ(), "SOCAT_SNDBUF="+strconv.Itoa(buffer))
	}

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

	log.Printf("Sending traffic to %s:%s\n", remoteIP, remotePort)

	err = cmd.Wait()
	if err != nil {
		log.Println("Got hit by a wall while running port forwarding:", err)
	}

	log.Println("Port forwarding just stopped.")
}
func maxGoro(remoteIP string, remotePort string) {
	max := runtime.GOMAXPROCS(0) 
	wg := sync.WaitGroup{}
	poison := make(chan struct{}, max)
	poolConn := make(chan net.Conn, max) 

	for i := 0; i < max; i++ {
		wg.Add(1)
		poison <- struct{}{}
		go func() {
			defer func() {
				wg.Done()
				<-poison
			}()

			conn := toPool(poolConn, remoteIP, remotePort)

			returnPool(conn, poolConn)
		}()
	}

	wg.Wait()
	close(poison)
	close(poolConn)
}

func toPool(pool chan net.Conn, remoteIP string, remotePort string) net.Conn {
	select {
	case conn := <-pool:
		return conn
	default:
		conn, err := net.Dial("udp", remoteIP+":"+remotePort)
		if err != nil {
			log.Println("Couldn't establish a connection:", err)
			return nil
		}
		return conn
	}
}

func returnPool(conn net.Conn, pool chan net.Conn) {
	select {
	case pool <- conn:
	default:
		conn.Close()
	}
}
