import { Tween } from "svelte/motion"

export type Bid = {
    id: string,
    impid: string,
    price: number,
    adomain?: string[],
    cid?: string,
    crid?: string,
    w?: number,
    h?: number,
}

export type AuctionResult = {
    request_id: string,
    timestamp: string,
    winner?: Bid,
    clearing_price: number,
    losers?: {
        bid: Bid,
        loss_reason: number
    }[]
}
// coordinate point tween object
export type NodeTween = {
    x0: Tween<number>;
    y0: Tween<number>;
    x1: Tween<number>;
    y1: Tween<number>;
}

export type LinkTween = {
    y0: Tween<number>;
    y1: Tween<number>;
    width: Tween<number>;
}
export type InputNode = {
    id: string;
}

export type InputLink = {
    source: string;
    target: string;
    value: number;
}

export type Scope =
    | { kind: "global"}
    | { kind: "seat"; id: string }
    | { kind: "campaign"; id: string }

export type NestedBuckets = Record<string, Record<string, unknown[]>>