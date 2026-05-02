<script lang="ts">
    import { onDestroy, onMount } from "svelte"
    import { createSocket } from "$lib/socket.svelte";
    import Sankey from "$lib/components/sankey/Sankey.svelte";

    const socket = createSocket("ws://localhost:1323/ws")

    let visibleResults = $derived(socket.auctionResults)

    onDestroy(() => {
        socket.disconnect()
    })

</script>

<div class="container mx-auto px-2">
    <div class="flex flex-col gap-4 my-4">
        <div class="flex justify-center">
            {#if socket.isOpen}
                <button type="button" class="btn preset-filled" onclick={socket.disconnect}>Disconnect</button>
            {:else}
                <button type="button" class="btn preset-filled" onclick={socket.connect}>Connect</button>
            {/if}
        </div>
        <div class="card w-fit mx-auto bg-zinc-900 py-6 px-6">
            <Sankey visibleResults={visibleResults} />
        </div>
    </div>
</div>    

