<script lang="ts">
    import type { Socket } from "$lib/socket.svelte";

    let { socket }: { socket: Socket } = $props()
</script>

{#if socket.isOpen}
    <button type="button" class="btn preset-filled" onclick={socket.disconnect}>Disconnect</button>
{:else}
    <button type="button" class="btn preset-filled" onclick={socket.connect}>Connect</button>
{/if}



