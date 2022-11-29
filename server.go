package main

import (
	"fmt"
	"log"
	"net"

	"strings"
	"sync"
)

// establish UDP server
// receive 2 requests
// exchange messages
// quit not affect A and B commuication
type peers struct {
	mu   sync.RWMutex
	addr map[string]*net.UDPAddr
	ch   chan struct {
		string
		*net.UDPAddr
	}
}

func main() {
	fmt.Println("begin server")

	//peers := make([]*net.UDPAddr,2)
	p := peers{
		addr: make(map[string]*net.UDPAddr),
	}
	listener, err := net.ListenUDP("udp", &net.UDPAddr{Port: 8081})
	if err != nil {
		log.Panic("Failed to ListenUDP", err)
	}
	listener1, err := net.ListenUDP("udp", &net.UDPAddr{Port: 8080})
	if err != nil {
		log.Panic("Failed to ListenUDP", err)
	}
	go p.receiveRegistration(listener)
	go p.findTargetAddr(listener1)
	for {
	}
	// n, addr, err = listener.ReadFromUDP(buf)
	// if err != nil {
	// 	log.Panic("2 Fail to readFromUDP")
	// }
	// msg = strings.Split(string(buf[:n]), ":")
	// peers[msg[0]] = addr.String()
	// fmt.Printf("read %d size, from %s, msg: %s",n,addr.String(),buf[:n])
	// // Exchange messages
	// listener.WriteToUDP([]byte(peers[0].String()), peers[1])
	// listener.WriteToUDP([]byte(peers[1].String()), peers[0])
	// Exit()
	// fmt.Println("the server will disconnect after 10 seconds")
	// time.Sleep(time.Second*10)

}

func (p *peers) receiveRegistration(listener *net.UDPConn) {
	for {
		buf := make([]byte, 256)
		n, addr, err := listener.ReadFromUDP(buf)
		if err != nil {
			log.Panic("Fail to readFromUDP")
		}
		msg := strings.Split(string(buf[:n]), ":")
		p.mu.Lock()
		p.addr[msg[1]] = addr
		p.mu.Unlock()
		fmt.Printf("read %d size, from %s, msg: %s \n", n, addr.String(), buf[:n])

		listener.WriteToUDP([]byte("register successful"), addr)
	}

}

func (p *peers) findTargetAddr(listener *net.UDPConn) {
	for {
		buf := make([]byte, 256)
		n, _, err := listener.ReadFromUDP(buf)
		if err != nil {
			log.Panic("Fail to readFromUDP")
		}

		msg := strings.Split(string(buf[:n]), ":")
		for {
			p.mu.Lock()
			t, ok := p.addr[msg[1]]

			if ok {
				listener.WriteToUDP([]byte(t.String()), p.addr[msg[0]])
				p.mu.Unlock()
				return
			}
			p.mu.Unlock()
		}
	}
}

// // func (p *peers) findTargetAddr(listener *net.UDPConn){
// // 	source, dest := p.receiveTargetAddr(listener)

// // }
