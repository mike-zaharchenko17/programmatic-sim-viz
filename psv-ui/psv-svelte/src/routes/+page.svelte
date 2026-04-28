<!--
	RequestID     string       `json:"request_id"`
	Timestamp     time.Time    `json:"timestamp"`
	Winner        *Bid         `json:"winner,omitempty"`
	ClearingPrice float64      `json:"clearing_price"`
	Losers        []LossRecord `json:"losers,omitempty"`
-->

<script lang="ts">
    import { onDestroy } from "svelte"
    import type { AuctionResult } from "$lib/types/types";

    let socket : WebSocket;
    let auctionResults = $state<Array<AuctionResult>>([]);
    let socketIsOpen = $state<Boolean>(false)

    const handleOpenSocket = () => {
        socket = new WebSocket('ws://localhost:1323/ws')
        socket.onmessage = (auctionEvent: any) => {
            const res = JSON.parse(auctionEvent.data)
            auctionResults = [...auctionResults, res]
        }
        socketIsOpen = true
    }

    const handleCloseSocket = () => {
        socket?.close()
        socketIsOpen = false
    }

    onDestroy(() => {
        handleCloseSocket()
    })

</script>

<div class="container">

    {#if socketIsOpen}
        <button onclick={handleCloseSocket}>Disconnect</button>
    {:else}
        <button onclick={handleOpenSocket}>Connect</button>
    {/if}

    {#each auctionResults as ar}
        <div class="card auction-result">
            <br>
            <p>{JSON.stringify(ar, null, 2)}</p>
        </div>
    {/each}
</div>

