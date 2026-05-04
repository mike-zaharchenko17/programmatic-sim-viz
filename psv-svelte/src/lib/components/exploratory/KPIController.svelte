<script>
    import Pie from "./charts/Pie.svelte";
    let { visibleResults } = $props() 
</script>

<div>
    <Pie resultSet={visibleResults} />
</div>