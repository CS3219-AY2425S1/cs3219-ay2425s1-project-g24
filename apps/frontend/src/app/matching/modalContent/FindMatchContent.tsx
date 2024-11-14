import React, { useState } from 'react';
import {
    Tag,
    Select,
    Space,
    Form,
    Button,
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
import { MatchParams } from '../MatchingModal';

interface Props {
    beginMatch(request: MatchRequestParams): void
}

const FindMatchContent: React.FC<Props> = ({ beginMatch }) => {
    return (
        <div>
            <div className="find-match-title">Find Match</div>
            <div className="difficulty-label">Difficulty</div>
            <div className="difficulty-selector">
                <Form.Item<MatchParams> name="difficulties">
                    {/* @ts-ignore Not required to pass value and onChange as since this is handled by Form.Item wrapper*/}
                    <TagInput/>
                </Form.Item>
            </div>
            <div className="topic-label">Topic</div>
            <div className="topic-selector">
                <Form.Item<MatchParams>
                    name="topics"
                    >
                    <Select mode="multiple" allowClear placeholder="Select Topic(s)" options={CategoriesOption}/>
                </Form.Item>
            </div>
            <button className="find-match-button" type="submit">
                FIND MATCH
            </button>
        </div>
    )
}

const TagInput: React.FC<{
    value: string[],
    onChange(value: string[]): void,
}> = ({ value, onChange }) => {
    return <>
        {DifficultyOption.map(difficultyOption => (
            <Tag.CheckableTag
                className={`difficulty-tag ${difficultyOption.value}-tag`}
                key={difficultyOption.value}
                checked={value.includes(difficultyOption.value)}
                onChange={(enabled) => {
                    onChange(
                        enabled 
                            ? [...value, difficultyOption.value] 
                            : value.filter(diff => diff !== difficultyOption.value)
                    )
                }}
            >
                {difficultyOption.label}
            </Tag.CheckableTag>
        ))}
    </>
}

export default FindMatchContent;
