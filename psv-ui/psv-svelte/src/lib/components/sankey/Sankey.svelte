<script lang="ts">
    import { sankey, sankeyLeft, sankeyLinkHorizontal, type SankeyGraph } from "d3-sankey";
    import type { InputNode, InputLink } from "$lib/types/types"
    import { auctionResultsToSankeyLinks, linksToNodes } from "$lib/data-processing-sankey";
    import { scaleOrdinal, schemeTableau10 } from "d3";
    import SankeyLink from "./SankeyLink.svelte";

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
                <g class="node">
                    <rect
                        x={node.x0}
                        y={node.y0}
                        width={(node.x1 ?? 0) - (node.x0 ?? 0)}
                        height={(node.y1 ?? 0) - (node.y0 ?? 0)}
                    >
                        <title>{node.id}: {node.value}</title>
                    </rect>
                    <text
                        x={(node.x0 ?? 0) < width / 2 ? (node.x1 ?? 0) + 6 : (node.x0 ?? 0) - 6}
                        y={((node.y0 ?? 0) + (node.y1 ?? 0)) / 2}
                        dy="0.35em"
                        text-anchor={(node.x0 ?? 0) / 2 ? "start" : "end"}
                        dominant-baseline="middle"
                        font-size="11"
                    >
                        {node.id}
                    </text>
                </g>
            {/each}
        </g>
    {/if}
</svg>