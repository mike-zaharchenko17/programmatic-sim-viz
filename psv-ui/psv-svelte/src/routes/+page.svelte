<script lang="ts">
    import { onDestroy, onMount } from "svelte"
    import { createSocket } from "$lib/socket.svelte";
    import { sankey, sankeyLeft, sankeyLinkHorizontal, type SankeyGraph } from "d3-sankey";
    import type { InputNode, InputLink } from "$lib/types/types";
    import { schemeTableau10 } from "d3";
    import { scaleOrdinal } from "d3";
    import { auctionResultsToSankeyLinks, linksToNodes } from "$lib/data-processing-sankey";

    const socket = createSocket("ws://localhost:1323/ws")

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

    let visibleResults = $derived(socket.auctionResults.slice(-30))
    
    let graph = $derived.by(() => {
        const links = auctionResultsToSankeyLinks(visibleResults)

        if (links.length === 0) {
            return null
        }

        const nodes = linksToNodes(links);

        return layout({
            nodes: nodes.map((d) => ({ ...d })),
            links: links.map((d) => ({ ...d }))
        });
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

    onMount(() => {
        socket.connect()
    })

    onDestroy(() => {
        socket.disconnect()
    })

</script>

<div class="container">
    <div class="container">
        <div class="toolbar">

          {#if socket.isOpen}
            <button onclick={socket.disconnect}>Disconnect</button>
          {:else}
            <button onclick={socket.connect}>Connect</button>
          {/if}

          <svg {width} {height} role="img" aria-label="Auction Sankey">
            {#if graph}
                <g fill="none" stroke-opacity="0.4">
                    {#each graph.links as link (linkKey(link))}
                        <path
                            class="link"
                            d={linkPath(link) ?? ""}
                            stroke={color(linkSourceId(link))}
                            stroke-width={Math.max(1, link.width ?? 1)}
                        >
                            <title>
                                {typeof link.source === "string" ? link.source : (link.source as InputNode).id}
                                →
                                {typeof link.target === "string" ? link.target : (link.target as InputNode).id}
                                :
                                {link.value}
                            </title>
                        </path>
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

        </div>
    </div>    
</div>

