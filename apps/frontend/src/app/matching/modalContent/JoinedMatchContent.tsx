import React from 'react';
import {
    Avatar,
 } from 'antd';
import {
    UserOutlined,
} from '@ant-design/icons';
import 'typeface-montserrat';
import './styles.scss';
import { handleCancelMatch } from '../handlers';
import { formatTime } from '@/utils/DateTime';


interface Props {
    cancel(): void
    name1: string, // user's username
    name2: string, // matched user's username
}

const JoinedMatchContent: React.FC<Props> = ({cancel, name1: me, name2: you}) => {
    const matchAlreadyJoined = () => {
        throw new Error('Match already joined.');
    }

    return (
        <div className="joined-match-content">
            <div className="matched-profiles">
                <div className="avatar-caption-container">
                    <Avatar size={64} icon={<UserOutlined />} />
                    <div className="user-caption">{me}</div>
                </div>
                <svg xmlns="http://www.w3.org/2000/svg" 
                    height="24px" 
                    viewBox="0 -960 960 960" 
                    width="24px" 
                    fill="#e8eaed"
                    className="bolt-icon"
                >
                    <path d="m422-232 207-248H469l29-227-185 267h139l-30 208ZM320-80l40-280H160l360-520h80l-40 320h240L400-80h-80Zm151-390Z"/>
                </svg>
                <div className="avatar-caption-container">
                    <Avatar size={64} icon={<UserOutlined />} />
                    <div className="user-caption">{you}</div>
                </div>
            </div>
            <div className="match-status-label">Match Found!</div>
            <div className="match-status-message">
                Waiting for others... {formatTime(83)}
            </div>
            <button className="joined-match-deactivated-button" 
                disabled
            >
                Joined
            </button>
            <button className="cancel-match-button"
                onClick={cancel}
            >
                Cancel
            </button>
        </div>
    )
}

export default JoinedMatchContent;
