import React, { useState } from 'react';
import {
    Tag,
    Select,
    Space,
 } from 'antd';
import {
    DifficultyOption,
    CategoriesOption,
} from '@/utils/SelectOptions';
import type { SelectProps } from 'antd';
import 'typeface-montserrat';
import './styles.scss';
import { ValidateUser } from "@/app/services/user"
import { type MatchRequestParams } from '@/app/services/use-matching';

interface DifficultySelectorProps {
    className?: string;
    selectedDifficulties: string[];
    onChange: (difficulties: string[]) => void;
}

interface TopicSelectorProps {
    className?: string;
    selectedTopics: string[];
    onChange: (topics: string[]) => void;
}

interface Props {
    beginMatch(request: MatchRequestParams): void
}

const FindMatchContent: React.FC<Props> = ({ beginMatch }) => {
    const [selectedDifficulties, setSelectedDifficulties] = useState<string[]>([]);
    const [selectedTopics, setSelectedTopics] = useState<string[]>([]);
    const [isLoading, setIsLoading] = useState(false);

    const handleDifficultyChange = (difficulties: string[]) => {
        setSelectedDifficulties(difficulties);
    };

    const handleTopicChange = (topics: string[]) => {
        setSelectedTopics(topics);
    };

    return (
        <div>
            <div className="find-match-title">Find Match</div>
            <div className="difficulty-label">Difficulty</div>
            <div className="difficulty-selector">
                <DifficultySelector
                    selectedDifficulties={selectedDifficulties}
                    onChange={handleDifficultyChange}
                />
            </div>
            <div className="topic-label">Topic</div>
            <div className="topic-selector">
                <TopicSelector
                    selectedTopics={selectedTopics}
                    onChange={handleTopicChange}
                />          
            </div>
            <button className="find-match-button"
                onClick={async () => {
                    setIsLoading(true);
                    const user = await ValidateUser();
                    beginMatch({
                        email: user.data.email,
                        username: user.data.username,
                        type: "match_request",
                        difficulties: selectedDifficulties,
                        topics: selectedTopics,
                    })
                }}
                disabled={isLoading}
            >
                Find Match
            </button>
        </div>
    )
}

const DifficultySelector: React.FC<DifficultySelectorProps> = ({ selectedDifficulties, onChange}) => {
    const handleChange = (difficulty: string) => {
        const newSelectedDifficulties = selectedDifficulties.includes(difficulty)
            ? selectedDifficulties.filter(selectedDifficulty => selectedDifficulty !== difficulty)
            : [...selectedDifficulties, difficulty];
        onChange(newSelectedDifficulties);
    }

    return (
        <div>
            {DifficultyOption.map(difficultyOption => (
                <Tag.CheckableTag
                    className={`difficulty-tag ${difficultyOption.value}-tag`}
                    key={difficultyOption.value}
                    checked={selectedDifficulties.includes(difficultyOption.label)}
                    onChange={() => handleChange(difficultyOption.label)}
                >
                    {difficultyOption.label}
                </Tag.CheckableTag>
            ))}
        </div>
    )
}

const TopicSelector: React.FC<TopicSelectorProps> = ({ selectedTopics, onChange}) => {
    const topicOptions: SelectProps[] = CategoriesOption;
    
    const handleChange = (topics: string[]) => {
        onChange(topics);
    }

    return (
       <div className="find-match-content">
            <Space className="select-space" direction="vertical">
                <Select
                    className="select-topic"
                    mode="multiple"
                    allowClear
                    style={{ width: '100%' }}
                    placeholder="Select Topic(s)"
                    onChange={handleChange}
                    options={topicOptions}
                />
            </Space>
       </div>
    )
}

export default FindMatchContent;
