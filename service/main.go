package main

import (
	"flag"
	"log"
	"log/slog"
	"net"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

var (
	useSystemD bool
	socketPath string

	numberConnected int
	quit            chan interface{}
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
		conn, err := socket.Accept()
		if err != nil {
			log.Fatal(err)
		}

		// Sync is used here because the socket must be sync
		handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	// Create a buffer for incoming data.
	buf := make([]byte, 4096)

	n, err := conn.Read(buf)
	if err != nil {
		slog.Error(err.Error(), "position", "cannot read")
		_, err = conn.Write([]byte("1"))
		if err != nil {
			slog.Error(err.Error(), "position", "cannot write error")
		}
		return
	}

	s := string(buf[:n])
	nbr, err := strconv.Atoi(s)
	if err != nil {
		slog.Warn(err.Error(), "position", "converting to int")
		_, err = conn.Write([]byte("2"))
		if err != nil {
			slog.Error(err.Error(), "position", "cannot write error")
		}
		return
	}

	if nbr < 0 {
		slog.Warn("negative number")
		_, err = conn.Write([]byte("3"))
		if err != nil {
			slog.Error(err.Error(), "position", "cannot write error")
		}
		return
	}

	updateConnected(nbr)

	_, err = conn.Write([]byte("0"))
	if err != nil {
		slog.Error(err.Error(), "position", "cannot write response")
	}
}

func updateConnected(n int) {
	if numberConnected == 0 {
		if n == 0 {
			slog.Warn("numberConnected already set to 0")
			return
		}
		quit <- true
	}
	numberConnected = n
	if n != 0 {
		return
	}

	ticker := time.NewTicker(5 * time.Minute)
	if quit != nil {
		quit <- true
	}
	quit = make(chan interface{})
	go func() {
		for {
			select {
			case <-ticker.C:
				stop()
			case <-quit:
				ticker.Stop()
				close(quit)
				return
			}
		}
	}()
}

func stop() {
	quit <- true
	var cli string
	if useSystemD {
		cli = "systemctl poweroff"
	} else {
		cli = "poweroff"
	}
	cmd := exec.Command(cli)
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		panic(err)
	}
}
