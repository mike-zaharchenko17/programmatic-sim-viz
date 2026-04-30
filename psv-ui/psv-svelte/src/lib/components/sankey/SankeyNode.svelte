<script lang="ts">
    import type { InputNode } from "$lib/types/types"
    let { link, path, color } = $props()
</script>

<path
    class="link"
    d={path}
    stroke={color}
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


