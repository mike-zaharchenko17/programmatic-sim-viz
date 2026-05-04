<script>
    import Pie from "./charts/Pie.svelte";
    let { visibleResults } = $props() 
</script>

<div>
    <Pie {visibleResults} />
</div>