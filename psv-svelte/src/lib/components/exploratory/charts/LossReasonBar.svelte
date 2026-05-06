<script lang="ts">
    import type { AuctionResult, Scope } from "$lib/types/types";
    import { hcl } from "d3";

    let { resultSet, scope }: {
        resultSet: AuctionResult[];
        scope: Scope;
    } = $props();

    function lossReasonLabel(code: number): string {
        if (code === 2) return "Impression Expired";
        if (code === 100) return "Below Floor";
        return "Outbid";
    }

    // ordering controls left-to-right in the bar and top-to-bottom in the legend
    const REASONS = ["Outbid", "Below Floor", "Impression Expired"] as const;
    type Reason = typeof REASONS[number];

    function shadesOf(base: string, n: number): string[] {
        const out: string[] = [];
        const lo = 45, hi = 75;
        for (let i = 0; i < n; i++) {
            const c = hcl(base);
            const t = n > 1 ? i / (n - 1) : 0.5;
            c.l = lo + t * (hi - lo);
            out.push(c.formatHex());
        }
        return out;
    }

    const SHADES = shadesOf("#fa594d", REASONS.length);
    const REASON_COLORS: Record<Reason, string> = Object.fromEntries(
        REASONS.map((r, i) => [r, SHADES[i]])
    ) as Record<Reason, string>;

    function lossMatchesScope(bid: AuctionResult["winner"]): boolean {
        if (!bid) return false;
        if (scope.kind === "global") return true;
        if (scope.kind === "seat") return bid.adomain?.[0] === scope.id;
        if (scope.kind === "campaign") return bid.cid === scope.id;
        return false;
    }

    let counts = $derived.by(() => {
        const c: Record<Reason, number> = {
            "Outbid": 0,
            "Below Floor": 0,
            "Impression Expired": 0,
        };
        for (const r of resultSet) {
            for (const loser of r.losers ?? []) {
                if (!lossMatchesScope(loser.bid)) continue;
                c[lossReasonLabel(loser.loss_reason) as Reason]++;
            }
        }
        return c;
    });

    let total = $derived(
        (Object.values(counts) as number[]).reduce((a, b) => a + b, 0)
    );

    const W = 280;
    const H = 16;

    let segments = $derived.by(() => {
        let acc = 0;
        return REASONS.map((r) => {
            const value = counts[r];
            const pct = total > 0 ? value / total : 0;
            const w = pct * W;
            const seg = {
                label: r,
                value,
                color: REASON_COLORS[r],
                pct,
                x: acc,
                w,
            };
            acc += w;
            return seg;
        });
    });
</script>

{#if total > 0}
    <svg
        viewBox={`0 0 ${W} ${H}`}
        style:width="100%"
        style:height="auto"
        style:display="block"
        role="img"
        aria-label="Loss reason breakdown"
    >
        {#each segments as s (s.label)}
            {#if s.w > 0}
                <rect x={s.x} y={0} width={s.w} height={H} fill={s.color} />
            {/if}
        {/each}
    </svg>

    <ul class="flex flex-col gap-0.5 mt-2 text-xs">
        {#each segments as s (s.label)}
            <li class="flex items-center gap-2">
                <span
                    class="inline-block w-2 h-2 rounded-sm shrink-0"
                    style:background-color={s.color}
                    aria-hidden="true"
                ></span>
                <span class="flex-1 opacity-80">{s.label}</span>
                <span class="opacity-70 tabular-nums">{s.value}</span>
                <span class="opacity-50 tabular-nums w-10 text-right">
                    {(s.pct * 100).toFixed(0)}%
                </span>
            </li>
        {/each}
    </ul>
{:else}
    <div class="text-xs opacity-50 italic py-2">No losses yet</div>
{/if}
