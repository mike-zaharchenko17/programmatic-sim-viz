/* Node Classification Utils */

import type { InputNode } from "../types/types"

export function isCampaign(id: string) {
    return id.startsWith("Campaign:") || id.startsWith("Creative:")
}

export function isOutcome(id: string) {
    return id === "Lost" || id === "Won"
}

export function nodeId(e: string | InputNode) {
    return typeof e === "string" ? e : e.id
}
