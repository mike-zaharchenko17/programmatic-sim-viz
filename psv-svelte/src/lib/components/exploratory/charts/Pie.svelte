
<script lang="ts">
    import type { AuctionResult } from "$lib/types/types"
    import { pie, arc, type PieArcDatum } from "d3";
    let { resultSet } = $props()
    let numAuctions = $derived(resultSet.length)
    let numWinners = $derived(
        resultSet.filter((res: AuctionResult) => res.winner !== undefined).length
    )
    let numLosers = $derived(
        resultSet.reduce(
            (acc: number, res: AuctionResult) => acc + (res.losers?.length ?? 0),
            0
        )
    )

    type Slice = { label: string; value: number; color: string };

    let data = $derived<Slice[]>([
        { label: "Won",  value: numWinners, color: "#55cc54" },
        { label: "Lost", value: numLosers,  color: "#fa594d" },
    ])

    const size = 240
    const radius = size / 2
    const pieLayout = pie<Slice>().sort(null).value((d) => d.value)

    const arcGen = arc<PieArcDatum<Slice>>()
        .innerRadius(radius * 0.6)
        .outerRadius(radius - 8)

    let arcs = $derived(pieLayout(data))
</script>

<svg
  viewBox={`0 0 ${size} ${size}`}
  style:width="100%"
  style:height="auto"
  style:display="block"
>
    <g transform={`translate(${radius}, ${radius})`}>
        {#each arcs as a (a.data.label)}
            <path d={arcGen(a)} fill={a.data.color} />
            {#if a.endAngle - a.startAngle > 0.15}
                {@const [cx, cy] = arcGen.centroid(a)}
                <text
                    text-anchor="middle"
                    dominant-baseline="middle"
                    fill="white"
                    font-size="12"
                    font-weight="600"
                    style:pointer-events="none"
                >
                    <tspan x={cx} y={cy} dy="-0.3em">{a.data.label}</tspan>
                    <tspan x={cx} dy="1.2em">{a.data.value}</tspan>
                </text>
            {/if}
        {/each}
        <text
            text-anchor="middle"
            dominant-baseline="middle"
            fill="currentColor"
            style:pointer-events="none"
        >
            <tspan x="0" y="0" dy="-0.3em" font-size="11">Total</tspan>
            <tspan x="0" dy="1.2em" font-size="16" font-weight="600">{numAuctions}</tspan>
        </text>
    </g>
</svg>
