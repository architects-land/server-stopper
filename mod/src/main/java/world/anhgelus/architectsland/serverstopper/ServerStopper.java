package world.anhgelus.architectsland.serverstopper;

import net.fabricmc.api.ModInitializer;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import world.anhgelus.architectsland.serverstopper.config.ModConfig;


public class ServerStopper implements ModInitializer {
    public static final String MOD_ID = "server-stopper";
    public static final Logger LOGGER = LoggerFactory.getLogger(MOD_ID);
    public static ModConfig CONFIG = new ModConfig();

    @Override
    public void onInitialize() {
        LOGGER.info("Initializing Server Stopper");
    }
}
