<!--
	RequestID     string       `json:"request_id"`
	Timestamp     time.Time    `json:"timestamp"`
	Winner        *Bid         `json:"winner,omitempty"`
	ClearingPrice float64      `json:"clearing_price"`
	Losers        []LossRecord `json:"losers,omitempty"`
-->

<script lang="ts">
    import { onMount, onDestroy } from "svelte"
    import type { AuctionResult } from "$lib/types/types";
    
    let socket : WebSocket;
    let auctionResults = $state<Array<AuctionResult>>([]);

    onMount(() => {
        socket = new WebSocket('ws://localhost:1323/ws')
        socket.onmessage = (auctionEvent: any) => {
            const res = JSON.parse(auctionEvent.data)
            auctionResults = [...auctionResults, res]
        }
    })

    onDestroy(() => {
        socket?.close()
    })
</script>

<div class="container">
    {#each auctionResults as ar}
        <div class="auction-result">
            <br>
            <p>{JSON.stringify(ar, null, 2)}</p>
        </div>
    {/each}
</div>

