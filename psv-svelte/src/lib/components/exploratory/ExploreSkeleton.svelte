<script lang="ts">
    import type { AuctionResult, Scope } from "$lib/types/types";
    import Pie from "./charts/Pie.svelte";
    let { resultSet, scope = { kind: "global" } }: {
        resultSet: AuctionResult[];
        scope?: Scope;
    } = $props();
    
    let title = $derived(scope.kind === "global" ? "Global" : scope.id)
</script>

<div class="flex flex-col gap-4 h-full">
    <header class="flex items-center justify-between">
        <h2 class="h6">{title}</h2>
        {#if scope.kind !== "global"}
            <button class="btn btn-sm">x</button>
        {/if}
    </header>
    <section aria-label="Outcome share">
        <Pie resultSet={resultSet} scope={scope} />
    </section>
</div>



