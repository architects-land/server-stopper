# Server Stopper - Service

This is the service creating the UNIX socket. 
This service only works with UNIX and UNIX-like OS (like Linux distribution).

## Usage

CLI args:
- `-systemd` if the host uses systemd (default: `false`)
- `-socket [path]` is the path to the UNIX socket (default: `/run/service-stopper.sock`)
- `-minute-before-poweroff [int > 0]` is the time in minutes of inactivity to wait before stopping the server 
(default: `5`)

A service file is provided. You can compile the program with `go build .` (Go 1.24+).