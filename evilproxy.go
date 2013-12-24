package main

import (
	"evilproxy/simulation"
	"flag"
	"fmt"
	"log"
	"net"
	"time"
)

func main() {
	var client = flag.String("connect", ":80", "Client connection address")
	var server = flag.String("server", ":8080", "Server connection address")
	flag.Parse()

	serv, err := net.Listen("tcp", *server)
	if err != nil {
		log.Fatalf("Unable start server on \"%s\". %v\n", *server, err)
	}

	for {
		ssock, err := serv.Accept()
		if err != nil {
			log.Fatalf("Unable accept client connection. %v\n", err)
		}

		// Complete the connection
		go func(client string) {
			csock, err := net.DialTimeout("tcp", client, time.Second*3)
			if err != nil {
				log.Printf("Unable to connect to \"%s\". %v\n", client, err)
			}

			cconn, sconn = simulation.ConstructConnections("")

			// TODO start streaming packets between the connections
		}(*client)
	}
}
