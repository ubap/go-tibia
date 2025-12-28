<script>
    import { onMount } from 'svelte';
    import { connect } from '$lib/botStore';

    onMount(() => {
        connect();
    });
</script>

<nav style="padding: 1rem; background: #333; color: white;">
    <a href="/stats" style="color: white; margin-right: 10px;">Character Stats</a>
    <a href="/fishing" style="color: white;">Fishing</a>
</nav>

<slot />

