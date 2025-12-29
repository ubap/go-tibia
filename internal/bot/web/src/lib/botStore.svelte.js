// src/lib/botStore.svelte.js
import { socket } from './socket.js'; // We'll move socket logic here

class BotStore {
    // These are reactive properties ($state)
    name = $state("Connecting...");
    hp = $state(0);
    mana = $state(0);
    x = $state(0);
    y = $state(0);
    z = $state(0);
    fishingEnabled = $state(false);
    lighthackEnabled = $state(false);
    lighthackLevel = $state(15);
    lighthackColor = $state(0xD7);

    // Waypoint list
    waypoints = $state([]);

    isDraggingWaypoint = false;

    // Methods to update state
    updateFromSnapshot(data) {
        this.name = data.name;
        this.hp = data.hp;
        this.mana = data.mana;
        this.x = data.x;
        this.y = data.y;
        this.z = data.z;

        this.fishingEnabled = data.fishingEnabled;

        this.lighthackEnabled = data.lighthackEnabled;
        this.lighthackLevel = data.lighthackLevel;
        this.lighthackColor = data.lighthackColor;

        // This is needed to prevent breaking the drag-and-drop UI
        if (!this.isDraggingWaypoint) {
            this.waypoints = data.waypoints;
        }
    }

    toggleFishing() {
        socket.send(JSON.stringify({ type: "TOGGLE_FISHING" }));
    }

    toggleLighthack = () => {
        this.lighthackEnabled = !this.lighthackEnabled;
        this.sendLighthackUpdate();
    };

    setLighthackLevel = (val) => {
        let level = parseInt(val);

        // Handle empty input or NaN
        if (isNaN(level)) level = 0;

        // Clamp the value between 0 and 16 so the slider/logic doesn't break
        if (level < 0) level = 0;
        if (level > 16) level = 16;

        this.lighthackLevel = level;
        this.sendLighthackUpdate();
    };

    setLighthackColor = (val) => {
        this.lighthackColor = parseInt(val);
        this.sendLighthackUpdate();
    };

    sendLighthackUpdate = () => {
        // Now 'this' will correctly point to the BotStore instance
        socket.send(JSON.stringify({
            type: "SET_LIGHTHACK",
            data: {
                enabled: this.lighthackEnabled,
                level: this.lighthackLevel,
                color: this.lighthackColor
            }
        }));
    };

    reorderWaypoints(newList) {
        this.waypoints = newList;
        socket.send(JSON.stringify({
            type: "REORDER_WAYPOINTS",
            data: newList
        }));
    }
}

export const bot = new BotStore();