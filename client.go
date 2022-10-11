package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("./client port name")
	}
	name := os.Args[2]
	port, _ := strconv.Atoi(os.Args[1])
	localAddr := net.UDPAddr{Port: port,}
	remoteAddr := net.UDPAddr{
		IP: net.ParseIP("10.20.3.135"),
		Port: 9527,
	}

	conn, err := net.DialUDP("udp",&localAddr, &remoteAddr)
	if err != nil{
		log.Panic("failed to DialUDP",err)
	}
	conn.Write([]byte("I am a peer: " + name))

	buf := make([]byte,256)
	n, _ , err := conn.ReadFromUDP(buf)
	if err != nil{
		log.Panic("failed to ReadFromUDP",err)
	}
	toAddr:= parseIP(string(buf[:n]))

	fmt.Println("target addr", toAddr)
	conn.Close()

	p2pchat(&localAddr, &toAddr)
}

func parseIP(addr string) (net.UDPAddr){
	strs := strings.Split(addr, ":")
	// if len(strs) != 2 {
	// 	fmt.Println("ip addr is not valid")
	// 	return nil , errors.New("ipaddr err")
	// }

	ip := strs[0]
	port, _ := strconv.Atoi(strs[1])
	return net.UDPAddr{
		IP: net.ParseIP(ip),
		Port: port,
	}

}

func p2pchat(fromAddr, toAddr *net.UDPAddr) {
	conn, err := net.DialUDP("udp", fromAddr, toAddr)
	if err != nil{
		log.Panic("failed to DialUDP",err)
	}
	n, err := conn.Write([]byte("Hole punching \n"))
	fmt.Println(n,err)

	// goroutine handle 
	go func()  {
		buf := make([]byte, 256)
		for{
			n , _ , err := conn.ReadFromUDP(buf)
			if err != nil{
				fmt.Println("readFromUDP err",err)
				continue
			}
			fmt.Printf("receive: %sp2p> \n", buf[:n])
		}
		
	}()

	reader := bufio.NewReader(os.Stdin)
	for{
		fmt.Println("p2p>")
		data, err := reader.ReadString('\n')
		if err != nil{
			log.Panic("failed to read string",err)
		}
		conn.Write([]byte(data))
	}
}