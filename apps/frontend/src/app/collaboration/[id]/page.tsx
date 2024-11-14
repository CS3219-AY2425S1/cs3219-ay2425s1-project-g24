"use client";
import Header from "@/components/Header/header";
import {
  Button,
  Col,
  Input,
  Layout,
  Modal,
  message,
  Row,
  TabsProps,
  Tag,
  Typography,
  Spin,
  Tooltip,
} from "antd";
import { Content } from "antd/es/layout/layout";
import "./styles.scss";
import { useRouter } from "next/navigation";
import { useEffect, useRef, useState } from "react";
import { GetSingleQuestion, Question } from "@/app/services/question";
import {
  ClockCircleOutlined,
  CodeOutlined,
  InfoCircleFilled,
  SendOutlined,
  VideoCameraOutlined,
} from "@ant-design/icons";
import CollaborativeEditor, {
  CollaborativeEditorHandle,
} from "@/components/CollaborativeEditor/CollaborativeEditor";
import { WebrtcProvider } from "y-webrtc";
import {
  ExecuteVisibleAndCustomTests,
  ExecuteVisibleAndHiddenTestsAndSubmit,
  ExecutionResults,
  GetVisibleTests,
  isTestResult,
  SubmissionHiddenTestResultsAndStatus,
  SubmissionResults,
  Test,
  TestData,
  TestResult,
} from "@/app/services/execute";
import { QuestionDetailFull } from "@/components/question/QuestionDetailFull/QuestionDetailFull";
import VideoPanel from "@/components/VideoPanel/VideoPanel";

interface CollaborationProps {}

