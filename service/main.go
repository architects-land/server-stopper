package main

import (
	"flag"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

var (
	useSystemD bool
	socketPath string
)

func init() {
	flag.BoolVar(&useSystemD, "systemd", false, "use systemd")
	flag.StringVar(&socketPath, "socket", "/run/service-stopper.sock", "path to socket")
}

func main() {
	flag.Parse()

	socket, err := net.Listen("unix", socketPath)
	if err != nil {
		panic(err)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		err = os.Remove(socketPath)
		if err != nil {
			panic(err)
		}
		os.Exit(0)
	}()

	for {
		// Accept an incoming connection.
		conn, err := socket.Accept()
		if err != nil {
			log.Fatal(err)
		}

		// Handle the connection in a separate goroutine.
		go func(conn net.Conn) {
			defer conn.Close()
			// Create a buffer for incoming data.
			buf := make([]byte, 4096)

			// Read data from the connection.
			n, err := conn.Read(buf)
			if err != nil {
				log.Fatal(err)
			}

			// Echo the data back to the connection.
			_, err = conn.Write(buf[:n])
			if err != nil {
				log.Fatal(err)
			}
		}(conn)
	}
}
