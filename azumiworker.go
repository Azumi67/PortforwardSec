package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"sync"
)

const (
	bufferSize         = 65535
	maxGoroutines      = 100
	pipeBufferSize     = 8192
	maxErrorLogEntries = 10 
)

type ErrorCounter struct {
	counter  int
	maxCount int
}

func (ec *ErrorCounter) Increment() bool {
	ec.counter++
	return ec.counter <= ec.maxCount
}

func forwardTCPPacket(sourceSocket net.Conn, dstSocket net.Conn, errorCounters map[string]*ErrorCounter) {
	buffer := make([]byte, bufferSize)
	for {
		n, err := sourceSocket.Read(buffer)
		if err != nil {
			if err.Error() != "EOF" {
				if err.(*net.OpError).Err.Error() == "read: connection reset by user" || err.(*net.OpError).Err.Error() == "use of closed network connection" {
					return
				}

				if errCounter, ok := errorCounters[err.Error()]; ok {
					if !errCounter.Increment() {
						continue
					}
				} else {
					errorCounters[err.Error()] = &ErrorCounter{1, maxErrorLogEntries}
				}
				log.Println("Error occurred while reading TCP packet:", err)
			}
			return
		}
		if n == 0 {
			break
		}
		_, err = dstSocket.Write(buffer[:n])
		if err != nil {
			log.Println("Error occurred while writing TCP packet:", err)
			return
		}
	}
}
func handleTCPIran(iranSocket net.Conn, remoteHost string, remotePort string, errorCounters map[string]*ErrorCounter) {
	remoteAddr := net.JoinHostPort(remoteHost, remotePort)

	remoteSocket, err := net.Dial("tcp", remoteAddr)
	if err != nil {
		log.Println("Error occurred while connecting with TCP Proto:", err)
		iranSocket.Close()
		return
	}

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		forwardTCPPacket(iranSocket, remoteSocket, errorCounters)
		iranSocket.Close()
		remoteSocket.Close()
	}()

	go func() {
		defer wg.Done()
		forwardTCPPacket(remoteSocket, iranSocket, errorCounters)
		iranSocket.Close()
		remoteSocket.Close()
	}()

	wg.Wait()
}

func portForwardTCP(localHost string, localPort string, remoteHost string, remotePort string) {
	localAddr := net.JoinHostPort(localHost, localPort)

	tcpServerSocket, err := net.Listen("tcp", localAddr)
	if err != nil {
		log.Println("Error occurred while listening for TCP:", err)
		return
	}
	defer tcpServerSocket.Close()

	log.Printf("[*] Azumi is Listening TCP on %s:%s\n", localHost, localPort)

	goroutinePool := make(chan struct{}, maxGoroutines)
	errorCounters := make(map[string]*ErrorCounter)
	for {
		iranSocket, err := tcpServerSocket.Accept()
		if err != nil {
			log.Println("Error occurred while accepting TCP connection:", err)
			continue
		}
		iranAddress := iranSocket.RemoteAddr().(*net.TCPAddr)
		log.Printf("[*] Azumi has Accepted TCP connection from %s:%d\n", iranAddress.IP.String(), iranAddress.Port)

		goroutinePool <- struct{}{}
		go func() {
			defer func() { <-goroutinePool }()
			handleTCPIran(iranSocket, remoteHost, remotePort, errorCounters)
		}()
	}
}

func forwardUDPPacket(sourceConn *net.UDPConn, dstAddr *net.UDPAddr) {
	buffer := make([]byte, bufferSize)
	for {
		n, _, err := sourceConn.ReadFromUDP(buffer)
		if err != nil {
			log.Println("Error occurred while sending UDP packets:", err)
			return
		}
		_, err = sourceConn.WriteToUDP(buffer[:n], dstAddr)
		if err != nil {
			log.Println("Error occurred while sending UDP packets:", err)
			return
		}
	}
}

func handleUDPIran(iranConn *net.UDPConn, remoteHost string, remotePort string) {
	remoteAddr := net.JoinHostPort(remoteHost, remotePort)
	dstAddr, err := net.ResolveUDPAddr("udp", remoteAddr)
	if err != nil {
		log.Println("Error occurred while trying to resolve UDP address:", err)
		return
	}

	forwardUDPPacket(iranConn, dstAddr)
}

func portForwardUDP(localHost string, localPort string, remoteHost string, remotePort string) {
	localAddr := net.JoinHostPort(localHost, localPort)

	udpServerAddr, err := net.ResolveUDPAddr("udp", localAddr)
	if err != nil {
		log.Println("Error occurred while resolving UDP address:", err)
		return
	}

	udpServerConn, err := net.ListenUDP("udp", udpServerAddr)
	if err != nil {
		log.Println("Error occurred while listening for UDP:", err)
		return
	}
	defer udpServerConn.Close()

	log.Printf("[*] Azumi is Listening UDP on %s:%s\n", localHost, localPort)

	for {
		buffer := make([]byte, bufferSize)
		n, _, err := udpServerConn.ReadFromUDP(buffer)
		if err != nil {
			log.Println("Error occurred while receiving UDP packet:", err)
			continue
		}
		log.Printf("[*] Azumi has Accepted UDP packet\n")

		go handleUDPIran(udpServerConn, remoteHost, remotePort)
		_, err = udpServerConn.WriteToUDP(buffer[:n], udpServerAddr)
		if err != nil {
			log.Println("Error occurred while sending UDP packet:", err)
		}
	}
}

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
		portForwardTCP(localHost, localPort, remoteHost, remotePort)
	case "udp":
		portForwardUDP(localHost, localPort, remoteHost, remotePort)
	default:
		fmt.Println("Invalid protocol. Supported protocols are tcp and udp.More methods coming!")
	}
}
