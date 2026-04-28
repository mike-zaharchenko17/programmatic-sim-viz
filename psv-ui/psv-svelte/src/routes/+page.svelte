<!--
	RequestID     string       `json:"request_id"`
	Timestamp     time.Time    `json:"timestamp"`
	Winner        *Bid         `json:"winner,omitempty"`
	ClearingPrice float64      `json:"clearing_price"`
	Losers        []LossRecord `json:"losers,omitempty"`
-->

<script lang="ts">
    import { onDestroy } from "svelte"
    import { createSocket } from "$lib/socket.svelte";

    const socket = createSocket("ws://localhost:1323/ws")

    onDestroy(() => {
        socket.disconnect()
    })

</script>

<div class="container">

    {#if socket.isOpen}
        <button onclick={socket.disconnect}>Disconnect</button>
    {:else}
        <button onclick={socket.connect}>Connect</button>
    {/if}

    {#each socket.auctionResults as ar}
        <div class="card auction-result">
            <br>
            <p>{JSON.stringify(ar, null, 2)}</p>
        </div>
    {/each}
</div>

