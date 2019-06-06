package main

import (
	"fmt"
	"net"
	"os"

	"github.com/isundaylee/simple_proxy/pkg/server"
)

const listenAddress string = "0.0.0.0:9270"

func main() {
	fmt.Println("Welcome to simple proxy!")

	ln, err := net.Listen("tcp", listenAddress)
	if err != nil {
		fmt.Printf("Failed to listen on %s: %s\n", listenAddress, err.Error())
		os.Exit(1)
	} else {
		fmt.Printf("Started listening on %s.\n", listenAddress)
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Printf("Error while accepting a client: %s\n", err.Error())
		}

		go server.HandleProtocol(conn, conn)
	}
}
