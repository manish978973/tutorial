package main

import (
	"encoding/binary"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/scionproto/scion/go/lib/snet"
	"github.com/scionproto/scion/go/lib/sciond"
)

func check(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func printUsage() {
	fmt.Println("\ntimestamp_server -s ServerSCIONAddress")
	fmt.Println("\tListens for incoming connections and responds back to them right away")
	fmt.Println("\tThe SCION address is specified as ISD-AS,[IP Address]:Port")
	fmt.Println("\tIf server listening port unspecified, a random available one will be used")
	fmt.Println("\tExample SCION address 1-1,[127.0.0.1]:42002\n")
}

func main() {
	var (
		serverAddress string

		err    error
		server *snet.Addr

		udpConnection *snet.Conn
	)

	// Fetch arguments from command line
	flag.StringVar(&serverAddress, "s", "", "Server SCION Address")
	flag.Parse()

	// Create the SCION UDP socket
	if len(serverAddress) > 0 {
		server, err = snet.AddrFromString(serverAddress)
		check(err)
	} else {
		printUsage()
		check(fmt.Errorf("Error, server address needs to be specified with -s"))
	}

	dispatcherAddr := "/run/shm/dispatcher/default.sock"
	snet.Init(server.IA, sciond.GetDefaultSCIONDPath(nil), dispatcherAddr)

	udpConnection, err = snet.ListenSCION("udp4", server)
	check(err)

	receivePacketBuffer := make([]byte, 2500)
	for {
		n, clientAddress, err := udpConnection.ReadFrom(receivePacketBuffer)
		check(err)

		// Packet received, send back response to same client
		m := binary.PutVarint(receivePacketBuffer[n:], time.Now().UnixNano())
		_, err = udpConnection.WriteTo(receivePacketBuffer[: n+m], clientAddress)
		check(err)
		fmt.Println("Received connection from", clientAddress)
	}
}
