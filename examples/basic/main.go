package main

import (
	"fmt"
	"github.com/miladrahimi/gorelay"
	"os"
	"strconv"
)

func main() {
	if len(os.Args) != 5 {
		fmt.Println("Usage: ./tcp_udp_relay <sourcePort> <destinationIP> <destinationTCPPort> <destinationUDPPort>")
		os.Exit(1)
	}

	sourcePort, _ := strconv.Atoi(os.Args[1])
	destinationIP := os.Args[2]
	destinationTCPPort, _ := strconv.Atoi(os.Args[3])
	destinationUDPPort, _ := strconv.Atoi(os.Args[4])

	go func() {
		tcp := gorelay.NewTcpRelay()
		err := tcp.Relay(sourcePort, destinationTCPPort, destinationIP)
		if err != nil {
			fmt.Println("TCP Relay Error:", err)
		}
	}()

	udp := gorelay.NewUdpRelay()
	err := udp.Relay(sourcePort, destinationUDPPort, destinationIP)
	if err != nil {
		fmt.Println("UDP Relay Error:", err)
	}
}
