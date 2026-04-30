<script lang="ts">
    let { node, width } = $props()
</script>

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