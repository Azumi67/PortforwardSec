package udp

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

func forwardUDPPacket(sourceConn *net.UDPConn, dstConn *net.UDPConn, errorCounters map[string]*ErrorCounter) {
	for {
		buffer := make([]byte, bufferSize)
		n, addr, err := sourceConn.ReadFromUDP(buffer)
		if err != nil {
			if err.Error() != "EOF" && !PeerError(err) {
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
				log.Println("Error occurred while reading UDP packet:", err)
			}
			return
		}
		if n == 0 {
			break
		}
		_, err = dstConn.WriteToUDP(buffer[:n], addr)
		if err != nil {
			log.Println("Error occurred while writing UDP packet:", err)
			return
		}
	}
}

func PeerError(err error) bool {
	opErr, ok := err.(*net.OpError)
	if ok {
		errStr := opErr.Err.Error()
		if errStr == "read: connection reset by peer" || errStr == "use of closed network connection" {
			return true
		}
	}
	return false
}

func handleUDPIran(iranConn *net.UDPConn, remoteHost string, remotePort string, errorCounters map[string]*ErrorCounter) {
	remoteAddr, err := net.ResolveUDPAddr("udp", net.JoinHostPort(remoteHost, remotePort))
	if err != nil {
		log.Println("Error occurred while resolving UDP remote address:", err)
		iranConn.Close()
		return
	}

	remoteConn, err := net.DialUDP("udp", nil, remoteAddr)
	if err != nil {
		log.Println("Error occurred while connecting with UDP remote host:", err)
		iranConn.Close()
		return
	}

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer func() {
			iranConn.Close()
			remoteConn.Close()
			wg.Done()
		}()
		forwardUDPPacket(iranConn, remoteConn, errorCounters)
		log.Println("Finished forwarding UDP packet from Iran to remote host.")
	}()

	go func() {
		defer func() {
			iranConn.Close()
			remoteConn.Close()
			wg.Done()
		}()
		forwardUDPPacket(remoteConn, iranConn, errorCounters)
		log.Println("Finished forwarding UDP packet from remote host to Iran.")
	}()

	wg.Wait()
	log.Println("UDP connection closed.")
}

func PortForwardUDP(localHost string, localPort string, remoteHost string, remotePort string) {
	localAddr, err := net.ResolveUDPAddr("udp", net.JoinHostPort(localHost, localPort))
	if err != nil {
		log.Println("Error occurred while resolving UDP local address:", err)
		return
	}

	udpServerConn, err := net.ListenUDP("udp", localAddr)
	if err != nil {
		log.Println("Error occurred while listening for UDP:", err)
		return
	}
	defer udpServerConn.Close()

	log.Printf("[*] Azumi is listening UDP on %s:%s\n", localHost, localPort)

	var wg sync.WaitGroup
	goroutinePool := make(chan struct{}, maxGoroutines)
	errorCounters := make(map[string]*ErrorCounter)

	for {
		buffer := make([]byte, bufferSize)
		n, addr, err := udpServerConn.ReadFromUDP(buffer)
		if err != nil {
			log.Println("Error occurred while reading UDP connection:", err)
			continue
		}
		log.Printf("[*] Azumi has accepted UDP connection from %s:%d\n", addr.IP.String(), addr.Port)

		iranConn, err := net.DialUDP("udp", nil, &net.UDPAddr{IP: net.IPv4(0, 0, 0, 0), Port: 0})
		if err != nil {
			log.Println("Error occurred while connecting with UDPIran:", err)
			continue
		}

		wg.Add(1)
		goroutinePool <- struct{}{}
		go func() {
			defer func() {
				<-goroutinePool
				wg.Done()
			}()
			handleUDPIran(iranConn, remoteHost, remotePort, errorCounters)
		}()
	}

	wg.Wait()
}
