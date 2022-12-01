package main

import (
	"fmt"
	"log"
	"net"
	"time"

	"strings"
	"sync"
)

// establish UDP server
// receive 2 requests
// exchange messages
// quit not affect A and B commuication
type peers struct {
	mu    sync.RWMutex
	addr  map[string]*net.UDPAddr
	state map[string]bool
	// socketRegistration  *net.UDPConn
	// socketReceive *net.UDPConn
}

func main() {
	fmt.Println("begin server")

	//peers := make([]*net.UDPAddr,2)
	p := peers{
		addr:  make(map[string]*net.UDPAddr),
		state: make(map[string]bool),
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
		if _, ok := p.addr[msg[1]]; ok {
			fmt.Printf("%s has already registered \n", addr.String())
			listener.WriteToUDP([]byte("already registered"), addr)
		} else {
			p.addr[msg[1]] = addr
			p.state[msg[1]] = false
			fmt.Printf("read %d size, from %s, msg: %s \n", n, addr.String(), buf[:n])
			listener.WriteToUDP([]byte("register successful"), addr)
		}
		p.mu.Unlock()
	}

}

// func (p *peers) receiveStop(listener *net.UDPConn) {
// 	for {
// 		buf := make([]byte, 256)
// 		n, addr, err := listener.ReadFromUDP(buf)
// 		if err != nil {
// 			log.Panic("Fail to readFromUDP")
// 		}
// 		msg := strings.Split(string(buf[:n]), ":")
// 		p.mu.Lock()
// 		p.state[msg[1]] = false
// 		p.mu.Unlock()
// 		fmt.Printf("from %s, msg: %s  \n", n, addr.String(), buf[:n])
// 		listener.WriteToUDP([]byte("ACK"), addr)
// 	}

// }

func (p *peers) findTargetAddr(listener *net.UDPConn) {
	for {
		buf := make([]byte, 256)
		n, _, err := listener.ReadFromUDP(buf)
		if err != nil {
			log.Panic("Fail to readFromUDP")
		}

		msg := strings.Split(string(buf[:n]), ":")
		done := make(chan string, 1)
		go p.checkPeerMap(msg, done)
		select {
		case <-time.After(time.Second * 10):
			fmt.Println("Timeout")
			listener.WriteToUDP([]byte("currently not online"), p.addr[msg[0]])
		case m := <-done:
			listener.WriteToUDP([]byte(m), p.addr[msg[0]])
		}
	}
}

func (p *peers) checkPeerMap(msg []string, done chan string) {
	p.mu.Lock()
	t, ok := p.addr[msg[1]]

	if ok {
		if p.state[msg[1]] {
			done <- "busy"
		} else {
			p.state[msg[1]] = true
			done <- t.String()
		}
	}
	p.mu.Unlock()
}
