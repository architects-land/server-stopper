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

	stopOnlyMinecraft = false
	minecraftService  = "minecraft.service"

	numberConnected = -1
	quit            chan interface{}
)

func init() {
	flag.BoolVar(&useSystemD, "systemd", useSystemD, "use systemd")
	flag.BoolVar(&stopOnlyMinecraft, "stop-only-minecraft", stopOnlyMinecraft, "shutdown only minecraft")
	flag.StringVar(&socketPath, "socket", socketPath, "path to socket")
	flag.StringVar(&minecraftService, "minecraft-service", minecraftService, "name of the minecraft service")
	flag.IntVar(&minuteBeforePowerOff, "minute-before-poweroff", minuteBeforePowerOff, "minutes before poweroff")
}

func main() {
	flag.Parse()

	if minuteBeforePowerOff < 1 {
		slog.Error("Minutes before poweroff is < 1")
		return
	}

	socket, err := net.Listen("unix", socketPath)
	if err != nil {
		panic(err)
	}
	if err = os.Chmod(socketPath, 0777); err != nil {
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
		return
	}

	s := string(buf[:n])
	nbr, err := strconv.Atoi(s)
	if err != nil {
		slog.Error(err.Error(), "position", "converting to int")
		return
	}

	if nbr < 0 {
		slog.Warn("negative number")
		return
	}

	updateConnected(nbr)
}

func updateConnected(n int) {
	if numberConnected == 0 {
		if n == 0 {
			slog.Warn("numberConnected already set to 0")
			return
		}
		quit <- true
	}
	slog.Info("Updating number connected", "new", n, "old", numberConnected)
	numberConnected = n
	if n != 0 {
		return
	}
	slog.Info("Starting timer to shutdown the server")

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
				ticker.Stop()
				close(quit)
				return
			case <-quit:
				slog.Info("Stopping timer to shutdown the server")
				ticker.Stop()
				close(quit)
				return
			}
		}
	}()
}

func stop() {
	slog.Info("Stopping the server...")
	var cmd *exec.Cmd
	if stopOnlyMinecraft {
		if !useSystemD {
			slog.Error("stopping only minecraft is not supported without systemd")
			return
		}
		cmd = exec.Command("systemctl", "stop", minecraftService)
		if err := cmd.Run(); err != nil {
			slog.Error(err.Error())
		}
		return
	}
	if useSystemD {
		cmd = exec.Command("systemctl", "poweroff")
	} else {
		cmd = exec.Command("poweroff")
	}
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		err2 := os.Remove(socketPath)
		if err2 != nil {
			slog.Error(err2.Error(), "position", "cannot remove socket")
		}
		panic(err)
	}
}
