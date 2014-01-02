package main

import (
	"evilproxy/connection"
	"evilproxy/parser"
	"flag"
	"io"
	"log"
	"net"
	"sync"
	"time"
)

func main() {
	var client = flag.String("client", ":80", "Client connection address")
	var server = flag.String("server", ":8080", "Server connection address")
	var connections = flag.Int("connections", -1, "Number of connections to allow")
	flag.Parse()

	var outstandingConns sync.WaitGroup

	serv, err := net.Listen("tcp", *server)
	if err != nil {
		log.Fatalf("Unable start server on \"%s\". %v\n", *server, err)
	}

	for i := 0; i != *connections; i++ {

		ssock, err := serv.Accept()
		if err != nil {
			log.Fatalf("Unable accept client connection. %v\n", err)
		}

		outstandingConns.Add(1)
		// Complete the connection
		go func(client string) {
			defer outstandingConns.Done()

			csock, err := net.DialTimeout("tcp", client, time.Second*3)
			if err != nil {
				log.Printf("Unable to connect to \"%s\". %v\n", client, err)
			}

			rule := ""
			cconn, sconn, err := parser.ConstructConnections(rule)
			if err != nil {
				log.Printf("Error (%v) unable to parse rule. \"%v\"\n", rule, err)
			}

			go io.Copy(csock, connection.ConnectionReaderAdaptor(cconn))
			go io.Copy(connection.ConnectionWriterAdaptor(cconn), csock)

			go io.Copy(ssock, connection.ConnectionReaderAdaptor(sconn))
			io.Copy(connection.ConnectionWriterAdaptor(sconn), ssock)

            // TODO Make sure all socket/connections get closed
            // TODO Add a debug option that makes sure all all goroutines are
            // shutdown runtime.NumGoroutine

		}(*client)

		outstandingConns.Wait()
	}
}
