import {
  FileDoneOutlined,
  InfoCircleFilled,
  PlayCircleOutlined,
} from "@ant-design/icons";
import { Button, Input, Spin, Tabs, TabsProps, Tooltip, Typography } from "antd";
import "./styles.scss";
import { isTestResult, Test } from "@/app/services/execute";

interface TestcaseDetailProps {
  visibleTestcases?: Test[];
  shouldShowSubmitButton?: boolean;
  isLoadingTestCase?: boolean;
  handleRunTestCases?: () => void;
  buttonIsDisabled?: boolean;
}

export const TestcaseDetail = (props: TestcaseDetailProps) => {
  const defaultTestcaseTabs: TabsProps["items"] = [
    {
      key: "1",
      label: "Case 1",
      children: (
        <Input.TextArea disabled placeholder="Insert Test Case 1" rows={6} />
      ),
    },
    {
      key: "2",
      label: "Case 2",
      children: (
        <Input.TextArea disabled placeholder="Insert Test Case 2" rows={6} />
      ),
    },
    // {
    //   key: "3",
    //   label: "Case 3",
    //   children: (
    //     <Input.TextArea placeholder="Input Manual Test Case" rows={6} />
    //   ),
    // },
  ];

  var visibleTestcasesTabs: TabsProps["items"] = props.visibleTestcases?.map(
    (item, index) => {
      return {
        key: index.toString(),
        label: (
          <span
            style={{
              color: !isTestResult(item) ? "" : item.passed ? "green" : "red",
            }}
          >
            Case {index + 1}
          </span>
        ),
        children: (
          <div>
            <Input.TextArea
              disabled
              placeholder={`Stdin: ${item.input}\nStdout: ${item.expected}`}
              rows={4}
            />
            {isTestResult(item) && (
              <div className="test-result-container">
                <InfoCircleFilled className="hidden-test-icon" />
                <Typography.Text
                  strong
                  style={{ color: item.passed ? "green" : "red" }}
                >
                  {item.passed ? "Passed" : "Failed"}
                </Typography.Text>
                <br />
                <Typography.Text strong>Actual Output:</Typography.Text>{" "}
                {item.actual}
                <br />
                {item.error && (
                  <>
                    <Typography.Text strong>Error:</Typography.Text>
                    <div className="error-message">{item.error}</div>
                  </>
                )}
              </div>
            )}
          </div>
        ),
      };
    }
  );

  return (
    <div className="test-section">
      <div className="test-top-box">
        <div className="test-title">
          <FileDoneOutlined className="title-icons" />
          Test Cases
        </div>
        {/* TODO: Link to execution service for running code against test-cases */}
        {props.shouldShowSubmitButton && (
          <div className="test-button-container">
            <div className="spinner-container">
              {props.isLoadingTestCase && <Spin tip="Running test cases..." />}
            </div>
            <Tooltip title="Run test cases will run public test cases">
              <InfoCircleFilled className="run-test-cases-info" />
            </Tooltip>
            <Button
              icon={<PlayCircleOutlined />}
              iconPosition="end"
              className="test-case-button"
              onClick={props.handleRunTestCases}
              disabled={props.buttonIsDisabled}
            >
              Run Test Cases
            </Button>
          </div>
        )}
      </div>
      <div className="test-cases-container">
        <Tabs
          items={visibleTestcasesTabs || defaultTestcaseTabs}
          defaultActiveKey="1"
        />
      </div>
    </div>
  );
};
