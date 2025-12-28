<script>
    import { botStore } from '$lib/botStore';
</script>

<h1>Character Stats</h1>

<div class="card">
    <!-- $botStore.hp will update automatically as Go sends data -->
    <p>Name: <span style="color: red;">{$botStore.name}</span></p>

    <p>Bot Status:
        {#if $botStore.fishingEnabled}
            <b style="color: green;">FISHING</b>
        {:else}
            <b style="color: gray;">IDLE</b>
        {/if}
    </p>
</div>