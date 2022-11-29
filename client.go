package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

const serverIP = "20.208.131.198"

func main() {
	if len(os.Args) < 4 {
		fmt.Println("./client port name target")
	}
	name := os.Args[2]
	target := os.Args[3]
	port, _ := strconv.Atoi(os.Args[1])
	peerChat(name, target, port, serverIP)

}


func peerChat(source, dest string, port int, serverIP string) {
	localAddr := net.UDPAddr{Port: port}
	registerRemoteAddr := net.UDPAddr{
		IP:   net.ParseIP(serverIP),
		Port: 8081,
	}
	targetRemoteAddr := net.UDPAddr{
		IP:   net.ParseIP(serverIP),
		Port: 8080,
	}
	register(localAddr, registerRemoteAddr, source)
	time.Sleep(1 * time.Second)
	var toAddr net.UDPAddr
	for {
		msgReceived := getDestAddr(localAddr, targetRemoteAddr, source, dest)
		if len(msgReceived) != 4 {
			toAddr = parseIP(string(msgReceived))
			fmt.Println("target addr", toAddr)
			break
		} else {
			fmt.Println("target peer is busy")
		}
		time.Sleep(10 * time.Second)
	}
	p2pchat(&localAddr, &toAddr)

}

func parseIP(addr string) net.UDPAddr {
	strs := strings.Split(addr, ":")
	ip := strs[0]
	port, _ := strconv.Atoi(strs[1])
	return net.UDPAddr{
		IP:   net.ParseIP(ip),
		Port: port,
	}
}

func register(localAddr, remoteAddr net.UDPAddr, source string) {
	conn, err := net.DialUDP("udp", &localAddr, &remoteAddr)
	if err != nil {
		log.Panic("failed to DialUDP", err)
	}
	conn.Write([]byte("This is peer:" + source))
	buf := make([]byte, 256)
	n, _, err := conn.ReadFromUDP(buf)
	if err != nil {
		log.Panic("failed to ReadFromUDP", err)
	}
	fmt.Printf("%s", buf[:n])
	conn.Close()
}

func getDestAddr(localAddr, remoteAddr net.UDPAddr, source, dest string) string {
	conn, err := net.DialUDP("udp", &localAddr, &remoteAddr)
	if err != nil {
		log.Panic("failed to DialUDP", err)
	}
	conn.Write([]byte(source + ":" + dest))

	buf := make([]byte, 256)
	n, _, err := conn.ReadFromUDP(buf)
	if err != nil {
		log.Panic("failed to ReadFromUDP", err)
	}
	conn.Close()
	return string(buf[:n])
}

// func stop(localAddr, remoteAddr net.UDPAddr) {

// }

func p2pchat(fromAddr, toAddr *net.UDPAddr) {
	conn, err := net.DialUDP("udp", fromAddr, toAddr)
	if err != nil {
		log.Panic("failed to DialUDP", err)
	}
	n, err := conn.Write([]byte("Hole punching \n"))
	fmt.Println(n, err)

	// goroutine handle
	go func() {
		buf := make([]byte, 256)
		for {
			n, _, err := conn.ReadFromUDP(buf)
			if err != nil {
				fmt.Println("readFromUDP err", err)
				continue
			}
			fmt.Printf("receive: %sp2p> \n", buf[:n])
		}
	}()

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Println("p2p>")
		data, err := reader.ReadString('\n')
		if err != nil {
			log.Panic("failed to read string", err)
		}
		conn.Write([]byte(data))
		if data == "stop" {
			break
		}
	}
}
