<script lang="ts">
    import type { AuctionResult, Scope } from "$lib/types/types";
    import Pie from "./charts/Pie.svelte";
    import ThroughputSparkline from "./charts/ThroughputSparkline.svelte";
    import PriceHistogram from "./charts/PriceHistogram.svelte";
    import LossReasonBar from "./charts/LossReasonBar.svelte";
    let { resultSet, clearScope, scope = { kind: "global" } }: {
        resultSet: AuctionResult[];
        clearScope: () => void
        scope?: Scope;
    } = $props();
    
    let title = $derived(scope.kind === "global" ? "Global" : scope.id)
</script>

<div class="flex flex-col gap-4 h-full">
    <header class="flex items-center justify-between">
        <h2 class="h6">{title}</h2>
        {#if scope.kind !== "global"}
            <button class="btn btn-sm" onclick={() => clearScope()}>x</button>
        {/if}
    </header>
    <section aria-label="Outcome share">
        <Pie resultSet={resultSet} scope={scope} />
        <div class="border border-surface-200-800 my-2"></div>
    </section>
    <section aria-label="Loss reasons">
        <h3 class="text-xs uppercase opacity-70 mb-1">Loss reasons</h3>
        <LossReasonBar resultSet={resultSet} scope={scope} />
    </section>
    <section aria-label="Clearing price distribution">
        <h3 class="text-xs uppercase opacity-70 mb-1">Clearing price</h3>
        <PriceHistogram resultSet={resultSet} scope={scope} />
    </section>
    <section aria-label="Throughput">
        <h3 class="text-xs uppercase opacity-70 mb-1">Throughput</h3>
        <ThroughputSparkline resultSet={resultSet} scope={scope} />
    </section>
</div>



