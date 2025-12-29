<script>
    import { bot } from '$lib/botStore.svelte';
    import { dndzone } from 'svelte-dnd-action';
    import { flip } from 'svelte/animate';

    const flipDurationMs = 200;

    // Available waypoint types for Tibia
    const types = ["Walk", "Node", "Rope", "Ladder", "Shovel", "Machete"];

    function handleDndConsider(e) {
        console.log("Handle DnD:", e.detail);
        bot.isDraggingWaypoint = true;
        bot.waypoints = e.detail.items;
    }

    function handleDndFinalize(e) {
        console.log("Finalized DnD:", e.detail);
        bot.waypoints = e.detail.items;
        bot.isDraggingWaypoint = false;

        // Send the new order back to Go proxy
        // socket.send(JSON.stringify({
        //     type: "REORDER_WAYPOINTS",
        //     data: bot.waypoints
        // }));
    }

    function removeWaypoint(id) {
        bot.waypoints = bot.waypoints.filter(w => w.id !== id);
        // Notify Go backend of deletion...
    }
</script>

<div class="max-w-2xl mx-auto space-y-4">
    <div class="flex justify-between items-center">
        <h2 class="text-2xl font-bold text-white">Cavebot Waypoints</h2>
        <button class="bg-orange-600 hover:bg-orange-700 text-white px-4 py-2 rounded-lg text-sm font-bold">
            + ADD CURRENT POS
        </button>
    </div>

    <!-- Reorderable List Container -->
    <section
            use:dndzone={{items: bot.waypoints, flipDurationMs}}
            onconsider={handleDndConsider}
            onfinalize={handleDndFinalize}
            class="space-y-2 min-h-[100px]"
    >
        {#each bot.waypoints as wp (wp.id)}
            <div
                    animate:flip={{duration: flipDurationMs}}
                    class="bg-slate-900 border border-slate-800 p-3 rounded-xl flex items-center gap-4 group hover:border-orange-500/50 transition-colors cursor-grab active:cursor-grabbing"
            >
                <!-- Drag Handle Icon -->
                <div class="text-slate-600 group-hover:text-slate-400">
                    <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="9" cy="5" r="1"/><circle cx="9" cy="12" r="1"/><circle cx="9" cy="19" r="1"/><circle cx="15" cy="5" r="1"/><circle cx="15" cy="12" r="1"/><circle cx="15" cy="19" r="1"/></svg>
                </div>

                <!-- Waypoint Type Selector -->
                <select class="bg-slate-800 border-none text-xs rounded-md text-orange-400 font-bold px-2 py-1 focus:ring-1 focus:ring-orange-500">
                    {#each types as t}
                        <option selected={wp.type === t}>{t}</option>
                    {/each}
                </select>

                <!-- Coordinates -->
                <div class="flex gap-2 font-mono text-sm text-slate-300 flex-1">
                    <span class="bg-slate-950 px-2 py-0.5 rounded border border-slate-800"><span class="text-slate-500 mr-1">X:</span>{wp.x}</span>
                    <span class="bg-slate-950 px-2 py-0.5 rounded border border-slate-800"><span class="text-slate-500 mr-1">Y:</span>{wp.y}</span>
                    <span class="bg-slate-950 px-2 py-0.5 rounded border border-slate-800"><span class="text-slate-500 mr-1">Z:</span>{wp.z}</span>
                </div>

                <!-- Delete Button -->
                <button
                        onclick={() => removeWaypoint(wp.id)}
                        class="text-slate-600 hover:text-red-500 p-1 opacity-0 group-hover:opacity-100 transition-opacity"
                >
                    <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M3 6h18"/><path d="M19 6v14c0 1-1 2-2 2H7c-1 0-2-1-2-2V6"/><path d="M8 6V4c0-1 1-2 2-2h4c1 0 2 1 2 2v2"/></svg>
                </button>
            </div>
        {/each}
    </section>

    {#if bot.waypoints.length === 0}
        <div class="text-center py-10 border-2 border-dashed border-slate-800 rounded-2xl text-slate-500">
            No waypoints yet. Walk around and click "Add Current Pos".
        </div>
    {/if}
</div>