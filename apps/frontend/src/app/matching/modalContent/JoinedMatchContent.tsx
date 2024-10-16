import React, { useState } from 'react';
import { 
    Modal,
    Tag,
    Select,
    Space,
    Avatar,
 } from 'antd';
import {
    LoadingOutlined,
    UserOutlined,
} from '@ant-design/icons';
import {
    DifficultyOption,
    CategoriesOption,
} from '@/utils/SelectOptions';
import type { SelectProps } from 'antd';
import 'typeface-montserrat';
import './styles.scss';
import { handleCancelMatch } from '../handlers';
import { formatTime } from '@/utils/DateTime';

const JoinedMatchContent: React.FC = () => {
    const matchAlreadyJoined = () => {
        throw new Error('Match already joined.');
    }

    return (
        <div className="joined-match-content">
            <div className="matched-profiles">
                <Avatar size={64} icon={<UserOutlined />} />
                <svg xmlns="http://www.w3.org/2000/svg" 
                    height="24px" 
                    viewBox="0 -960 960 960" 
                    width="24px" 
                    fill="#e8eaed"
                    className="bolt-icon"
                >
                    <path d="m422-232 207-248H469l29-227-185 267h139l-30 208ZM320-80l40-280H160l360-520h80l-40 320h240L400-80h-80Zm151-390Z"/>
                </svg>
                <Avatar size={64} icon={<UserOutlined />} />
            </div>
            <div className="match-status-label">Match Found!</div>
            <div className="match-status-message">
                Waiting for others... {formatTime(83)}
            </div>
            <button className="joined-match-deactivated-button" 
                disabled
                onClick={() => matchAlreadyJoined()}
            >
                Joined
            </button>
            <button className="cancel-match-button"
                onClick={() => handleCancelMatch()}
            >
                Cancel
            </button>
        </div>
    )
}

export default JoinedMatchContent;
