import { writable } from 'svelte/store';

export const botStore = writable({
    fishingEnabled: false,
    name: "Connecting...",
    x: 0,
    y: 0,
    z: 0
});

let socket;

export function connect() {
    socket = new WebSocket('ws://localhost:8080/ws');

    socket.onmessage = (event) => {
        const data = JSON.parse(event.data);
        botStore.set(data);
    };

    socket.onclose = () => {
        setTimeout(connect, 1000); // Auto-reconnect if Go restarts
    };
}

export function sendToggleFishing() {
    console.log("Toggling fishing state. Current Socket State:", socket?.readyState);
    if (socket && socket.readyState === WebSocket.OPEN) {
        const command = { type: "TOGGLE_FISHING" };
        socket.send(JSON.stringify(command));
    }
}