import React, { useState } from 'react';
import { LoadingOutlined } from '@ant-design/icons';
import 'typeface-montserrat';
import './styles.scss';
import { handleCancelMatch } from '../handlers';
import { useTimer } from "react-timer-hook"
import {formatTime} from '@/utils/DateTime';

const TIMEOUT = 10;

interface Props {
    cancelMatch(cancelledIn: number): void
    timeout(timeoutIn: number): void
}

const MatchingInProgressContent: React.FC<Props> = ({cancelMatch: cancel, timeout}) => {
    const { totalSeconds } = useTimer({
        expiryTimestamp: new Date(Date.now() + 10 * 1000),
        onExpire: () => timeout(TIMEOUT - totalSeconds),
    });

    return (
        <div className="matching-in-progess-content">
            <LoadingOutlined className="loading-icon"/>
            <div className="match-status-label">Matching in Progress</div>
            <div className="match-status-message">
                Please be patient as we match you with other online users<br/><br/>
                {formatTime(TIMEOUT - totalSeconds)}
            </div>
            <button className="cancel-match-button"
                onClick={() => {
                    cancel(TIMEOUT - totalSeconds);
                }}
            >
                Cancel
            </button>
        </div>
    )
}

export default MatchingInProgressContent;
