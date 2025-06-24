# Server Stopper

Minecraft mod combined with a service aiming to stop the host server if the Minecraft server is inactive.

This mono-repo contains a Minecraft mod and the service managing the server.
The Minecraft mod is in the `mod` folder and the service is in the `service` folder.

It only works with UNIX and UNIX-like OS (like Linux distributions).

The full documentations is available [here](https://architects-land.github.io/minecraft-scaleway-frontend/server-stopper.html).

## Technologies

Minecraft mod:
- Fabric + YARN mappings
- Fabric API
- Java 21
- Minecraft 1.21.6

Service:
- Go 1.24
