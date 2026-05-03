<script>
    import KPIDefault from "./KPIDefault.svelte";
    let { visibleResults } = $props() 
</script>

<div>
    <KPIDefault {visibleResults} />
</div>