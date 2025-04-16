package main

import (
	"log"
	"math/rand"
	"net"
	"time"
)

const (
	addr          = "0.0.0.0:12345"
	proto         = "tcp4"
	proverbPeriod = 3 * time.Second
)

type ProverbServer struct {
	proverbs []string
	rnd      *rand.Rand
}

func NewProverbServer() *ProverbServer {
	proverbs := []string{
		"Don't communicate by sharing memory, share memory by communicating.",
		"Concurrency is not parallelism.",
		"Channels orchestrate; mutexes serialize.",
		"The bigger the interface, the weaker the abstraction.",
		"Make the zero value useful.",
		"interface{} says nothing.",
		"Gofmt's style is no one's favorite, yet gofmt is everyone's favorite.",
		"A little copying is better than a little dependency.",
		"Syscall must always be guarded with build tags.",
		"Cgo must always be guarded with build tags.",
		"Cgo is not Go.",
		"With the unsafe package there are no guarantees.",
		"Clear is better than clever.",
		"Reflection is never clear.",
		"Errors are values.",
		"Don't just check errors, handle them gracefully.",
		"Design the architecture, name the components, document the details.",
		"Documentation is for users.",
		"Don't panic.",
	}

	return &ProverbServer{
		proverbs: proverbs,
		rnd:      rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (ps *ProverbServer) getRandomProverb() string {
	return ps.proverbs[ps.rnd.Intn(len(ps.proverbs))]
}

func (ps *ProverbServer) handleConn(conn net.Conn) {
	defer conn.Close()

	header := "\033[2J\033[H" +
		"\033[1;1H=== Go Proverbs Server ===" +
		"\033[2;1H(press 'q' to quit)" +
		"\033[4;1HNew proverb every 3 seconds:" +
		"\033[5;1H---------------------------" +
		"\033[7;1H"

	_, err := conn.Write([]byte(header))
	if err != nil {
		log.Printf("Error writing header: %v", err)
		return
	}

	quit := make(chan struct{})

	go func() {
		buf := make([]byte, 1)
		for {
			_, err := conn.Read(buf)
			if err != nil || buf[0] == 'q' {
				close(quit)
				return
			}
		}
	}()

	ticker := time.NewTicker(proverbPeriod)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			proverb := ps.getRandomProverb()
			_, err := conn.Write([]byte("\r\033[K" + proverb + "\n"))
			if err != nil {
				log.Printf("Client disconnected: %v", conn.RemoteAddr())
				return
			}
		case <-quit:
			conn.Write([]byte("\n\033[1;1HConnection closed\n"))
			return
		}
	}
}

func main() {
	server := NewProverbServer()

	listener, err := net.Listen(proto, addr)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
	defer listener.Close()

	log.Printf("Go Proverb Server is listening on %s (%s)", addr, proto)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Error accepting connection: %v", err)
			continue
		}

		log.Printf("New client connected: %v", conn.RemoteAddr())
		go server.handleConn(conn)
	}
}
