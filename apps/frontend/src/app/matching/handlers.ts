export const handleFindMatch = () => {
    // start timer
    // open socket connection
    // change matching state to 'matching'
}

export const handleCancelMatch = () => {
    // close socket connection
    // stop timer
    // change matching state to 'cancelled'
}

// TODO: wait for collaboration service
export const handleJoinMatch = () => {
    // notify server that user has joined the match (close connection?)
    // change matching state to 'joined'
}

export const handleRetryMatch = () => {
    // retain current difficulty and topics
    // same as handleFindMatch
}

export const handleReselectMatchOptions = () => {
    // reset selected difficulties and topics
    // same as handleFindMatch
}
