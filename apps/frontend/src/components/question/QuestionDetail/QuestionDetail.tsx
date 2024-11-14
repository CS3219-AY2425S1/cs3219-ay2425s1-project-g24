import { Tag } from "antd";
import "./styles.scss";

interface QuestionDetailProps {
  questionTitle?: string;
  complexity?: string;
  categories?: string[];
  description?: string;
}

// Returns the html of question details (without testcase information)
export const QuestionDetail = (props: QuestionDetailProps) => {
  return (
    <div className="question-container">
      <div className="question-title">{props.questionTitle}</div>
      <div className="question-difficulty">
        <Tag
          className="complexity-tag"
          style={{
            color:
              props.complexity === "easy"
                ? "#2DB55D"
                : props.complexity === "medium"
                ? "orange"
                : "red",
          }}
        >
          {props.complexity &&
            props.complexity.charAt(0).toUpperCase() +
              props.complexity.slice(1)}
        </Tag>
      </div>
      <div className="question-topic">
        <span className="topic-label">Topics: </span>
        {props.categories?.map((category) => (
          <Tag key={category}>{category}</Tag>
        ))}
      </div>
      <div className="question-description">{props.description}</div>
    </div>
  );
};
