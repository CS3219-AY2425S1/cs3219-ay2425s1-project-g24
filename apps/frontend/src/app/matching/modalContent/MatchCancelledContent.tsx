import React from 'react';
import 'typeface-montserrat';
import './styles.scss';
import { handleReselectMatchOptions, handleRetryMatch } from '../handlers';
import { formatTime } from '@/utils/DateTime';
import { Button, Form } from 'antd';

interface Props {
    reselect(): void,
    canceledIn: number,
}

const MatchCancelledContent: React.FC<Props> = ({reselect, canceledIn}) => {
    return (
        <div className="match-cancelled-content">
            <div className="cancel-icon-container">
                <svg xmlns="http://www.w3.org/2000/svg"
                    height="24px" 
                    viewBox="0 -960 960 960" 
                    width="24px"
                    fill="#e8eaed"
                    className="cancel-icon"
                >
                    <path d="m336-280 144-144 144 144 56-56-144-144 144-144-56-56-144 144-144-144-56 56 144 144-144 144 56 56ZM480-80q-83 0-156-31.5T197-197q-54-54-85.5-127T80-480q0-83 31.5-156T197-763q54-54 127-85.5T480-880q83 0 156 31.5T763-763q54 54 85.5 127T880-480q0 83-31.5 156T763-197q-54 54-127 85.5T480-80Zm0-80q134 0 227-93t93-227q0-134-93-227t-227-93q-134 0-227 93t-93 227q0 134 93 227t227 93Zm0-320Z"/>
                </svg>
            </div>
            <div className="match-status-label">Match Cancelled!</div>
            <div className="match-status-message">
                Your match request has been cancelled after waiting {formatTime(canceledIn)}
            </div>
            <Form.Item noStyle>
                <button className="retry-match-button" type="submit">
                    Retry
                </button>
            </Form.Item>
            <button className="reselect-match-options-button"
                onClick={reselect}
            >
                Reselect Options
            </button>
        </div>
    )
}

export default MatchCancelledContent;
