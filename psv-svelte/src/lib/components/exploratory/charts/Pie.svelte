
<script lang="ts">
    import type { AuctionResult } from "$lib/types/types"
    import { pie, arc, type PieArcDatum } from "d3";
    import { Tween } from "svelte/motion";
    import { cubicOut } from "svelte/easing";
    let { resultSet, scope } = $props()

    // filters
    let { nw: countWinningBids, nl: countLosingBids, na: countTotalBids } = $derived.by(() : { nw: number; nl: number, na: number } => {
        let nw, nl, na = 0

        if (scope.kind === "global") {
            nw = resultSet.filter((res: AuctionResult) => res.winner !== undefined).length
            nl = resultSet.reduce(
                (acc: number, res: AuctionResult) => acc + (res.losers?.length ?? 0),
                0
            )
        }

        if (scope.kind === "seat") {
            nw = resultSet.filter((r: AuctionResult) => r.winner?.adomain?.[0] === scope.id).length
            nl = resultSet.filter((r: AuctionResult) => {
                // case won: exclude from loser count
                if (r.winner?.adomain?.[0] === scope.id) {
                    return false
                }

                // case didn't win: see whether appears in loser array. if so, count 
                if (r.losers) {
                    for (const loss of r.losers) {
                        if (loss.bid.adomain?.[0] === scope.id) {
                            return true
                        }
                    }
                }
                // case: neither won nor lost; didn't bid
                return false
            }).length
        }

        if (scope.kind === "campaign") {
            nw = resultSet.filter((r: AuctionResult) => 
                r.winner?.cid ? r.winner.cid === scope.id : false
            ).length

            nl = resultSet.filter((r: AuctionResult) => {
                // case won: exclude from loser count
                if (r.winner?.cid && r.winner.cid === scope.id) {
                    return false
                }

                // case didn't win: see whether appears in loser array. if so, count 
                if (r.losers) {
                    for (const loss of r.losers) {
                        if (loss.bid.cid && loss.bid.cid === scope.id) {
                            return true
                        }
                    }
                }
                // case: neither won nor lost; didn't bid
                return false
            }).length
        }

        if (na === 0) {
            na = nw + nl
        }

        return {nw, nl, na}
    })


    type Slice = { label: string; value: number; color: string };

    let data = $derived<Slice[]>([
        { label: "Won",  value: countWinningBids, color: "#55cc54" },
        { label: "Lost", value: countLosingBids,  color: "#fa594d" },
    ])

    const size = 240
    const radius = size / 2
    const pieLayout = pie<Slice>().sort(null).value((d) => d.value)

    const arcGen = arc<PieArcDatum<Slice>>()
        .innerRadius(radius * 0.6)
        .outerRadius(radius - 8)

    let arcs = $derived(pieLayout(data))

    type ArcTween = {
        startAngle: Tween<number>;
        endAngle: Tween<number>;
    };
    const arcTweens = new Map<string, ArcTween>();

    $effect(() => {
        for (const a of arcs) {
            let t = arcTweens.get(a.data.label);
            if (!t) {
                t = {
                    startAngle: new Tween(a.startAngle, { duration: 300, easing: cubicOut }),
                    endAngle: new Tween(a.endAngle, { duration: 300, easing: cubicOut }),
                };
                arcTweens.set(a.data.label, t);
            } else {
                t.startAngle.target = a.startAngle;
                t.endAngle.target = a.endAngle;
            }
        }
    });
</script>

<svg
  viewBox={`0 0 ${size} ${size}`}
  class="w-1/2 lg:w-11/12 h-auto block mx-auto"
>
    <g transform={`translate(${radius}, ${radius})`}>
        {#each arcs as a (a.data.label)}
            {@const t = arcTweens.get(a.data.label)}
            {#if t}
                {@const tweened = {
                    ...a,
                    startAngle: t.startAngle.current,
                    endAngle: t.endAngle.current,
                }}
                <path d={arcGen(tweened)} fill={a.data.color} />
                {#if tweened.endAngle - tweened.startAngle > 0.15}
                    {@const [cx, cy] = arcGen.centroid(tweened)}
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
            {/if}
        {/each}
        <text
            text-anchor="middle"
            dominant-baseline="middle"
            fill="currentColor"
            style:pointer-events="none"
        >
            {#if resultSet.length > 0}
                <tspan x="0" y="0" dy="-0.3em" font-size="11">Total Bids</tspan>
                <tspan x="0" dy="1.2em" font-size="16" font-weight="600">{countTotalBids}</tspan>
            {:else}
                <tspan class="text-xs opacity-50 italic ">No results yet</tspan>
            {/if}
        </text>
    </g>
</svg>
