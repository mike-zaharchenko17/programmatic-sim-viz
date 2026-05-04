<script lang="ts">
    import type { AuctionResult, Scope } from "$lib/types/types";
    import Pie from "./charts/Pie.svelte";
    let { resultSet, scope = { kind: "global" } }: {
        resultSet: AuctionResult[];
        scope?: Scope;
    } = $props();

    let scopedResults = $derived.by(() => {
        if (scope.kind === "global") return resultSet

        if (scope.kind === "seat") {
            return resultSet.filter((r: AuctionResult) => r.winner?.adomain?.[0] === scope.id)
        }

        if (scope.kind === "campaign") {
            return resultSet.filter((r: AuctionResult) => r.winner?.cid === scope.id)
        }
    })
    
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
        <Pie resultSet={scopedResults} />
    </section>
</div>



