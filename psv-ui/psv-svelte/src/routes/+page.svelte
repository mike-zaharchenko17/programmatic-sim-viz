<script lang="ts">
    import { onDestroy, onMount } from "svelte"
    import { createSocket } from "$lib/socket.svelte";
    import Sankey from "$lib/components/sankey/Sankey.svelte";

    const socket = createSocket("ws://localhost:1323/ws")

    let visibleResults = $derived(socket.auctionResults)

    onMount(() => {
        socket.connect()
    })

    onDestroy(() => {
        socket.disconnect()
    })

</script>

<div class="container">
    <div class="container">
        <div class="toolbar">
            {#if socket.isOpen}
                <button onclick={socket.disconnect}>Disconnect</button>
            {:else}
                <button onclick={socket.connect}>Connect</button>
            {/if}
        </div>
        <div class="sankey-chart">
            <Sankey visibleResults={visibleResults} />
        </div>
    </div>    
</div>

