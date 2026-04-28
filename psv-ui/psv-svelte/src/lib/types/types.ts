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