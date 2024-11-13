"use client";
import Header from "@/components/Header/header";
import { Col, Layout, message, PaginationProps, Row, Spin, Table } from "antd";
import { Content } from "antd/es/layout/layout";
import { CodeOutlined, HistoryOutlined } from "@ant-design/icons";
import "./styles.scss";
import { useEffect, useRef, useState } from "react";
import { GetSingleQuestion } from "../../services/question";
import React from "react";
import { useSearchParams } from "next/navigation";
import { useRouter } from "next/navigation";
import { QuestionDetailFull } from "@/components/question/QuestionDetailFull/QuestionDetailFull";
import { Compartment, EditorState } from "@codemirror/state";
import { basicSetup, EditorView } from "codemirror";
import { javascript } from "@codemirror/lang-javascript";
import { GetHistory, GetUserQuestionHistories } from "@/app/services/history";
import { ValidateUser, VerifyTokenResponseType } from "@/app/services/user";
import { GetVisibleTests, Test } from "@/app/services/execute";

interface Submission {
  submittedAt: string;
  language: string;
  matchedUser: string;
  otherUser: string;
  historyDocRefId: string;
  code: string;
}

interface TablePagination {
  totalCount: number;
  totalPages: number;
  currentPage: number;
  limit: number;
}

export default function QuestionPage() {
  const [isLoading, setIsLoading] = useState<boolean>(true); // Store the states related to table's loading

  // Message States
  const [messageApi, contextHolder] = message.useMessage();

  const errorMessage = (message: string) => {
    messageApi.open({
      type: "error",
      content: message,
    });
  };

  const editorRef = useRef(null);
  const languageConf = new Compartment();

  // Retrieve the questionDocRefId and historyDocRefId from query params during page navigation
  const searchParams = useSearchParams();
  const questionDocRefId: string = searchParams?.get("data") ?? "";
  const historyDocRefId: string = searchParams?.get("history") ?? "";
  // Code Editor States
  const [questionTitle, setQuestionTitle] = useState<string | undefined>(
    undefined
  );
  const [complexity, setComplexity] = useState<string | undefined>(undefined);
  const [categories, setCategories] = useState<string[]>([]); // Store the selected filter categories
  const [description, setDescription] = useState<string | undefined>(undefined);
  const [username, setUsername] = useState<string | undefined>(undefined);
  const [userQuestionHistories, setUserQuestionHistories] =
    useState<History[]>();
  const [submission, setSubmission] = useState<Submission>();
  const [isHistoryLoading, setIsHistoryLoading] = useState<boolean>(true);
  const [isSubmissionLoading, setIsSubmissionLoading] = useState<boolean>(true);
  const [currentSubmissionId, setCurrentSubmissionId] = useState<
    string | undefined
  >(historyDocRefId == "" ? undefined : historyDocRefId);
  const [paginationParams, setPaginationParams] = useState<TablePagination>({
    totalCount: 0,
    totalPages: 0,
    currentPage: 1,
    limit: 3,
  });
  const [visibleTestCases, setVisibleTestCases] = useState<Test[]>([]);

  const state = EditorState.create({
    doc: "",
    extensions: [
      basicSetup,
      languageConf.of(javascript()),
      EditorView.theme({
        "&": { height: "100%", overflow: "hidden" }, // Enable Scroll
      }),
      EditorView.editable.of(false), // Disable editing
    ],
  });

  // Handler for change in page jumper
  const onPageJump: PaginationProps["onChange"] = (pageNumber) => {
    setPaginationParams((prev) => {
      loadQuestionHistories(pageNumber, prev.limit);
      return { ...paginationParams, currentPage: pageNumber };
    });
  };

  async function loadQuestionHistories(currentPage: number, limit: number) {
    if (username === undefined) return;
    setIsHistoryLoading(true);
    GetUserQuestionHistories(username, questionDocRefId, currentPage, limit)
      .then((data: any) => {
        setUserQuestionHistories(data.histories);
        setPaginationParams({
          ...paginationParams,
          totalCount: data.totalCount,
          totalPages: data.totalPages,
          currentPage: data.currentPage,
          limit: data.limit,
        });
      })
      .finally(() => {
        setIsHistoryLoading(false);
      });
  }

  // When code editor page is initialised, fetch the particular question, and display in code editor
  useEffect(() => {
    if (!isLoading) {
      setIsLoading(true);
    }

    GetSingleQuestion(questionDocRefId)
      .then((data: any) => {
        setQuestionTitle(data.title);
        setComplexity(data.complexity);
        setCategories(data.categories);
        setDescription(data.description);
      })
      .catch((e) => {
        errorMessage(e.message);
      })
      .finally(() => {
        setIsLoading(false);
      });

    GetVisibleTests(questionDocRefId)
      .then((data: Test[]) => {
        setVisibleTestCases(data);
      })
      .catch((e) => {
        errorMessage(e.message);
      });
  }, [questionDocRefId]);

  useEffect(() => {
    loadQuestionHistories(paginationParams.currentPage, paginationParams.limit);
  }, [username]);

  useEffect(() => {
    ValidateUser().then((data: VerifyTokenResponseType) => {
      setUsername(data.data.username);
    });
  }, []);

  useEffect(() => {
    // Only show history if a history is selected
    if (currentSubmissionId === undefined) return;

    const view = new EditorView({
      state,
      parent: editorRef.current || undefined,
    });

    setIsSubmissionLoading(true);
    GetHistory(currentSubmissionId)
      .then((data: any) => {
        const submittedAt = new Date(data.createdAt);
        setSubmission({
          submittedAt: submittedAt.toLocaleString("en-US"),
          language: data.language,
          matchedUser:
            username == data.matchedUser ? data.user : data.matchedUser,
          otherUser: data.user,
          historyDocRefId: data.historyDocRefId,
          code: data.code,
        });
        setIsSubmissionLoading(false);

        view.dispatch(
          state.update({
            changes: { from: 0, to: state.doc.length, insert: data.code },
          })
        );
      })
      .finally(() => {});

    return () => {
      // Cleanup on component unmount
      view.destroy();
    };
  }, [currentSubmissionId]);

  const columns = [
    {
      title: "Submitted at",
      dataIndex: "createdAt",
      key: "createdAt",
      render: (date: string) => {
        return new Date(date).toLocaleString();
      },
    },
    {
      title: "Language",
      dataIndex: "language",
      key: "language",
    },
    {
      title: "Matched with",
      dataIndex: "matchedUser",
      key: "matchedUser",
    },
  ];

  const handleRowClick = (s: Submission) => {
    setCurrentSubmissionId(s.historyDocRefId);
  };

  return (
    <div>
      {contextHolder}
      <Layout className="question-layout">
        <Header selectedKey={undefined} />
        <Content className="question-content">
          <Row gutter={0} className="question-row">
            <Col span={12} className="first-col">
              <QuestionDetailFull
                questionTitle={questionTitle}
                complexity={complexity}
                categories={categories}
                description={description}
                visibleTestcases={visibleTestCases}
              />
            </Col>
            <Col span={12} className="second-col">
              <Row className="history-row">
                <div className="history-container">
                  <div className="history-top-container">
                    <div className="history-title">
                      <HistoryOutlined className="title-icons" />
                      Submission History
                    </div>
                    <div style={{ margin: "10px" }}>
                      <Table
                        rowKey="id"
                        dataSource={userQuestionHistories}
                        columns={columns}
                        onRow={(record: any) => {
                          return {
                            onClick: () => handleRowClick(record),
                            style: { cursor: "pointer" },
                          };
                        }}
                        loading={isHistoryLoading}
                        pagination={{
                          size: "small",
                          current: paginationParams.currentPage,
                          total: paginationParams.totalCount,
                          pageSize: paginationParams.limit,
                          onChange: onPageJump,
                        }}
                        scroll={{ y: 200 }}
                      />
                    </div>
                  </div>
                </div>
              </Row>
              {currentSubmissionId && (
                <Row className="code-row">
                  <div className="code-container">
                    {isSubmissionLoading && (
                      <div
                        style={{
                          display: "flex",
                          justifyContent: "center",
                          alignItems: "center",
                          height: "100%",
                        }}
                      >
                        <Spin />
                      </div>
                    )}
                    <div
                      style={{
                        visibility: `${
                          isSubmissionLoading ? "hidden" : "visible"
                        }`,

                        display: "flex",
                        flexDirection: "column",
                        height: "100%",
                      }}
                    >
                      <div className="code-top-container">
                        <div className="code-title">
                          <CodeOutlined className="title-icons" />
                          Submitted Code
                        </div>
                      </div>

                      {/* Details of submission */}
                      <div
                        style={{
                          margin: "10px",
                          display: "flex",
                          justifyContent: "space-between",
                          flexDirection: "row",
                        }}
                      >
                        <div className="submission-header-detail">
                          <div style={{ fontWeight: "bold" }}>
                            Submitted at:
                          </div>
                          <div style={{ paddingLeft: "5px" }}>
                            {submission?.submittedAt || "-"}
                          </div>
                        </div>
                        <div className="submission-header-detail">
                          <div style={{ fontWeight: "bold" }}>Language:</div>
                          <div style={{ paddingLeft: "5px" }}>
                            {submission?.language || "-"}
                          </div>
                        </div>
                        <div className="submission-header-detail">
                          <div style={{ fontWeight: "bold" }}>
                            Matched with:
                          </div>
                          <div style={{ paddingLeft: "5px" }}>
                            {submission?.matchedUser
                              ? // Check to ensure that matched user is correct, otherwise swap with otherUser
                                username == submission.matchedUser
                                ? submission.otherUser
                                : submission.matchedUser
                              : "-"}
                          </div>
                        </div>
                      </div>

                      {/* Code Editor */}
                      <div
                        style={{
                          marginTop: "10px",
                          height: "100%",
                        }}
                      >
                        <div
                          ref={editorRef}
                          style={{
                            height: "100%",
                            overflow: "scroll",
                            border: "1px solid #ddd",
                          }}
                        ></div>
                      </div>
                    </div>
                  </div>
                </Row>
              )}
            </Col>
          </Row>
        </Content>
      </Layout>
    </div>
  );
}
