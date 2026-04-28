import type { AuctionResult } from "./types/types";

export function createSocket(url: string) {

    let socket : WebSocket | null = $state(null)
    let isOpen = $state(false)
    let auctionResults = $state<AuctionResult[]>([])

    function connect() {
        socket = new WebSocket('ws://localhost:1323/ws')

        socket.onopen = () => {
            isOpen = true
        }

        socket.onmessage = (auctionEvent: any) => {
            const res = JSON.parse(auctionEvent.data)
            auctionResults = [...auctionResults, res]
        }

        socket.onclose = () => {
            isOpen = false
        }

        socket.onerror = (err) => {
            console.error("WebSocket error:", err);
        };
    }

    function disconnect() {
        socket?.close()
        socket = null
        isOpen = false
    }

    return {
        get isOpen() { return isOpen },
        get auctionResults() { return auctionResults },
        connect,
        disconnect

    }
}