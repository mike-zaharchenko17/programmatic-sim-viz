import type { AuctionResult } from "./types/types";

export type Socket = {
    readonly isOpen: boolean;
    readonly auctionResults: AuctionResult[];
    connect: () => void;
    disconnect: () => void;
    clear: () => void;
}

export function createSocket(url: string): Socket {
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

    function clear() {
        auctionResults = []
    }

    return {
        get isOpen() { return isOpen },
        get auctionResults() { return auctionResults },
        connect,
        disconnect,
        clear
    }
}