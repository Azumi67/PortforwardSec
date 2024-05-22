package tcp

import (
	"log"
	"net"
	"sync"
)

const (
	bufferSize         = 65535
	maxGoroutines      = 100
	maxErrorLogEntries = 10
)

type ErrorCounter struct {
	counter  int
	maxCount int
	mu       sync.Mutex
}

func (ec *ErrorCounter) Increment() bool {
	ec.mu.Lock()
	defer ec.mu.Unlock()
	ec.counter++
	return ec.counter <= ec.maxCount
}

func forwardTCPPacket(sourceSocket net.Conn, dstSocket net.Conn, errorCounters map[string]*ErrorCounter) {
	for {
		buffer := make([]byte, bufferSize)
		n, err := sourceSocket.Read(buffer)
		if err != nil {
			if err.Error() != "EOF" && !isConnectionResetByPeerError(err) {
				errStr := err.Error()
				if errCounter, ok := errorCounters[errStr]; ok {
					if !errCounter.Increment() {
						continue
					}
				} else {
					ec := &ErrorCounter{counter: 1, maxCount: maxErrorLogEntries}
					ec.Increment()
					errorCounters[errStr] = ec
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

func isConnectionResetByPeerError(err error) bool {
	opErr, ok := err.(*net.OpError)
	if ok {
		errStr := opErr.Err.Error()
		if errStr == "read: connection reset by peer" || errStr == "use of closed network connection" {
			return true
		}
	}
	return false
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
		defer func() {
			iranSocket.Close()
			remoteSocket.Close()
			wg.Done()
		}()
		forwardTCPPacket(iranSocket, remoteSocket, errorCounters)
	}()

	go func() {
		defer func() {
			iranSocket.Close()
			remoteSocket.Close()
			wg.Done()
		}()
		forwardTCPPacket(remoteSocket, iranSocket, errorCounters)
	}()

	wg.Wait()
}

func PortForwardTCP(localHost string, localPort string, remoteHost string, remotePort string) {
	localAddr := net.JoinHostPort(localHost, localPort)

	tcpServerSocket, err := net.Listen("tcp", localAddr)
	if err != nil {
		log.Println("Error occurred while listening for TCP:", err)
		return
	}
	defer tcpServerSocket.Close()

	log.Printf("[*] Azumi is Listening TCP on %s:%s\n", localHost, localPort)

	var wg sync.WaitGroup
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

		wg.Add(1)
		goroutinePool <- struct{}{}
		go func() {
			defer func() {
				<-goroutinePool
				wg.Done()
			}()
			handleTCPIran(iranSocket, remoteHost, remotePort, errorCounters)
		}()
	}
	wg.Wait()
}
