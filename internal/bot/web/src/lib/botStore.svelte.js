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

        // This is needed to prevent breaking the drag-and-drop UI
        if (!this.isDraggingWaypoint) {
            this.waypoints = data.waypoints;
        }
    }

    toggleFishing() {
        socket.send(JSON.stringify({ type: "TOGGLE_FISHING" }));
    }

    reorderWaypoints(newList) {
        this.waypoints = newList;
        socket.send(JSON.stringify({
            type: "REORDER_WAYPOINTS",
            data: newList
        }));
    }
}

export const bot = new BotStore();