<script lang="ts">
    import { onDestroy } from "svelte"
    import { createSocket } from "$lib/socket.svelte";

    const socket = createSocket("ws://localhost:1323/ws")

    onDestroy(() => {
        socket.disconnect()
    })

</script>

<div class="container">

    {#if socket.isOpen}
        <button onclick={socket.disconnect}>Disconnect</button>
    {:else}
        <button onclick={socket.connect}>Connect</button>
    {/if}

    {#each socket.auctionResults as ar}
        <div class="card auction-result">
            <br>
            <p>{JSON.stringify(ar, null, 2)}</p>
        </div>
    {/each}
</div>

