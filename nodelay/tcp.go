package nodelay

import (
	
	"log"
	"net"
	"strings"
	"sync"
)

const (
	goro      = 100
	logEntries = 10
)

type logCounter struct {
	counter  int
	maxCounter int
	mu       sync.Mutex
}

func (ec *logCounter) Increment() bool {
	ec.mu.Lock()
	defer ec.mu.Unlock()
	ec.counter++
	return ec.counter <= ec.maxCounter
}

func forwardTCPNagle(iranSocket net.Conn, dstSocket net.Conn, errorCounters map[string]*logCounter, bufferSize int) {
	buffer := make([]byte, bufferSize)
	for {
		n, err := iranSocket.Read(buffer)
		if err != nil {
			if err.Error() != "EOF" && !userError(err) {
				errStr := err.Error()
				if errCounter, ok := errorCounters[errStr]; ok {
					if !errCounter.Increment() {
						continue
					}
				} else {
					ec := &logCounter{counter: 1, maxCounter: logEntries}
					ec.Increment()
					errorCounters[errStr] = ec
				}
			}
			if strings.Contains(err.Error(), "reuse of closed connections") {
				continue
			}
			return
		}
		if n == 0 {
			break
		}
		_, err = dstSocket.Write(buffer[:n])
		if err != nil {
			return
		}
	}
}

func userError(err error) bool {
	opErr, ok := err.(*net.OpError)
	if ok {
		errStr := opErr.Err.Error()
		if errStr == "read: connection reset by user" || errStr == "reuse of closed connections" {
			return true
		}
	}
	return false
}

func handleTCPIran(iranSocket net.Conn, remoteIP string, remotePort string, noNagle bool, errorCounters map[string]*logCounter, bufferSize int) {
	remoteAddr := net.JoinHostPort(remoteIP, remotePort)

	dstSocket, err := net.Dial("tcp", remoteAddr)
	if err != nil {
		log.Println("Error occoured while connecting with TCP proto:", err)
		iranSocket.Close()
		return
	}

	if tcpConn, ok := iranSocket.(*net.TCPConn); ok && noNagle {
		tcpConn.SetNoDelay(true) 
	}

	var wg sync.WaitGroup
	wg.Add(2)

	go func(dstSocket net.Conn) {
		defer func() {
			iranSocket.Close()
			dstSocket.Close()
			wg.Done()
		}()
		forwardTCPNagle(iranSocket, dstSocket, errorCounters, bufferSize)
		log.Println("Finished Forwarding TCP packet from Iran to kharej.")
	}(dstSocket)

	go func(dstSocket net.Conn) {
		defer func() {
			iranSocket.Close()
			dstSocket.Close()
			wg.Done()
		}()
		forwardTCPNagle(dstSocket, iranSocket, errorCounters, bufferSize)
		log.Println("Finished Forwarding TCP packet from kharej to Iran.")
	}(dstSocket)

	wg.Wait()
	log.Println("TCP proto connection closed.It will be used again")
}

func PortForwardTCP(iranIP string, localPort string, remoteIP string, remotePort string, noNagle bool, bufferSize int) {
	localAddr := net.JoinHostPort(iranIP, localPort)

	kharejSocket, err := net.Listen("tcp", localAddr)
	if err != nil {
		log.Println("Error occurred while listening for TCP:", err)
		return
	}
	defer kharejSocket.Close()

	log.Printf("[*] Azumi is Listening TCP on %s:%s\n", iranIP, localPort)

	var wg sync.WaitGroup
	goroPool := make(chan struct{}, goro)
	errorCounters := make(map[string]*logCounter)
	for {
		iranSocket, err := kharejSocket.Accept()
		if err != nil {
			log.Println("Error occurred while accepting TCP connection:", err)
			continue
		}
		iranAddress := iranSocket.RemoteAddr().(*net.TCPAddr)
		log.Printf("[*] Azumi has Accepted TCP connection from %s:%d\n", iranAddress.IP.String(), iranAddress.Port)

		wg.Add(1)
		goroPool <- struct{}{}
		go func(iranSocket net.Conn) {
			remoteAddr := net.JoinHostPort(remoteIP, remotePort)

			dstSocket, err := net.Dial("tcp", remoteAddr)
			if err != nil{
    log.Println("Error occurred while connecting with TCP Proto:", err)
    iranSocket.Close()
    wg.Done()
    <-goroPool
    return
}

if tcpConn, ok := iranSocket.(*net.TCPConn); ok && noNagle {
    tcpConn.SetNoDelay(true) 
}

defer func() {
    iranSocket.Close()
    dstSocket.Close()
    wg.Done()
    <-goroPool
}()
handleTCPIran(iranSocket, remoteIP, remotePort, noNagle, errorCounters, bufferSize) 
}(iranSocket)
}
wg.Wait()
}
