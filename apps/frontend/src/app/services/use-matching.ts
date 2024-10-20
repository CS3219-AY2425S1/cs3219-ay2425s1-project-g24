import { MatchState } from "@/contexts/websocketcontext";
import { useEffect, useState } from "react";
import useWebSocket, { Options, ReadyState } from "react-use-websocket";

const MATCHING_SERVICE_URL = process.env.NEXT_PUBLIC_MATCHING_SERVICE_URL;

if (MATCHING_SERVICE_URL == undefined) {
    throw "NEXT_PUBLIC_MATCHING_SERVICE_URL was not defined in .env";
}

export type MatchRequestParams = {
    type: "match_request",
    username: string,
    email: string,
    topics: string[],
    difficulties: string[],
}

export type MatchFoundResponse = {
    type: "match_found",
    matchId: number,
    partnerId: number,
    partnerName: string,
}

export type MatchTimeoutResponse = {
    type: "timeout",
    message: string,
}

export type MatchRejectedResponse = {
    type: "match_rejected",
    message: string,
}

type MatchResponse = MatchFoundResponse | MatchTimeoutResponse;

export default function useMatching(): MatchState {
    const [isSocket, setIsSocket] = useState<boolean>(false);
    const [ste, setSte] = useState<MatchState>({
        state: "closed",
        start,
    });

    const options: Options = {
        onClose() {
            setIsSocket(false);
        },
        onMessage({data: response}) {
            const responseJson: MatchResponse = JSON.parse(response);
            if (responseJson.type == "timeout") {
                timeout();
                return;
            }

            if (responseJson.type == "match_found") {
                setIsSocket(false);
                setSte({
                    state: "found",
                    info: {
                        matchId: responseJson.matchId.toString(),
                        partnerId: responseJson.partnerId.toString(),
                        partnerName: responseJson.partnerName,
                    },
                    ok: cancel
                })
                return;
            }
        }
    }

    const { 
        readyState: socketState, 
        sendJsonMessage, 
    } = useWebSocket<MatchResponse>(MATCHING_SERVICE_URL as string, options, isSocket);
    
    function timeout() {
        setIsSocket(false);
        setSte({
            state: "timeout",
            ok: cancel,
        });
    }

    function cancel() {
        setIsSocket(false)
        setSte({
            state: "closed",
            start,
        })
    }
    
    function start(request: MatchRequestParams) {
        setIsSocket(true)
        sendJsonMessage(request);
    }
    
    let matchState: MatchState;
    switch (socketState) {
        case ReadyState.CLOSED:
        case ReadyState.UNINSTANTIATED:
            matchState = {state: "closed", start}
            break;
        case ReadyState.OPEN:
            matchState = {state: "matching", cancel, timeout}
            break;
        case ReadyState.CONNECTING:
            matchState = {state: "starting"}
            break;
        case ReadyState.CLOSING:
            matchState = {state: "cancelling"}
            break;
    }

    return isSocket ? matchState : ste;
}