package world.anhgelus.architectsland.serverstopper.config;

import world.anhgelus.architectsland.serverstopper.ServerStopper;

public class ModConfig {
    public String socketPath = "/run/server-stopper.sock";

    public ModConfig() {
        SimpleConfig config = SimpleConfig.of(ServerStopper.MOD_ID).provider(this::provider).request();

        this.socketPath = config.getOrDefault( "socket_path", this.socketPath);
    }

    private String provider( String filename ) {
        return "# Server Stopper configuration file\n\n" +
                "# Path to the UNIX socket.\n" +
                "socket_path = " + this.socketPath;
    }
}

