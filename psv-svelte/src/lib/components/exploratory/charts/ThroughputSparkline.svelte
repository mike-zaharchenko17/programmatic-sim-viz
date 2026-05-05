<script lang="ts">
    import type { AuctionResult, Scope } from "$lib/types/types";
    import { scaleLinear, line, max as d3Max, curveMonotoneX } from "d3";
    import { inScope } from "$lib/utils/scope";
    let { resultSet, scope }: {
        resultSet: AuctionResult[];
        scope: Scope;
    } = $props();

    const WINDOW_SEC = 60;
    const BIN_MS = 1000;

    let now = $state(Date.now());
    
    $effect(() => {
        const id = setInterval(() => (now = Date.now()), 1000);
        return () => clearInterval(id);
    });

    let bins = $derived.by(() => {
        const windowStart = now - WINDOW_SEC * 1000;
        const counts = new Array(WINDOW_SEC).fill(0);

        for (const r of resultSet) {
            if (!inScope(r, scope)) continue;
            const t = new Date(r.timestamp).getTime();
            if (t < windowStart || t > now) continue;
            const idx = Math.floor((t - windowStart) / BIN_MS);
            if (idx >= 0 && idx < WINDOW_SEC) counts[idx]++;
        }

        return counts;
    });

    let totalInWindow = $derived(bins.reduce((a, b) => a + b, 0));

    const W = 280;
    const H = 64;
    const PAD = 4;

    let x = $derived(
        scaleLinear()
            .domain([0, WINDOW_SEC - 1])
            .range([PAD, W - PAD])
    );

    let y = $derived(
        scaleLinear()
            .domain([0, Math.max(1, d3Max(bins) ?? 1)])
            .range([H - PAD, PAD])
    );

    let lineGen = $derived(
        line<number>()
            .x((_, i) => x(i))
            .y((d) => y(d))
            .curve(curveMonotoneX)
    );

    let pathD = $derived(lineGen(bins) ?? "");
    let areaD = $derived(
        pathD
            ? `${pathD} L ${x(WINDOW_SEC - 1)} ${H - PAD} L ${x(0)} ${H - PAD} Z`
            : ""
    );
</script>

<div class="flex justify-between text-xs opacity-70 mb-1">
    <span>last {WINDOW_SEC}s</span>
    <span>{totalInWindow} auctions</span>
</div>

<svg
    viewBox={`0 0 ${W} ${H}`}
    style:width="100%"
    style:height="auto"
    style:display="block"
    role="img"
    aria-label="Auctions per second over the last {WINDOW_SEC} seconds"
>
    <path d={areaD} fill="currentColor" fill-opacity="0.15" />
    <path
        d={pathD}
        fill="none"
        stroke="currentColor"
        stroke-width="1.5"
        stroke-linejoin="round"
        stroke-linecap="round"
    />
</svg>
