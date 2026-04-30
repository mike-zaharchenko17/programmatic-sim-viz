/* uses the primitives defined in data-processing.ts to compute sankey vals */

import type { InputLink, InputNode, NestedBuckets } from "./types/types";
import type { Bid, AuctionResult } from "./types/types"

export function nestedToLinks(
    nested: NestedBuckets
) : InputLink[] {

    const links: InputLink[] = [];

    for (const [source, targets] of Object.entries(nested)) {
        for (const [target, items] of Object.entries(targets)) {
            links.push({
                source,
                target,
                value: items.length
            })
        }
    }

    // deterministic sort for stable sankey
    links.sort((a, b) => {
        const sourceCmp = a.source.localeCompare(b.source)
        if (sourceCmp !== 0) {
            return sourceCmp
        }
        return a.target.localeCompare(b.target)
    })

    return links
}

export function linksToNodes(
    links: InputLink[]
) : InputNode[] {
    const ids = new Set<string>()

    for (const link of links) {
        ids.add(link.source)
        ids.add(link.target)
    }

    return Array.from(ids)
        .sort((a, b) => a.localeCompare(b))
        .map(id => ({ id }))
}

export function getADomain(bid: Bid) : string {
    return bid.adomain?.[0] ?? "Unknown domain"
}

export function getBidder(bid: Bid) : string {
    // `bid.id` is usually unique per bid event; prefer stable identifiers
    // so Sankey middle nodes aggregate instead of exploding in cardinality.
    if (bid.cid) {
        return `Campaign:${bid.cid}`
    }
    if (bid.crid) {
        return `Creative:${bid.crid}`
    }
    return "Unknown bidder"
}

export function getOutcome(lossReason?: number): string {
    if (lossReason === undefined) {
        return "Won"
    }

    return "Lost: Outbid"
}

export function auctionResultsToSankeyLinks(
    results: AuctionResult[]
  ): InputLink[] {
    const counts = new Map<string, number>();
  
    function increment(source: string, target: string) {
      const key = `${source}|||${target}`;
      counts.set(key, (counts.get(key) ?? 0) + 1);
    }
  
    for (const result of results) {
      if (result.winner) {
        const domain = getADomain(result.winner);
        const bidder = getBidder(result.winner);
  
        increment(domain, bidder);
        increment(bidder, getOutcome());
      }
  
      for (const loser of result.losers ?? []) {
        const domain = getADomain(loser.bid);
        const bidder = getBidder(loser.bid);
        const outcome = getOutcome(loser.loss_reason);
  
        increment(domain, bidder);
        increment(bidder, outcome);
      }
    }
  
    return Array.from(counts.entries())
      .map(([key, value]) => {
        const [source, target] = key.split("|||");
  
        return {
          source,
          target,
          value
        };
      })
      .sort((a, b) => {
        const sourceCmp = a.source.localeCompare(b.source);
        if (sourceCmp !== 0) return sourceCmp;
        return a.target.localeCompare(b.target);
      });
}