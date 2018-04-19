package main

import (
	"bufio"
	"log"
	"net"
	"os"
	"sync"
	"time"
)

var (
	outFileName = "/var/log/flowlog.json"
	lAddr       = ":55555"
	lock        sync.Mutex
	writer      *bufio.Writer
)

func flusher() {
	var lastNonEmpty bool
	for {
		lock.Lock()
		if writer.Buffered() > 0 {
			if lastNonEmpty {
				writer.Flush()
			} else {
				lastNonEmpty = true
			}
		} else {
			lastNonEmpty = false
		}
		lock.Unlock()
		time.Sleep(1 * time.Second)
	}
}

func handleConn(c net.Conn) {
	log.Printf("New connection from %s\n", c.RemoteAddr())

	defer c.Close()
	s := bufio.NewScanner(c)

	for s.Scan() {
		lock.Lock()
		_, err := writer.WriteString(s.Text() + "\n")
		lock.Unlock()
		if err != nil {
			log.Printf("Error reading from remote endpoint %s: %v\n", c.RemoteAddr(), err)
			return
		}
	}

	if err := s.Err(); err != nil {
		log.Printf("Error reading from remote endpoint %s: %v\n", c.RemoteAddr(), err)
	}
	log.Printf("Connection from %s closed\n", c.RemoteAddr())
}

func main() {
	if len(os.Args) > 1 {
		outFileName = os.Args[1]
	}
	if len(os.Args) > 2 {
		lAddr = ":" + os.Args[2]
	}
	sock, err := net.Listen("tcp", lAddr)
	if err != nil {
		log.Fatalf("Error opening listening connection to %s: %v", lAddr, err)
	}
	defer sock.Close()
	log.Printf("Listening on %s\n", lAddr)

	of, err := os.Create(outFileName)
	if err != nil {
		log.Fatalf("Error opening %q: %v", outFileName, err)
	}
	writer = bufio.NewWriter(of)
	log.Printf("Output file at %s\n", outFileName)

	go flusher()

	for {
		c, err := sock.Accept()
		if err != nil {
			log.Fatalf("Acceping connection: %v", err)
		}
		go handleConn(c)
	}

}
