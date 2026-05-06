<script lang="ts">
    import { onDestroy } from "svelte"
    import { createSocket } from "$lib/socket.svelte";
    import Sankey from "$lib/components/sankey/Sankey.svelte";
    import ExploreSkeleton from "$lib/components/exploratory/ExploreSkeleton.svelte";
    import type { Scope } from "$lib/types/types";
    import Toolbar from "./Toolbar.svelte";

    const socket = createSocket("ws://localhost:1323/ws")

    let currentScope = $state<Scope>({kind: "global", id: null})

    const setScope = ((newScope: Scope) => {
        if ((newScope.kind === currentScope.kind) && (newScope.id === currentScope.id)) {
            currentScope = {kind: "global", id: null};
            return;
        }
        currentScope = newScope;
    })

    
    const clearScope = () => {
        currentScope = {kind: "global", id: null}
    }

    let visibleResults = $derived(socket.auctionResults)

    onDestroy(() => {
        socket.disconnect()
    })

</script>

<div class="flex flex-col my-4">
    <Toolbar {socket} />
    <div class="flex flex-col lg:flex-row justify-center">
        <div class="card rounded-none flex-1 min-w-0 preset-filled-surface-100-900 border border-surface-200-800 divide-surface-200-800 py-6 px-6 flex flex-col">
            <h2 class="h5 mb-2">Auction Flow</h2>
            <Sankey 
                visibleResults={visibleResults} 
                setScope={setScope} 
                scope={currentScope}
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
