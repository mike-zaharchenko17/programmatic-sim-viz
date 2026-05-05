<script lang="ts">
    import type { AuctionResult, Scope } from "$lib/types/types";
    import { wonByScope } from "$lib/utils/scope";
    import { bin, extent, max as d3Max, median, scaleLinear } from "d3";

    let { resultSet, scope }: {
        resultSet: AuctionResult[];
        scope: Scope;
    } = $props();


    let prices = $derived.by(() => {
        const out: number[] = [];
        for (const r of resultSet) {
            if (wonByScope(r, scope)) out.push(r.clearing_price);
        }
        return out;
    });

    const MIN_POINTS = 5;
    const BIN_COUNT = 20;

    const W = 280;
    const H = 80;
    const PAD = { top: 4, right: 2, bottom: 12, left: 2 };

    let domain = $derived.by((): [number, number] | null => {
        const [lo, hi] = extent(prices);
        if (lo === undefined || hi === undefined) return null;
        if (lo === hi) return [lo - 0.5, hi + 0.5];
        return [lo, hi];
    });

    let bins = $derived.by(() => {
        if (!domain || prices.length < MIN_POINTS) return null;
        return bin<number, number>()
            .domain(domain)
            .thresholds(BIN_COUNT)(prices);
    });

    let med = $derived(prices.length >= MIN_POINTS ? median(prices) : null);

    let xScale = $derived(
        domain
            ? scaleLinear().domain(domain).range([PAD.left, W - PAD.right])
            : null
    );

    let yScale = $derived.by(() => {
        if (!bins) return null;
        const maxCount = d3Max(bins, (b) => b.length) ?? 1;
        return scaleLinear()
            .domain([0, Math.max(1, maxCount)])
            .range([H - PAD.bottom, PAD.top]);
    });

    function fmtPrice(n: number): string {
        return `$${n.toFixed(2)}`;
    }
</script>

<div class="flex justify-between text-xs opacity-70 mb-1">
    <span>{prices.length} wins</span>
    {#if med !== null && med !== undefined}
        <span>med {fmtPrice(med)}</span>
    {/if}
</div>

{#if bins && xScale && yScale && domain}
    <svg
        viewBox={`0 0 ${W} ${H}`}
        style:width="100%"
        style:height="auto"
        style:display="block"
        role="img"
        aria-label="Clearing price distribution"
    >
        {#each bins as b (b.x0)}
            {#if b.x0 !== undefined && b.x1 !== undefined && b.length > 0}
                {@const bx = xScale(b.x0)}
                {@const bw = Math.max(0, xScale(b.x1) - xScale(b.x0) - 1)}
                {@const by = yScale(b.length)}
                {@const bh = H - PAD.bottom - by}
                <rect
                    x={bx}
                    y={by}
                    width={bw}
                    height={bh}
                    fill="currentColor"
                    fill-opacity="0.7"
                />
            {/if}
        {/each}

        {#if med !== null && med !== undefined}
            {@const mx = xScale(med)}
            <line
                x1={mx}
                x2={mx}
                y1={PAD.top}
                y2={H - PAD.bottom}
                stroke="currentColor"
                stroke-width="1"
                stroke-dasharray="2 2"
                opacity="0.6"
            />
        {/if}

        <text
            x={PAD.left}
            y={H - 2}
            font-size="9"
            fill="currentColor"
            opacity="0.6"
        >
            {fmtPrice(domain[0])}
        </text>
        <text
            x={W - PAD.right}
            y={H - 2}
            font-size="9"
            fill="currentColor"
            opacity="0.6"
            text-anchor="end"
        >
            {fmtPrice(domain[1])}
        </text>
    </svg>
{:else}
    <div class="text-xs opacity-50 italic py-2">
        {prices.length === 0 ? "No wins yet" : "Need a few more wins to plot"}
    </div>
{/if}
