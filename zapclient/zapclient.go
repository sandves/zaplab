package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", ":12110")
	checkError(err)
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	checkError(err)

	SubscribeAndPrintStatistics(conn)
}

func SubscribeAndPrintStatistics(conn net.Conn) {
	for {
		var stats [512]byte
		_, err := conn.Read(stats[0:])
		checkError(err)
		fmt.Println(string(stats[0:]))
	}
}

func checkError(err error) {
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
		os.Exit(1)
	}
}
