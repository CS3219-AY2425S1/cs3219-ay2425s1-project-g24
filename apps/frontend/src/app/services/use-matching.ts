import { MatchInfo, MatchState } from "@/contexts/websocketcontext";
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
    type: "match_question_found",
    match_id: string,
    user: string,
    matched_user: string,
    matched_topics: string[],
    question_doc_ref_id: string,
    question_name: string,
    question_difficulty: string,
    question_topics: string[],
}

export type MatchTimeoutResponse = {
    type: "timeout",
    message: string,
}

export type MatchRejectedResponse = {
    type: "match_rejected",
    message: string,
}

type MatchResponse = MatchFoundResponse | MatchTimeoutResponse | MatchRejectedResponse;

export default function useMatching(): MatchState {
    const [isSocket, setIsSocket] = useState<boolean>(false);
    const [ste, setSte] = useState<MatchState>({
        state: "closed",
        info: null,
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

            if (responseJson.type == "match_question_found") {
                setIsSocket(false);

                const info: MatchInfo = parseInfoFromResponse(responseJson);
                setSte({
                    state: "found",
                    info: info,
                    ok: () => {
                        setIsSocket(false);
                        setSte({
                            state: "closed",
                            info: info,
                            start,
                        });
                    }
                })
                return;
            }

            if (responseJson.type == "match_rejected") {
                console.log("match rejected: " + responseJson.message);
                cancel();
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
            start
        });
    }

    function cancel() {
        setIsSocket(false)
        setSte({
            state: "closed",
            info: null,
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
            matchState = {state: "closed", info: null, start}
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

function parseInfoFromResponse(responseJson: MatchFoundResponse): MatchInfo {
    return {
        matchId: responseJson.match_id?.toString() ?? "unknown",
        matchedUser: responseJson.matched_user ?? "unknown",
        user: responseJson.user ?? "unknown",
        questionDocRefId: responseJson.question_doc_ref_id ?? "unknown",
        matchedTopics: responseJson.matched_topics ?? [],
    };
}
