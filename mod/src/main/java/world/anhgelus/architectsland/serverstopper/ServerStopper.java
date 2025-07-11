package world.anhgelus.architectsland.serverstopper;

import net.fabricmc.api.ModInitializer;
import net.fabricmc.fabric.api.entity.event.v1.ServerPlayerEvents;
import net.fabricmc.fabric.api.event.lifecycle.v1.ServerLifecycleEvents;
import net.minecraft.server.MinecraftServer;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import world.anhgelus.architectsland.serverstopper.config.ModConfig;

import java.io.IOException;
import java.net.StandardProtocolFamily;
import java.net.UnixDomainSocketAddress;
import java.nio.ByteBuffer;
import java.nio.channels.SocketChannel;
import java.nio.file.Path;


public class ServerStopper implements ModInitializer {
    public static final String MOD_ID = "server-stopper";
    public static final Logger LOGGER = LoggerFactory.getLogger(MOD_ID);
    public static ModConfig CONFIG = new ModConfig();

    @Override
    public void onInitialize() {
        LOGGER.info("Initializing Server Stopper");

        ServerLifecycleEvents.SERVER_STARTED.register(server -> sendPlayersConnected(server.getPlayerManager().getPlayerList().size()));

        ServerPlayerEvents.JOIN.register(player -> {
            final var server = player.getServer();
            if (server == null) {
                LOGGER.warn("Server is null during join");
                return;
            }
            sendPlayersConnected(server.getPlayerManager().getPlayerList().size());
        });

        ServerPlayerEvents.LEAVE.register(player -> {
            final var server = player.getServer();
            if (server == null) {
                LOGGER.warn("Server is null during leave");
                return;
            }
            sendPlayersConnected(server.getPlayerManager().getPlayerList().size() - 1);
        });
    }

    private void sendPlayersConnected(int n) {
        if (n < 0) {
            LOGGER.warn("Number of player connected is below 0");
            n = 0;
        }
        try (final var channel = SocketChannel.open(StandardProtocolFamily.UNIX)) {
            final var path = Path.of(CONFIG.socketPath);
            channel.connect(UnixDomainSocketAddress.of(path));

            final var buffer = ByteBuffer.allocate(1024);
            buffer.clear().put(String.format("%d", n).getBytes()).flip();

            while (buffer.hasRemaining()) channel.write(buffer);
        } catch (IOException e) {
            throw new RuntimeException(e);
        }
    }
}
