<script lang="ts">
    import { sankey, sankeyLeft } from "d3-sankey";
    import type { NodeTween, InputNode, InputLink } from "$lib/types/types"
    import { auctionResultsToSankeyLinks, linksToNodes } from "$lib/data-processing-sankey";
    import { scaleOrdinal, schemeTableau10 } from "d3";
    import SankeyNodeTween from "./SankeyNodeTween.svelte";

    let { visibleResults } = $props()

    const width = 1000
    const height = 600

    const layout = sankey<InputNode, InputLink>()
        .nodeId(d => d.id)
        .nodeAlign(sankeyLeft)
        .nodeWidth(16)
        .nodePadding(12)
        .nodeSort(null)
        .linkSort(null)
        .iterations(4)
        .extent([
            [0, 0],
            [width, height]
        ])

    const nodeTweens = new Map<string, NodeTween>();
    
    let graph = $derived.by(() => {
        const links = auctionResultsToSankeyLinks(visibleResults)

        if (links.length === 0) {
            return null
        }

        const nodes = linksToNodes(links)

        return layout({
            nodes: nodes.map((d) => ({ ...d })),
            links: links.map((d) => ({ ...d }))
        })
    })

    const color = scaleOrdinal<string, string>(schemeTableau10)

</script>

<svg 
    {width} 
    {height} 
    role="img" 
    aria-label="Auction Sankey"
    style="outline: 1px solid red; display: block"
>
    {#if graph}
        <SankeyNodeTween graph={graph} nodeTweens={nodeTweens} {width} {color} />
    {/if}
</svg>