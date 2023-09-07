package gorelay

import (
	"fmt"
	"net"
	"strconv"
	"sync/atomic"
)

type TcpRelay struct {
	logger           Logger
	bytesTransferred int64
	bufferSize       int
}

func (r *TcpRelay) SetBufferSize(bufferSize int) {
	r.bufferSize = bufferSize
}

func (r *TcpRelay) SetLogger(logger Logger) {
	r.logger = logger
}

func (r *TcpRelay) BytesTransferred() int64 {
	return r.bytesTransferred
}

func (r *TcpRelay) Relay(sourcePortNumber, destinationPortNumber int, destinationHost string) error {
	sourcePort := strconv.Itoa(sourcePortNumber)
	destinationPort := strconv.Itoa(destinationPortNumber)

	listener, err := net.Listen("tcp", ":"+sourcePort)
	if err != nil {
		return err
	}
	defer func(listener net.Listener) {
		if err := listener.Close(); err != nil {
			r.logger.Info("cannot close tcp listener: " + err.Error())
		}
	}(listener)

	r.logger.Info(fmt.Sprintf("relaying tcp :%s to %s:%s", sourcePort, destinationHost, destinationPort))

	for {
		sourceConnection, err := listener.Accept()
		if err != nil {
			r.logger.Error("cannot accept tcp source: " + err.Error())
			continue
		}

		destinationConnection, err := net.Dial("tcp", destinationHost+":"+destinationPort)
		if err != nil {
			r.logger.Error("cannot connect to tcp destination: " + err.Error())
			if err = sourceConnection.Close(); err != nil {
				r.logger.Info("cannot close source tcp connection: " + err.Error())
			}
			continue
		}

		go r.copy(sourceConnection, destinationConnection)
		go r.copy(destinationConnection, sourceConnection)
	}
}

func (r *TcpRelay) copy(source, destination net.Conn) {
	defer func(source, destination net.Conn) {
		if err := source.Close(); err != nil {
			r.logger.Info("cannot close source tcp connection: " + err.Error())
		}
		if err := destination.Close(); err != nil {
			r.logger.Info("cannot close destination tcp connection: " + err.Error())
		}
	}(source, destination)

	buffer := make([]byte, r.bufferSize)
	for {
		n, err := source.Read(buffer)
		if err != nil {
			return
		}

		if _, err = destination.Write(buffer[:n]); err != nil {
			return
		}

		go atomic.AddInt64(&r.bytesTransferred, int64(n))
	}
}

func NewTcpRelay() *TcpRelay {
	return &TcpRelay{bytesTransferred: 0, bufferSize: 1024, logger: NewBasicLogger()}
}
