<script lang="ts">
    import { Tween } from "svelte/motion";
    import { cubicOut } from "svelte/easing";
    import { scaleOrdinal, schemeTableau10 } from "d3";
    import type { SankeyGraph } from "d3-sankey";
    import { isCampaign, isOutcome, nodeId } from "$lib/utils/node-classification";
    import type { InputNode, InputLink, LinkTween, NodeTween, Scope } from "$lib/types/types";

    const defaultColor = scaleOrdinal<string, string>(schemeTableau10)

    let { graph, nodeTweens, width, handleNodeClick, color = defaultColor, }: {
        graph: SankeyGraph<InputNode, InputLink>;
        nodeTweens: Map<string, NodeTween>;
        width: number;
        handleNodeClick: (v: Scope) => void;
        color?: (value: string) => string;
    } = $props()

    const linkTweens = new Map<string, LinkTween>();

    function linkKey(link: any): string {
        return `${nodeId(link.source)}->${nodeId(link.target)}`;
    }

    function getNodeTween(node: any): NodeTween {
        let t = nodeTweens.get(node.id);
        if (!t) {
            t = {
                x0: new Tween(node.x0, { duration: 400, easing: cubicOut }),
                y0: new Tween(node.y0, { duration: 400, easing: cubicOut }),
                x1: new Tween(node.x1, { duration: 400, easing: cubicOut }),
                y1: new Tween(node.y1, { duration: 400, easing: cubicOut }),
            };
            nodeTweens.set(node.id, t);
        } else {
            t.x0.target = node.x0
            t.y0.target = node.y0
            t.x1.target = node.x1
            t.y1.target = node.y1
        }
        return t;
    }

    function getLinkTween(link: any): LinkTween {
        const key = linkKey(link);
        let t = linkTweens.get(key);

        if (!t) {
            t = {
                y0: new Tween(link.y0 ?? 0, { duration: 400, easing: cubicOut }),
                y1: new Tween(link.y1 ?? 0, { duration: 400, easing: cubicOut }),
                width: new Tween(link.width ?? 1, { duration: 400, easing: cubicOut })
            };
            linkTweens.set(key, t);
        } else {
            t.y0.target = link.y0 ?? 0;
            t.y1.target = link.y1 ?? 0;
            t.width.target = link.width ?? 1;
        }

        return t;
    }

    function linkPath(sourceX: number, targetX: number, y0: number, y1: number): string {
        const midX = (sourceX + targetX) / 2;
        return `M${sourceX},${y0}C${midX},${y0} ${midX},${y1} ${targetX},${y1}`;
    }

    let hovered = $state<string | null>(null);
    let pendingTimer: ReturnType<typeof setTimeout> | null = null;

    const HOVER_DELAY_MS = 400;

    function clearPending() {
        if (pendingTimer !== null) {
            clearTimeout(pendingTimer);
            pendingTimer = null;
        } 
    }

    function scheduleHover(id: string) {
        clearPending()
        pendingTimer = setTimeout(() => {
          hovered = id
          pendingTimer = null
        }, HOVER_DELAY_MS)
    }

    function cancelHover() {
        clearPending();
        hovered = null;
    }

    let adjacency = $derived.by(() => {
        const incomingByTarget = new Map<string, any[]>();
        const outgoingBySource = new Map<string, any[]>();
        if (!graph) return { incomingByTarget, outgoingBySource };

        for (const link of graph.links) {
            const s = nodeId(link.source as any);
            const t = nodeId(link.target as any);
            if (!incomingByTarget.has(t)) incomingByTarget.set(t, []);
            incomingByTarget.get(t)!.push(link);
            if (!outgoingBySource.has(s)) outgoingBySource.set(s, []);
            outgoingBySource.get(s)!.push(link);
        }
        return { incomingByTarget, outgoingBySource };
    });

    let related = $derived.by(() => {
        if (!graph || !hovered) return null;

        const { incomingByTarget, outgoingBySource } = adjacency;
        const nodes = new Set<string>([hovered]);
        const linkKeys = new Set<string>();

        const walk = (
            start: string,
            adj: Map<string, any[]>,
            getNext: (l: any) => string,
        ) => {
            const queue = [start];
            while (queue.length) {
                const n = queue.shift()!;
                for (const link of adj.get(n) ?? []) {
                    linkKeys.add(linkKey(link));
                    const next = getNext(link);
                    if (!nodes.has(next)) {
                        nodes.add(next);
                        queue.push(next);
                    }
                }
            }
        };

        walk(hovered, incomingByTarget, (l) => nodeId(l.source));
        walk(hovered, outgoingBySource, (l) => nodeId(l.target));

        return { nodes, linkKeys };
    });

    function nodeOpacity(id: string): number {
        if (!hovered) return 1;
        return related?.nodes.has(id) ? 1 : 0.15;
    }

    function linkOpacity(link: any): number {
        if (!hovered) return 0.55;
        return related?.linkKeys.has(linkKey(link)) ? 0.85 : 0.05;
    }

    $effect(() => {
        if (!graph) return;
        const liveNodeIds = new Set<string>();
        const liveLinkIds = new Set<string>();

        for (const n of graph.nodes) {
            liveNodeIds.add(n.id);
            getNodeTween(n);
        }

        for (const id of nodeTweens.keys()) {
            if (!liveNodeIds.has(id)) nodeTweens.delete(id);
        }

        for (const link of graph.links) {
            const key = linkKey(link);
            liveLinkIds.add(key);
            getLinkTween(link);
        }

        for (const id of linkTweens.keys()) {
            if (!liveLinkIds.has(id)) linkTweens.delete(id);
        }
    });

    $effect(() => () => clearPending())

