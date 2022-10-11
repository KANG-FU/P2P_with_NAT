package main

import (
	"fmt"
	"log"
	"net"
	"time"
)
// establish UDP server
// receive 2 requests
// exchange messages
// quit not affect A and B commuication
func main() {
	fmt.Println("begin server")
	//服务器启动侦听
	//listener, err := net.ListenUDP("udp", &net.UDPAddr{Port: 9527})
	listener, err := net.ListenUDP("udp", &net.UDPAddr{
		IP: net.ParseIP("185.25.192.150"),
		Port: 43346,
	})
	if err != nil{
		log.Panic("Failed to ListenUDP", err)
	}
	defer listener.Close()
	peers := make([]*net.UDPAddr,2)
	buf := make([]byte ,256)
	n, addr, err := listener.ReadFromUDP(buf)
	if err != nil {
		log.Panic("1 Fail to readFromUDP")
	}
	peers[0] = addr
	fmt.Printf("read %d size, from %s, msg: %s \n",n,addr.String(),buf[:n])

	n, addr, err = listener.ReadFromUDP(buf)
	if err != nil {
		log.Panic("2 Fail to readFromUDP")
	}
	peers[1] = addr
	fmt.Printf("read %d size, from %s, msg: %s",n,addr.String(),buf[:n])
	// Exchange messages
	listener.WriteToUDP([]byte(peers[0].String()), peers[1])
	listener.WriteToUDP([]byte(peers[1].String()), peers[0])
	// Exit()
	fmt.Println("the server will disconnect after 10 seconds")
	time.Sleep(time.Second*10)

}