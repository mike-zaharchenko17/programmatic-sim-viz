<script lang="ts">
    import { Tween } from "svelte/motion";
    import { cubicOut } from "svelte/easing";
    import { scaleOrdinal, schemeTableau10 } from "d3";
    import type { SankeyGraph } from "d3-sankey";
    import type { InputNode, InputLink } from "$lib/types/types";

    const defaultColor = scaleOrdinal<string, string>(schemeTableau10)

    let { graph, nodeTweens, width, color = defaultColor }: {
        graph: SankeyGraph<InputNode, InputLink>;
        nodeTweens: Map<string, NodeTween>;
        width: number;
        color?: (value: string) => string;
    } = $props()

    // coordinate point tween object
    type NodeTween = {
        x0: Tween<number>;
        y0: Tween<number>;
        x1: Tween<number>;
        y1: Tween<number>;
    }

    type LinkTween = {
        y0: Tween<number>;
        y1: Tween<number>;
        width: Tween<number>;
    }

    const linkTweens = new Map<string, LinkTween>();

    function nodeId(endpoint: string | InputNode): string {
        return typeof endpoint === "string" ? endpoint : endpoint.id;
    }

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

</script>

{#if graph}
  <g class="links" fill="none" stroke-opacity="0.55">
    {#each graph.links as link (linkKey(link))}
      {@const source = nodeTweens.get(nodeId(link.source))}
      {@const target = nodeTweens.get(nodeId(link.target))}
      {@const lt = linkTweens.get(linkKey(link))}
      {#if source && target && lt}
        <path
          d={linkPath(source.x1.current, target.x0.current, lt.y0.current, lt.y1.current)}
          stroke={color(nodeId(link.source))}
          stroke-width={Math.max(1, lt.width.current)}
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
      <g class="node">
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