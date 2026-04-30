<script lang="ts">
    import { sankey, sankeyLeft, sankeyLinkHorizontal, type SankeyGraph } from "d3-sankey";
    import type { InputNode, InputLink } from "$lib/types/types"
    import { auctionResultsToSankeyLinks, linksToNodes } from "$lib/data-processing-sankey";
    import { scaleOrdinal, schemeTableau10 } from "d3";
    import SankeyLink from "./SankeyLink.svelte";
    import SankeyNode from "./SankeyNode.svelte";

    let { visibleResults } = $props()

    const width = 900
    const height = 500

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

    const linkPath = sankeyLinkHorizontal<any, any>()

    const color = scaleOrdinal<string, string>(schemeTableau10)

    const linkKey = (link: any) : string => {
        const source = typeof link.source === "string" ? link.source : (link.source as InputNode).id;
        const target = typeof link.target === "string" ? link.target : (link.target as InputNode).id;

        return `${source}->${target}`;
    }

    const linkSourceId = (link: any): string =>
        typeof link.source === "string" ? link.source : (link.source as InputNode).id

</script>

<svg {width} {height} role="img" aria-label="Auction Sankey">
    {#if graph}
        <g class="links" fill="none" stroke-opacity="0.4">
            {#each graph.links as link (linkKey(link))}
                <SankeyLink
                    link={link}
                    path={linkPath(link) ?? ""}
                    color={color(linkSourceId(link))}
                />
            {/each}
        </g>
        <g class="nodes">
            {#each graph.nodes as node (node.id)}
                <SankeyNode
                    node={node}
                    width={width}
                />
            {/each}
        </g>
    {/if}
</svg>