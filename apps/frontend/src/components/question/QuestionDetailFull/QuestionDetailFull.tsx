import { Row, TabsProps } from "antd";
import { QuestionDetail } from "../QuestionDetail/QuestionDetail";
import { TestcaseDetail } from "../TestcaseDetail/TestcaseDetail";
import "./styles.scss";
import { Test } from "@/app/services/execute";

// TODO: should add function for test case submission in props
interface QuestionDetailFullProps {
  questionTitle?: string;
  complexity?: string;
  categories?: string[];
  description?: string;
  visibleTestcases?: Test[];
  shouldShowSubmitButton?: boolean;
  handleRunTestCases?: () => void;
  isLoadingTestCase?: boolean;
  buttonIsDisabled?: boolean;
}

export const QuestionDetailFull = (props: QuestionDetailFullProps) => {
  return (
    <>
      <Row className="question-box">
        <QuestionDetail
          questionTitle={props.questionTitle}
          complexity={props.complexity}
          categories={props.categories}
          description={props.description}
        />
      </Row>
      <Row className="test-box">
        <TestcaseDetail
          visibleTestcases={props.visibleTestcases}
          shouldShowSubmitButton={props.shouldShowSubmitButton}
          handleRunTestCases={props.handleRunTestCases}
          isLoadingTestCase={props.isLoadingTestCase}
          buttonIsDisabled={props.buttonIsDisabled}
        />
      </Row>
    </>
  );
};
