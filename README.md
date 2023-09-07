# GoRelay

GoRelay is a Go (Golang) package designed for relaying TCP and UDP traffic. It provides the capability to listen on TCP
and UDP ports, then efficiently forward incoming network traffic to a designated destination (host:port).

## Documentation

### Requirements

It requires Go v1.21 or newer versions.

### Installation

To install this package, run the following command in your project directory.

```shell
go get github.com/miladrahimi/gorelay
```

Next, include it in your application:

```go
import "github.com/miladrahimi/gorelay"
```

## Quick Start

The following example shows how to use the package:

```go
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
```

## License

GoRelay is initially created by [Milad Rahimi](https://miladrahimi.com) and released under
the [MIT License](http://opensource.org/licenses/mit-license.php).

