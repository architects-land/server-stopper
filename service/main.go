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
	useSystemD           = false
	socketPath           = "/run/server-stopper.sock"
	minuteBeforePowerOff = 5

	numberConnected = -1
	quit            chan interface{}
)

func init() {
	flag.BoolVar(&useSystemD, "systemd", useSystemD, "use systemd")
	flag.StringVar(&socketPath, "socket", socketPath, "path to socket")
	flag.IntVar(&minuteBeforePowerOff, "minute-before-poweroff", minuteBeforePowerOff, "minutes before poweroff")
}

func main() {
	flag.Parse()

	if minuteBeforePowerOff < 1 {
		slog.Error("minutes before poweroff is < 1")
		return
	}

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

	slog.Info("Socket started")

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
	slog.Info("updating number connected", "new", n, "old", numberConnected)
	numberConnected = n
	if n != 0 {
		return
	}

	ticker := time.NewTicker(time.Duration(minuteBeforePowerOff) * time.Minute)
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