</script>

{#if graph}
  <g class="links" fill="none">
    {#each graph.links as link (linkKey(link))}
      {@const source = nodeTweens.get(nodeId(link.source))}
      {@const target = nodeTweens.get(nodeId(link.target))}
      {@const lt = linkTweens.get(linkKey(link))}
      {#if source && target && lt}
        <path
          d={linkPath(source.x1.current, target.x0.current, lt.y0.current, lt.y1.current)}
          stroke={color(nodeId(link.source))}
          stroke-width={Math.max(1, lt.width.current)}
          stroke-opacity={linkOpacity(link)}
          style:transition="stroke-opacity 150ms ease-out"
        >
          <title>{nodeId(link.source)} → {nodeId(link.target)}: {link.value}</title>
        </path>
      {/if}
    {/each}
  </g>

  <g class="nodes">
  {#each graph.nodes as node (node.id)}
    {@const t = nodeTweens.get(node.id)}
    {#if t}
      <g
        class="node"
        opacity={nodeOpacity(node.id)}
        style:transition="opacity 150ms ease-out"
        style:cursor="pointer"
        onmouseenter={() => scheduleHover(node.id)}
        onmouseleave={cancelHover}
        onclick={() => {
          if (isCampaign(node.id)) {
            handleNodeClick({ kind: "campaign", id: node.id })
          }
          if (!isCampaign(node.id) && !isOutcome(node.id)) {
            handleNodeClick({ kind: "seat", id: node.id })
          }
        }}
        role="presentation"
      >
        <rect
          x={t.x0.current}
          y={t.y0.current}
          width={t.x1.current - t.x0.current}
          height={t.y1.current - t.y0.current}
          fill={color(node.id)}
        >
          <title>{node.id}: {node.value}</title>
        </rect>
        <text
          x={t.x0.current < width / 2 ? t.x1.current + 6 : t.x0.current - 6}
          y={(t.y0.current + t.y1.current) / 2}
          dy="0.35em"
          text-anchor={t.x0.current < width / 2 ? "start" : "end"}
          dominant-baseline="middle"
          font-size="11"
          fill="#F2F0EF"
        >
          {node.id}
        </text>
      </g>
    {/if}
  {/each}
  </g>
{/if}