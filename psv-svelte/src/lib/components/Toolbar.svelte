<script lang="ts">
    import type { Socket } from "$lib/socket.svelte";

    let { socket }: { socket: Socket } = $props()

    let count = $derived(socket.auctionResults.length)
    let clearIsDisabled = $derived(count === 0)
</script>

<div
    class="card rounded-none border border-b-0 border-surface-200-800 preset-filled-surface-100-900 flex items-center justify-between gap-4 px-4 py-2 w-full"
    role="toolbar"
    aria-label="Dashboard controls"
>
    <div class="flex items-center gap-3 text-sm">
        <span class="relative flex w-2 h-2 shrink-0" aria-hidden="true">
            {#if socket.isOpen}
                <span class="absolute inset-0 rounded-full bg-green-500 opacity-60 animate-ping"></span>
            {/if}
            <span
                class="relative inline-block w-2 h-2 rounded-full {socket.isOpen ? 'bg-green-500' : 'bg-zinc-500'}"
            ></span>
        </span>
        <span class="opacity-80">
            {socket.isOpen ? 'Live' : 'Disconnected'}
        </span>
        {#if count > 0}
            <span class="opacity-30">·</span>
            <span class="tabular-nums opacity-70">
                {count.toLocaleString()} {count === 1 ? 'auction' : 'auctions'}
            </span>
        {/if}
    </div>

    <div class="flex items-center gap-2">
        <button 
            type="button" 
            class="btn btn-sm preset-filled"
            disabled={clearIsDisabled}
            onclick={socket.clear}
        >
            Clear
        </button>
        {#if socket.isOpen}
            <button type="button" class="btn btn-sm preset-tonal" onclick={socket.disconnect}>
                Disconnect
            </button>
        {:else}
            <button type="button" class="btn btn-sm preset-filled" onclick={socket.connect}>
                Connect
            </button>
        {/if}
    </div>
</div>



