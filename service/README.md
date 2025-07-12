# Server Stopper - Service

This is the service creating the UNIX socket. 
This service only works with UNIX and UNIX-like OS (like Linux distribution).

## Usage

CLI args:
- `-systemd` if the host uses systemd (default: `false`)
- `-socket [path]` is the path to the UNIX socket (default: `/run/service-stopper.sock`)
- `-minute-before-poweroff [int > 0]` is the time in minutes of inactivity to wait before stopping the server 
(default: `5`)
- `-stop-only-minecraft` if the service must only stop Minecraft, only works with systemd (default: `false`)
- `-minecraft-service` is the name of the Minecraft service, only works with systemd and only used if 
`-stop-only-minecraft` is true (default: `minecraft.service`)

A service file is provided. You can compile the program with `go build .` (Go 1.24+).

You can stop the Minecraft server via the RCON with the arg `-stop-minecraft`. 
You must set a password with the arg `-rcon-password` and if you are using a non-standard port, you can set it with
`-rcon-port`.

Use `github.com/willroberts/minecraft-client` to connect to the RCON.