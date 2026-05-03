<script>
    let { visibleResults } = $props()

    let numAuctions = $derived(visibleResults.length)
</script>

<div>
    <p>Num Auctions: {numAuctions}</p>
</div>
