<script lang="ts">
	import './layout.css';

	let { children } = $props();
</script>

<svelte:head></svelte:head>
<div class="container mx-auto px-2">
	{@render children()}
	<footer class="border-t border-surface-200-800 mt-8 py-6 px-2">
		<div class="flex flex-col sm:flex-row items-center justify-between gap-2 text-xs opacity-60">
			<span>Programmatic Sim Viz &middot; CS-GY 6313 &middot; 2026</span>
			<span class="hidden sm:inline">
				Built with 
				<a class="hover:opacity-80" href="https://svelte.dev/" target="_blank" rel="noopener">
					Svelte
				</a> 
				&amp; 
				<a class="hover:opacity-80" href="https://d3js.org" target="_blank" rel="noopener">
					D3
				</a>
			</span>
			<div class="flex items-center gap-4">
				<a
					href="https://github.com/mike-zaharchenko17/programmatic-sim-viz"
					class="hover:opacity-80 transition-opacity"
					target="_blank"
					rel="noopener"
				>
					GitHub
				</a>
				<span>&middot;</span>
				<a
					href="/about"
					class="hover:opacity-80 transition-opacity"
				>
					About
				</a>
			</div>
		</div>
	</footer>
</div>