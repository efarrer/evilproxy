package main

import (
	"bytes"
	"flag"
	"io"
	"log"
	"net"
	"runtime/pprof"
	"sync"
	"time"

	"github.com/efarrer/evilproxy/connection"
	"github.com/efarrer/evilproxy/debug"
	"github.com/efarrer/evilproxy/parser"
)

func main() {
	var client = flag.String("client", ":80", "Client connection address")
	var server = flag.String("server", ":8080", "Server connection address")
	var connections = flag.Int("connections", -1, "Number of connections to allow")
	var debugEnabled = flag.Bool("debug", false, "Enable additional debug functionality")
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

			defer ssock.Close()
			defer csock.Close()
			defer cconn.Close()
			defer sconn.Close()

			// TODO Make sure all socket/connections get closed
		}(*client)

		outstandingConns.Wait()
		if *debugEnabled {
			time.Sleep(1 * time.Second)
			buffer := &bytes.Buffer{}
			pprof.Lookup("goroutine").WriteTo(buffer, 2)
			if cnt, value := debug.OutstandingGoRoutines(buffer.String()); cnt != 1 {
				log.Printf("%v\n\nOutstanding goroutines %v", value, cnt)
			}
		}
	}
}