export default function CollaborationPage(props: CollaborationProps) {
  const router = useRouter();
  const providerRef = useRef<WebrtcProvider | null>(null);
  const submissionProviderRef = useRef<WebrtcProvider | null>(null);
  const executionProviderRef = useRef<WebrtcProvider | null>(null);

  const editorRef = useRef<CollaborativeEditorHandle>(null);

  const [isLoading, setIsLoading] = useState<boolean>(false);

  // Code Editor States
  const [historyDocRefId, setHistoryDocRefId] = useState<string | undefined>(
    undefined
  );
  const [code, setCode] = useState<string>("");
  const [questionTitle, setQuestionTitle] = useState<string | undefined>(
    undefined
  );
  const [questionDocRefId, setQuestionDocRefId] = useState<string | undefined>(
    undefined
  );
  const [complexity, setComplexity] = useState<string | undefined>(undefined);
  const [categories, setCategories] = useState<string[]>([]); // Store the selected filter categories
  const [description, setDescription] = useState<string | undefined>(undefined);
  const [selectedLanguage, setSelectedLanguage] = useState("python"); // State to hold the selected language item

  // Session states
  const [collaborationId, setCollaborationId] = useState<string | undefined>(
    undefined
  );
  const [currentUser, setCurrentUser] = useState<string | undefined>(undefined);
  const [matchedUser, setMatchedUser] = useState<string>("Loading...");
  const [sessionDuration, setSessionDuration] = useState<number>(0); // State for count-up timer (TODO: currently using localstorage to store time, change to db stored time in the future)
  const stopwatchRef = useRef<NodeJS.Timeout | null>(null);
  const [matchedTopics, setMatchedTopics] = useState<string[] | undefined>(
    undefined
  );

  useEffect(() => {
    const storedTime = localStorage.getItem("session-duration");
    setSessionDuration(storedTime ? parseInt(storedTime) : 0);
  }, []);

  // Chat states
  const [messageToSend, setMessageToSend] = useState<string | undefined>(
    undefined
  );

  // Test case states
  const [manualTestCase, setManualTestCase] = useState<string | undefined>(
    undefined
  );
  const [visibleTestCases, setVisibleTestCases] = useState<Test[]>([]);
  const [isLoadingTestCase, setIsLoadingTestCase] = useState<boolean>(false);
  const [isLoadingSubmission, setIsLoadingSubmission] =
    useState<boolean>(false);
  const [
    submissionHiddenTestResultsAndStatus,
    setSubmissionHiddenTestResultsAndStatus,
  ] = useState<SubmissionHiddenTestResultsAndStatus | undefined>(undefined);

  // End Button Modal state
  const [isModalOpen, setIsModalOpen] = useState<boolean>(false);
  // Session End Modal State
  const [isSessionEndModalOpen, setIsSessionEndModalOpen] =
    useState<boolean>(false);
  const [countDown, setCountDown] = useState<number>(5);

  // Stops the session duration stopwatch
  const stopStopwatch = () => {
    if (stopwatchRef.current) {
      clearInterval(stopwatchRef.current);
    }
  };

  // Starts the session duration stopwatch
  const startStopwatch = () => {
    if (stopwatchRef.current) {
      clearInterval(stopwatchRef.current);
    }

    stopwatchRef.current = setInterval(() => {
      setSessionDuration((prevTime) => {
        const newTime = prevTime + 1;
        localStorage.setItem("session-duration", newTime.toString());
        return newTime;
      });
    }, 1000);
  };

  // Convert seconds into time of format "hh:mm:ss"
  const formatTime = (seconds: number) => {
    const hours = Math.floor(seconds / 3600);
    const minutes = Math.floor((seconds % 3600) / 60);
    const secs = seconds % 60;

    return (
      (hours > 9 ? hours : "0" + hours) +
      ":" +
      (minutes > 9 ? minutes : "0" + minutes) +
      ":" +
      (secs > 9 ? secs : "0" + secs)
    );
  };

  // Message
  const [messageApi, contextHolder] = message.useMessage();

  const successMessage = (message: string) => {
    messageApi.open({
      type: "success",
      content: message,
    });
  };

  const infoMessage = (message: string) => {
    messageApi.open({
      type: "info",
      content: message,
    });
  };

  const errorMessage = (message: string) => {
    messageApi.open({
      type: "error",
      content: message,
    });
  };

  const sendSubmissionResultsToMatchedUser = (data: SubmissionResults) => {
    if (!providerRef.current) {
      throw new Error("Provider not initialized");
    }
    providerRef.current.awareness.setLocalStateField("submissionResultsState", {
      submissionResults: data,
      id: Date.now(),
    });
  };

  const sendExecutingStateToMatchedUser = (executing: boolean) => {
    if (!providerRef.current) {
      throw new Error("Provider not initialized");
    }
    providerRef.current.awareness.setLocalStateField("executingState", {
      executing: executing,
      id: Date.now(),
    });
  };

  const sendSubmittingStateToMatchedUser = (submitting: boolean) => {
    if (!providerRef.current) {
      throw new Error("Provider not initialized");
    }
    providerRef.current.awareness.setLocalStateField("submittingState", {
      submitting: submitting,
      id: Date.now(),
    });
  };

  const sendExecutionResultsToMatchedUser = (data: ExecutionResults) => {
    if (!providerRef.current) {
      throw new Error("Provider not initialized");
    }
    providerRef.current.awareness.setLocalStateField("executionResultsState", {
      executionResults: data,
      id: Date.now(),
    });
  };

  const updateSubmissionResults = (data: SubmissionResults) => {
    const submissionHiddenTestResultsAndStatus: SubmissionHiddenTestResultsAndStatus = {
      hiddenTestResults: data.hiddenTestResults,
      status: data.status,
    }
    setSubmissionHiddenTestResultsAndStatus(submissionHiddenTestResultsAndStatus);
    localStorage.setItem("submissionHiddenTestResultsAndStatus", JSON.stringify(submissionHiddenTestResultsAndStatus));
    setVisibleTestCases(data.visibleTestResults);
    localStorage.setItem("visibleTestResults", JSON.stringify(data.visibleTestResults));
  };

  const updateExecutionResults = (data: ExecutionResults) => {
    setVisibleTestCases(data.visibleTestResults);
    localStorage.setItem("visibleTestResults", JSON.stringify(data.visibleTestResults));
  };

  const updateLangauge = (data: string) => {
    setSelectedLanguage(data);
  }

  const handleRunTestCases = async () => {
    if (!questionDocRefId) {
      throw new Error("Question ID not found");
    }
    setIsLoadingTestCase(true);
    sendExecutingStateToMatchedUser(true);
    try {
      const data = await ExecuteVisibleAndCustomTests(questionDocRefId, {
        code: code,
        language: selectedLanguage,
        customTestCases: "",
      });
      updateExecutionResults(data);
      successMessage("Test cases executed. Review the results below.");
      sendExecutionResultsToMatchedUser(data);
    } finally {
      setIsLoadingTestCase(false);
      sendExecutingStateToMatchedUser(false);
    }
  };

  const handleSubmitCode = async () => {
    if (!questionDocRefId) {
      throw new Error("Question ID not found");
    }
    setIsLoadingSubmission(true);
    sendSubmittingStateToMatchedUser(true);
    try {
      const data = await ExecuteVisibleAndHiddenTestsAndSubmit(questionDocRefId, {
        code: code,
        language: selectedLanguage,
        user: currentUser ?? "",
        matchedUser: matchedUser ?? "",
        matchedTopics: matchedTopics ?? [],
        title: questionTitle ?? "",
        questionDifficulty: complexity ?? "",
        questionTopics: categories,
      });
      updateExecutionResults({
        visibleTestResults: data.visibleTestResults,
        customTestResults: [],
      });
      updateSubmissionResults(data);
      sendSubmissionResultsToMatchedUser(data);
      successMessage("Code saved successfully!");
    } finally {
      setIsLoadingSubmission(false);
      sendSubmittingStateToMatchedUser(false);
    }
  };

  const handleCodeChange = (code: string) => {
    setCode(code);
  };

  // Fetch the question on initialisation
  useEffect(() => {
    if (!isLoading) {
      setIsLoading(true);
    }

    // Retrieve details from localstorage
    const questionDocRefId: string =
      localStorage.getItem("questionDocRefId") ?? "";
    const collabId: string = localStorage.getItem("collabId") ?? "";
    const matchedUser: string = localStorage.getItem("matchedUser") ?? "";
    const currentUser: string = localStorage.getItem("user") ?? "";
    const matchedTopics: string[] =
      localStorage.getItem("matchedTopics")?.split(",") ?? [];
      const submissionHiddenTestResultsAndStatus: SubmissionHiddenTestResultsAndStatus | undefined = 
      localStorage.getItem("submissionHiddenTestResultsAndStatus")
        ? JSON.parse(localStorage.getItem("submissionHiddenTestResultsAndStatus") as string)
        : undefined;
    const visibleTestCases: Test[] = JSON.parse(localStorage.getItem("visibleTestResults") ?? "[]") ?? [];

    // Set states from localstorage
    setCollaborationId(collabId);
    setMatchedUser(matchedUser);
    setCurrentUser(currentUser);
    setMatchedTopics(matchedTopics);
    setQuestionDocRefId(questionDocRefId);
    setSubmissionHiddenTestResultsAndStatus(submissionHiddenTestResultsAndStatus);
    setVisibleTestCases(visibleTestCases);

    GetSingleQuestion(questionDocRefId).then((data: Question) => {
      setQuestionTitle(`${data.id}. ${data.title}`);
      setComplexity(data.complexity);
      setCategories(data.categories);
      setDescription(data.description);
    });

    if (visibleTestCases.length == 0) {
      GetVisibleTests(questionDocRefId).then((data: Test[]) => {
        setVisibleTestCases(data);
      }).catch((e) => {
        errorMessage(e.message);
      });
    }

    // Start stopwatch
    startStopwatch();
  }, []);

  // useEffect for timer
  useEffect(() => {
    if (isSessionEndModalOpen && countDown > 0) {
      const timer = setInterval(() => {
        setCountDown((prevCountDown) => prevCountDown - 1);
      }, 1000);

      return () => clearInterval(timer); // Clean up on component unmount or when countdown changes
    } else if (countDown === 0) {
      router.push("/matching"); // Redirect to matching page
    }
  }, [isSessionEndModalOpen, countDown]);

  // Tabs component items for visibleTestCases
  var items: TabsProps["items"] = visibleTestCases.map((item, index) => {
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
  });

  // Handles the cleaning of localstorage variables, stopping the timer & signalling collab user on webrtc
  // type: "initiator" | "peer"
  const handleCloseCollaboration = (type: string) => {
    // Stop stopwatch
    stopStopwatch();
    if (editorRef.current && type === "initiator") {
      editorRef.current.endSession(); // Call the method on the editor
    }

    // Trigger modal open showing session end details
    setIsSessionEndModalOpen(true);

    // Remove localstorage variables for collaboration
    localStorage.removeItem("session-duration"); // TODO: Remove this after collaboration backend data stored
    localStorage.removeItem("user");
    localStorage.removeItem("matchedUser");
    localStorage.removeItem("collabId");
    localStorage.removeItem("questionDocRefId");
    localStorage.removeItem("matchedTopics");
    localStorage.removeItem("submissionHiddenTestResultsAndStatus");
    localStorage.removeItem("visibleTestResults");
    localStorage.removeItem("editor-language"); // Remove editor language type when session closed
  };

  return (
    <Layout className="collaboration-layout">
      {contextHolder}
      <Header selectedKey={undefined} />
      <Content className="collaboration-content">
        <Modal
          height={500}
          title={"Session Ended"}
          footer={null}
          open={isSessionEndModalOpen}
          width={400}
          closable={false}
        >
          <p className="session-modal-description">
            The collaboration session has ended. You will be redirected in{" "}
            {countDown} seconds
          </p>
          <p className="session-modal-question">
            Question:{" "}
            <span className="session-modal-title">{questionTitle}</span>
          </p>
          <p className="session-modal-difficulty">
            Difficulty:{" "}
            <Tag
              className="complexity-tag"
              style={{
                color:
                  complexity === "easy"
                    ? "#2DB55D"
                    : complexity === "medium"
                    ? "orange"
                    : "red",
              }}
            >
              {complexity &&
                complexity.charAt(0).toUpperCase() + complexity.slice(1)}
            </Tag>
          </p>
          <p className="session-modal-duration">
            Duration:{" "}
            <span className="session-modal-time">
              {formatTime(sessionDuration)}
            </span>
          </p>
          <p className="session-modal-matched-user">
            Matched User:{" "}
            <span className="session-modal-matched-user-name">
              {matchedUser}
            </span>
          </p>
        </Modal>
        <Row gutter={0} className="collab-row">
          <Col span={7} className="first-col">
            <QuestionDetailFull
              questionTitle={questionTitle}
              complexity={complexity}
              categories={categories}
              description={description}
              visibleTestcases={visibleTestCases}
              shouldShowSubmitButton
              handleRunTestCases={handleRunTestCases}
              isLoadingTestCase={isLoadingTestCase}
              buttonIsDisabled={isLoadingTestCase || isLoadingSubmission}
            />
          </Col>
          <Col span={11} className="second-col">
            <Row className="code-row">
              <div className="code-container">
                <div className="code-top-container">
                  <div className="code-title">
                    <CodeOutlined className="title-icons" />
                    Code
                  </div>
                  {/* TODO: Link to execution service for code submission */}
                  <div className="test-button-container">
                    <div className="spinner-container">
                      {isLoadingSubmission && <Spin tip="Saving code..." />}
                    </div>
                    <Button
                      icon={<SendOutlined />}
                      iconPosition="end"
                      onClick={handleSubmitCode}
                      disabled={isLoadingSubmission || isLoadingTestCase}
                    >
                      Submit
                    </Button>
                  </div>
                </div>
                {collaborationId && currentUser && selectedLanguage && (
                  <CollaborativeEditor
                    ref={editorRef}
                    user={currentUser}
                    collaborationId={collaborationId}
                    updateLanguage={updateLangauge}
                    setMatchedUser={setMatchedUser}
                    handleCloseCollaboration={handleCloseCollaboration}
                    providerRef={providerRef}
                    updateSubmissionResults={updateSubmissionResults}
                    updateExecutionResults={updateExecutionResults}
                    matchedUser={matchedUser}
                    onCodeChange={handleCodeChange}
                    updateExecuting={setIsLoadingTestCase}
                    updateSubmitting={setIsLoadingSubmission}
                  />
                )}
                <div className="hidden-test-results">
                  <Tooltip title="Status applies only to this session">
                    <InfoCircleFilled className="hidden-test-icon" />
                  </Tooltip>
                  <Typography.Text
                    strong
                    style={{
                      color: submissionHiddenTestResultsAndStatus
                        ? submissionHiddenTestResultsAndStatus.status ===
                          "Accepted"
                          ? "green"
                          : submissionHiddenTestResultsAndStatus.status ===
                            "Attempted"
                          ? "orange"
                          : "black" // default color for any other status
                        : "gray", // color for "Not Attempted"
                    }}
                  >
                    Session Status:{" "}
                    {submissionHiddenTestResultsAndStatus
                      ? submissionHiddenTestResultsAndStatus.status
                      : "Not Attempted"}
                  </Typography.Text>
                  <br />
                  {submissionHiddenTestResultsAndStatus && (
                    <Typography.Text
                      strong
                      style={{
                        color:
                          submissionHiddenTestResultsAndStatus.hiddenTestResults
                            .passed ===
                          submissionHiddenTestResultsAndStatus.hiddenTestResults
                            .total
                            ? "green" // All test cases passed
                            : "red", // Some test cases failed
                      }}
                    >
                      Passed{" "}
                      {
                        submissionHiddenTestResultsAndStatus.hiddenTestResults
                          .passed
                      }{" "}
                      /{" "}
                      {
                        submissionHiddenTestResultsAndStatus.hiddenTestResults
                          .total
                      }{" "}
                      hidden test cases
                    </Typography.Text>
                  )}
                </div>
              </div>
            </Row>
          </Col>
          <Col span={6} className="third-col">
            <Row className="session-row">
              <div className="session-container">
                <div className="session-top-container">
                  <div className="session-title">
                    <ClockCircleOutlined className="title-icons" />
                    Session Details
                  </div>
                  <Modal
                    height={500}
                    title={"End Session"}
                    okText={"End"}
                    okButtonProps={{ danger: true }}
                    onOk={() => handleCloseCollaboration("initiator")}
                    open={isModalOpen}
                    onCancel={() => setIsModalOpen(false)}
                    width={400}
                  >
                    <p className="modal-description">
                      Are you sure you want to quit the existing collaboration
                      session? This will end the session for both users!
                    </p>
                  </Modal>
                  <Button
                    danger
                    onClick={() => setIsModalOpen(true)}
                    className="session-end-button"
                  >
                    End for All
                  </Button>
                </div>

                <div className="session-duration">
                  Duration:
                  <span className="session-duration-timer">
                    {formatTime(sessionDuration)}
                  </span>
                </div>
                <div className="session-matched-user-label">
                  Matched User:
                  <span className="session-matched-user-name">
                    {matchedUser}
                  </span>
                </div>
              </div>
            </Row>
            <Row className="chat-row">
              <div className="chat-container">
                <div className="chat-title">
                  <VideoCameraOutlined className="title-icons" />
                  Video
                </div>
                <VideoPanel />
                {/* <div className="chat-message-box">
                  <div className="chat-header-message">
                    Matched with {matchedUser}
                  </div>
                  <div></div>
                </div>
                <div className="chat-typing-box">
                  <Input.TextArea
                    onChange={(e) => setMessageToSend(e.target.value)}
                    placeholder="Send Message Here"
                    rows={4}
                  />
                </div> */}
              </div>
            </Row>
          </Col>
        </Row>
      </Content>
    </Layout>
  );
}
