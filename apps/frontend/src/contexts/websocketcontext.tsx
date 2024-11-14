import { MatchRequestParams } from "@/app/services/use-matching";
import { createContext } from "react";


export type SocketState = {
    state: "cancelling" | "starting"
} | {
    state: "closed";
    info: MatchInfo | null;
    start(req: MatchRequestParams): void;
} | {
    state: "matching";
    cancel(): void;
    timeout(): void;
};

export type MatchInfo = {
    matchId: string;
    user: string;
    matchedUser: string;
    questionDocRefId: string;
    matchedTopics: string[];
}

export type MatchState = SocketState | {
    state: "found";
    info: MatchInfo;
    ok(): void;
} | {
    state: "timeout";
    ok(): void;
    start(req: MatchRequestParams): void;
};

export const WebSocketContext = createContext<MatchState | null>(null);

