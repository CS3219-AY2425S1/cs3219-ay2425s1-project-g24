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
import { handleFindMatch } from '../handlers';

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

const FindMatchContent: React.FC = () => {
    const [selectedDifficulties, setSelectedDifficulties] = useState<string[]>([]);
    const [selectedTopics, setSelectedTopics] = useState<string[]>([]);

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
                onClick={handleFindMatch}
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
                    checked={selectedDifficulties.includes(difficultyOption.value)}
                    onChange={() => handleChange(difficultyOption.value)}
                >
                    {difficultyOption.label}
                </Tag.CheckableTag>
            ))}
        </div>
    )
}

const TopicSelector: React.FC<TopicSelectorProps> = ({ selectedTopics, onChange}) => {
    const topicOptions: SelectProps[] = CategoriesOption;
    
    const handleChange = (topic: string) => {
        const newSelectedTopics = selectedTopics.includes(topic)
            ? selectedTopics.filter(selectedTopic => selectedTopic !== topic)
            : [...selectedTopics, topic];
        onChange(newSelectedTopics);
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
