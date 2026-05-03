<script lang="ts">
    let { visibleResults } = $props()
    import type { AuctionResult } from "$lib/types/types";

    let numAuctions = $derived(visibleResults.length)
    let numWinners = $derived(visibleResults.filter((res: AuctionResult) => res.winner !== undefined).length)
    let numLosers = $derived(visibleResults.reduce((acc: number, res: AuctionResult) => res.losers ? acc + res.losers.length : acc += 0, 0))
</script>

<div>
    <p>Num Auctions: {numAuctions}</p>
    <p>Num Winners: {numWinners}</p>
    <p>Num Losers: {numLosers}</p>
</div>
