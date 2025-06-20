package main

import "flag"

var (
	useSystemD bool
	socketPath string
)

func init() {
	flag.BoolVar(&useSystemD, "systemd", false, "use systemd")
	flag.StringVar(&socketPath, "socket", "/run/service-stopped.sock", "path to socket")
}

func main() {
	flag.Parse()
}
