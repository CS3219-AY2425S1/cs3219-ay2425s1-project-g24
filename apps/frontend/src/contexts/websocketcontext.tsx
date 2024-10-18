import { MatchRequestParams } from "@/app/services/use-matching";
import { createContext } from "react";


export type SocketState = {
    state: "cancelling" | "starting"
} | {
    state: "closed";
    start(req: MatchRequestParams): void;
} | {
    state: "matching";
    cancel(): void;
    timeout(): void;
};

export type MatchState = SocketState | {
    state: "found";
    info: {
        matchId: string;
        partnerId: string;
        partnerName: string;
    };
    ok(): void;
} | {
    state: "timeout";
    ok(): void;
};

export const WebSocketContext = createContext<MatchState | null>(null);

