import type { AuctionResult, Scope } from "$lib/types/types";

export function inScope(r: AuctionResult, scope: Scope) {
    if (scope.kind === "global") return true

    if (scope.kind === "seat") {
        if (r.winner?.adomain?.[0] === scope.id) return true
        return r.losers?.some((l) => l.bid.adomain?.[0] === scope.id) ?? false
    }

    if (scope.kind === "campaign") {
        if (r.winner?.cid === scope.id) return true
        return r.losers?.some((l) => l.bid.cid === scope.id) ?? false
    }
}

export function wonByScope(r: AuctionResult, scope: Scope): boolean {
    if (!r.winner) return false;
    if (scope.kind === "global") return true;
    if (scope.kind === "seat") return r.winner.adomain?.[0] === scope.id;
    if (scope.kind === "campaign") return r.winner.cid === scope.id;
    return false;
}