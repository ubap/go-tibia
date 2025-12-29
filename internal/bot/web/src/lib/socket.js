// src/lib/socket.js
import { bot } from './botStore.svelte.js';

export let socket;

export function connect() {
    socket = new WebSocket('ws://127.0.0.1:8080/ws');

    socket.onmessage = (event) => {
        const data = JSON.parse(event.data);
        bot.updateFromSnapshot(data);
    };

    socket.onclose = () => setTimeout(connect, 1000);
}