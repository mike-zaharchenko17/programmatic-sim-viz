<script lang="ts">
    import { sankey, sankeyLeft } from "d3-sankey";
    import type { NodeTween, InputNode, InputLink } from "$lib/types/types"
    import { auctionResultsToSankeyLinks, linksToNodes } from "$lib/data-processing-sankey";
    import { scaleOrdinal, schemeObservable10 } from "d3";
    import { hcl } from "d3";
    import SankeyNodeTween from "./SankeyNodeTween.svelte";

    let { visibleResults } = $props()

    // removes red and green from the tableau scale
    const SEAT_PALETTE_FILTERED = schemeObservable10.filter((_, i) => i !== 2 && i !== 4)

    const width = 900
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

    let links = $derived(auctionResultsToSankeyLinks(visibleResults))
    
    let graph = $derived.by(() => {
        if (links.length === 0) {
            return null
        }

        const nodes = linksToNodes(links)

        return layout({
            nodes: nodes.map((d) => ({ ...d })),
            links: links.map((d) => ({ ...d }))
        })
    })

    function isCampaign(id: string) {
        return id.startsWith("Campaign:") || id.startsWith("Creative:")
    }

    function isOutcome(id: string) {
        return id === "Lost: Outbid" || id === "Won"
    }

    function endpointId(e: string | InputNode) {
        return typeof e === "string" ? e : e.id
    }

    function shade(seatColor: string, idx: number, total: number): string {
        const c = hcl(seatColor);
        const lo = 40, hi = 78;
        const t = total > 1 ? idx / (total - 1) : 0.5;
        c.l = lo + t * (hi - lo);
        return c.formatHex();
    }

    let color = $derived.by(() => {
        const fallback = (_: string) => "#888";
        if (!graph) return fallback;

        const seats = new Set<string>();
        const campaignToSeat = new Map<string, string>();
        const seatToCampaigns = new Map<string, string[]>();

        for (const link of graph.links) {
            const s = endpointId(link.source as any);
            const t = endpointId(link.target as any);
            if (isCampaign(t) && !isOutcome(s)) {
                seats.add(s);
                if (!campaignToSeat.has(t)) {
                    campaignToSeat.set(t, s);
                    const arr = seatToCampaigns.get(s) ?? [];
                    arr.push(t);
                    seatToCampaigns.set(s, arr);
                }
            } else if (!isCampaign(s) && !isOutcome(s) && !isOutcome(t)) {
                seats.add(s);
            }
        }

        const seatScale = scaleOrdinal<string, string>()
            .domain([...seats].sort())
            .range(SEAT_PALETTE_FILTERED);

        const shadeIdx = new Map<string, number>();

        for (const [, arr] of seatToCampaigns) {
            arr.sort();
            arr.forEach((c, i) => shadeIdx.set(c, i));
        }

        return (id: string): string => {
            if (id === "Won") return "#55cc54";
            if (id === "Lost: Outbid") return "#fa594d";
            const seat = campaignToSeat.get(id);
            if (seat !== undefined) {
                const total = seatToCampaigns.get(seat)!.length;
                return shade(seatScale(seat), shadeIdx.get(id)!, total);
            }
            return seatScale(id);
        };
    });


</script>

<svg 
    viewBox={`0 0 ${width} ${height}`}
    preserveAspectRatio="xMidYMid meet"
    role="img" 
    aria-label="Auction Sankey"
    style:width="100%"
    style:height="auto"
    style:display="block"
>
    {#if graph}
        <SankeyNodeTween 
            graph={graph} 
            nodeTweens={nodeTweens} 
            width={width} 
            color={color} 
        />
    {/if}
</svg>