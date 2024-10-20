"use client"

import { ReactNode } from "react"
import { WebSocketContext } from "@/contexts/websocketcontext";
import useMatching from "@/app/services/use-matching";
const MATCHING_SERVICE_URL = process.env.NEXT_PUBLIC_MATCHING_SERVICE_URL;

if (MATCHING_SERVICE_URL == undefined) {
    throw "NEXT_PUBLIC_MATCHING_SERVICE_URL was not defined in .env";
}

export default function WebSocketProvider({children}: {children: ReactNode}) {
    const matchState = useMatching();

    return <WebSocketContext.Provider value={matchState}>
        {children}
    </WebSocketContext.Provider>
}