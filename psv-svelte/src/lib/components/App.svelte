<script lang="ts">
    import { onDestroy } from "svelte"
    import { createSocket } from "$lib/socket.svelte";
    import Sankey from "$lib/components/sankey/Sankey.svelte";
    import ExploreSkeleton from "$lib/components/exploratory/ExploreSkeleton.svelte";
    import type { Scope } from "$lib/types/types";

    const socket = createSocket("ws://localhost:1323/ws")

    let currentScope = $state<Scope>({kind: "global"})

    const setScope = ((newScope: Scope) => {
        currentScope = newScope
    })

    const clearScope = () => {
        currentScope = {kind: "global"}
    }

    let visibleResults = $derived(socket.auctionResults)

    onDestroy(() => {
        socket.disconnect()
    })

</script>

<div class="flex flex-col gap-4 my-4">
    <div class="flex justify-center">
        {#if socket.isOpen}
            <button type="button" class="btn preset-filled" onclick={socket.disconnect}>Disconnect</button>
        {:else}
            <button type="button" class="btn preset-filled" onclick={socket.connect}>Connect</button>
        {/if}
    </div>
    <div class="flex flex-col lg:flex-row justify-center">
        <div class="card rounded-none flex-1 min-w-0 preset-filled-surface-100-900 border border-surface-200-800 divide-surface-200-800 py-6 px-6 flex flex-col">
            <h2 class="h5 mb-2">Auction Flow</h2>
            <Sankey 
                visibleResults={visibleResults} 
                setScope={setScope} 
            />
          </div>
        <div class="card rounded-none border border-t-0 lg:border-t lg:border-l-0 w-full lg:w-96 shrink-0 preset-filled-surface-100-900 py-6 px-6 border-surface-200-800 divide-surface-200-800">
            <h2 class="h5 mb-2">Performance</h2>
            <ExploreSkeleton 
                resultSet={visibleResults} 
                scope={currentScope}
                clearScope={clearScope}
            />
        </div>
    </div>
</div> 
