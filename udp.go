package gorelay

import (
	"fmt"
	"net"
	"strconv"
	"sync/atomic"
)

type UdpRelay struct {
	logger           Logger
	bytesTransferred int64
	bufferSize       int
}

func (r *UdpRelay) SetBufferSize(bufferSize int) {
	r.bufferSize = bufferSize
}

func (r *UdpRelay) SetLogger(logger Logger) {
	r.logger = logger
}

func (r *UdpRelay) Relay(sourcePortNumber, destinationPortNumber int, destinationHost string) error {
	sourcePort := strconv.Itoa(sourcePortNumber)
	destinationPort := strconv.Itoa(destinationPortNumber)

	listenAddr, err := net.ResolveUDPAddr("udp", ":"+sourcePort)
	if err != nil {
		return err
	}

	destinationAddr, err := net.ResolveUDPAddr("udp", destinationHost+":"+destinationPort)
	if err != nil {
		return err
	}

	connection, err := net.ListenUDP("udp", listenAddr)
	if err != nil {
		return err
	}
	defer func(connection *net.UDPConn) {
		if err := connection.Close(); err != nil {
			r.logger.Info("cannot close udp connection: " + err.Error())
		}
	}(connection)

	r.logger.Info(fmt.Sprintf("relaying udp :%s to %s:%s", sourcePort, destinationHost, destinationPort))

	buffer := make([]byte, r.bufferSize)
	for {
		n, _, err := connection.ReadFromUDP(buffer)
		if err != nil {
			continue
		}

		_, err = connection.WriteToUDP(buffer[:n], destinationAddr)
		if err != nil {
			continue
		}

		go atomic.AddInt64(&r.bytesTransferred, int64(n))
	}
}

func NewUdpRelay() *UdpRelay {
	return &UdpRelay{bytesTransferred: 0, bufferSize: 1024, logger: NewBasicLogger()}
}
