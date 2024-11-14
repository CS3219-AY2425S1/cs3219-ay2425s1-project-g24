import React, { useState } from 'react';
import { LoadingOutlined } from '@ant-design/icons';
import 'typeface-montserrat';
import './styles.scss';
import { handleCancelMatch } from '../handlers';
import { useTimer } from "react-timer-hook"
import {formatTime} from '@/utils/DateTime';

interface Props {
    cancelMatch(): void
    timePassed: number
}

const MatchingInProgressContent: React.FC<Props> = ({cancelMatch: cancel, timePassed}) => {

    return (
        <div className="matching-in-progess-content">
            <LoadingOutlined className="loading-icon"/>
            <div className="match-status-label">Matching in Progress</div>
            <div className="match-status-message">
                Please be patient as we match you with other online users<br/><br/>
                {formatTime(timePassed)}
            </div>
            <button className="cancel-match-button"
                onClick={() => {
                    cancel();
                }}
            >
                Cancel
            </button>
        </div>
    )
}

export default MatchingInProgressContent;
