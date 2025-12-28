import { writable } from 'svelte/store';

export const botStore = writable({
    fishingEnabled: false,
    name: "Connecting...",
    x: 0,
    y: 0,
    z: 0
});

export function connect() {
    const socket = new WebSocket('ws://localhost:8080/ws');

    socket.onmessage = (event) => {
        const data = JSON.parse(event.data);
        botStore.set(data);
    };

    socket.onclose = () => {
        setTimeout(connect, 1000); // Auto-reconnect if Go restarts
    };
}