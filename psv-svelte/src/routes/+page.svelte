<script lang="ts">
    import { onDestroy, onMount } from "svelte"
    import { createSocket } from "$lib/socket.svelte";
    import Sankey from "$lib/components/sankey/Sankey.svelte";
    import KPIController from "$lib/components/exploratory/KPIController.svelte";

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
        <div class="flex flex-col lg:flex-row justify-center">
            <div class="card rounded-none flex-1 min-w-0 preset-filled-surface-100-900 border border-surface-200-800 divide-surface-200-800 py-6 px-6">
                <h2 class="h5 mb-2">Auction Flow</h2>
                <Sankey visibleResults={visibleResults} />
            </div>
            <div class="card rounded-none w-full lg:w-72 shrink-0 preset-filled-surface-100-900 py-6 px-6 border-t border-r border-b border-surface-200-800 divide-surface-200-800">
                <h2 class="h5 mb-2">Performance</h2>
                <KPIController visibleResults={visibleResults} />
            </div>
        </div>
    </div>
</div>    

