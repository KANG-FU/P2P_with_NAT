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

func main() {
	if len(os.Args) < 4 {
		fmt.Println("./client port name target")
	}
	name := os.Args[2]
	target := os.Args[3]
	port, _ := strconv.Atoi(os.Args[1])
	localAddr := net.UDPAddr{Port: port}
	registerRemoteAddr := net.UDPAddr{
		IP:   net.ParseIP("10.20.3.135"),
		Port: 8081,
	}
	targetRemoteAddr := net.UDPAddr{
		IP:   net.ParseIP("10.20.3.135"),
		Port: 8080,
	}
	register(localAddr, registerRemoteAddr, name)
	time.Sleep(10 * time.Second)
	toAddr := getDestAddr(localAddr, targetRemoteAddr, name, target)
	fmt.Println("target addr", toAddr)
	p2pchat(&localAddr, &toAddr)
}

func parseIP(addr string) net.UDPAddr {
	strs := strings.Split(addr, ":")
	// if len(strs) != 2 {
	// 	fmt.Println("ip addr is not valid")
	// 	return nil , errors.New("ipaddr err")
	// }
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

func getDestAddr(localAddr, remoteAddr net.UDPAddr, source, dest string) net.UDPAddr {
	conn, err := net.DialUDP("udp", &localAddr, &remoteAddr)
	if err != nil {
		log.Panic("failed to DialUDP", err)
	}
	conn.Write([]byte(source + ":" + dest))
	fmt.Println("sending target")

	buf := make([]byte, 256)
	n, _, err := conn.ReadFromUDP(buf)
	if err != nil {
		log.Panic("failed to ReadFromUDP", err)
	}
	toAddr := parseIP(string(buf[:n]))
	conn.Close()
	return toAddr
}

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
			fmt.Println("received")
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
	}
}
