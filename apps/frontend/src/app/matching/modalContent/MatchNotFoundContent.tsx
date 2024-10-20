import React from 'react';
import 'typeface-montserrat';
import './styles.scss';
import { handleReselectMatchOptions, handleRetryMatch } from '../handlers';
import { formatTime } from '@/utils/DateTime';

const MatchNotFoundContent: React.FC<{
    retry(): void,
    reselect(): void,
    timedOutIn: number,
}> = ({
    retry, reselect, timedOutIn
}) => {
    return (
        <div className="joined-match-content">
            <div className="dissatisfied-icon-container">
                <svg xmlns="http://www.w3.org/2000/svg" 
                    height="24px" 
                    viewBox="0 -960 960 960" 
                    width="24px" 
                    fill="#e8eaed"
                    className="disatisfied-icon"
                >
                    <path d="M620-520q25 0 42.5-17.5T680-580q0-25-17.5-42.5T620-640q-25 0-42.5 17.5T560-580q0 25 17.5 42.5T620-520Zm-280 0q25 0 42.5-17.5T400-580q0-25-17.5-42.5T340-640q-25 0-42.5 17.5T280-580q0 25 17.5 42.5T340-520Zm140 100q-68 0-123.5 38.5T276-280h66q22-37 58.5-58.5T480-360q43 0 79.5 21.5T618-280h66q-25-63-80.5-101.5T480-420Zm0 340q-83 0-156-31.5T197-197q-54-54-85.5-127T80-480q0-83 31.5-156T197-763q54-54 127-85.5T480-880q83 0 156 31.5T763-763q54 54 85.5 127T880-480q0 83-31.5 156T763-197q-54 54-127 85.5T480-80Zm0-400Zm0 320q134 0 227-93t93-227q0-134-93-227t-227-93q-134 0-227 93t-93 227q0 134 93 227t227 93Z"/>
                </svg>
            </div>
            <div className="match-status-label">Match Not Found!</div>
            <div className="match-status-message">
                Sorry, we could not find a match after {formatTime(timedOutIn)}
            </div>
            {/* <button className="retry-match-button"
                onClick={retry}
            >
                Retry
            </button> */}
            <button className="reselect-match-options-button"
                onClick={reselect}
            >
                Reselect Options
            </button>
        </div>
    )
}

export default MatchNotFoundContent;
