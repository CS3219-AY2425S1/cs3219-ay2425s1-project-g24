"use client";
import Header from "@/components/Header/header";
import {
  Button,
  Col,
  Dropdown,
  Layout,
  message,
  Menu,
  Row,
  TableProps,
  Tag,
  Select,
} from "antd";
import { Content } from "antd/es/layout/layout";
import {
  DeleteOutlined,
  EditOutlined,
  PlusCircleOutlined,
  LeftOutlined,
  RightOutlined,
  CaretRightOutlined,
  ClockCircleOutlined,
  CommentOutlined,
  CheckCircleOutlined,
} from "@ant-design/icons";
import "./styles.scss";
import { useEffect, useState } from "react";
import { Question, GetSingleQuestion } from "../../services/question";
import React from "react";
import TextArea from "antd/es/input/TextArea";
import { useSearchParams } from "next/navigation";
import { ProgrammingLanguageOptions } from "@/utils/SelectOptions";

export default function Home() {
  const [isLoading, setIsLoading] = useState<boolean>(true); // Store the states related to table's loading

  // Message States
  const [messageApi, contextHolder] = message.useMessage();

  const error = (message: string) => {
    messageApi.open({
      type: "error",
      content: message,
    });
  };

  // Retrieve the docRefId from query params during page navigation
  const searchParams = useSearchParams();
  const docRefId: string = searchParams?.get("data") ?? "";

  // Code Editor States
  const [questionTitle, setQuestionTitle] = useState<string | undefined>(
    undefined
  );
  const [complexity, setComplexity] = useState<string | undefined>(undefined);
  const [categories, setCategories] = useState<string[]>([]); // Store the selected filter categories
  const [description, setDescription] = useState<string | undefined>(undefined);
  const [selectedItem, setSelectedItem] = useState("python"); // State to hold the selected language item

  // When code editor page is initialised, fetch the particular question, and display in code editor
  useEffect(() => {
    if (!isLoading) {
      setIsLoading(true);
    }

    GetSingleQuestion(docRefId).then((data: any) => {
      setQuestionTitle(data.title);
      setComplexity(data.complexity);
      setCategories(data.categories);
      setDescription(data.description);
    });
  }, [docRefId]);

  return (
    <div>
      {contextHolder}
      <Layout className="code-editor-layout">
        <Header />
        <Content className="code-editor-content">
          <Row className="entire-page">
            <Col className="col-boxes" span={7}>
              <Row className="problem-description boxes">
                <div className="problem-description-info">
                  <div className="problem-description-top">
                    <h3 className="problem-description-title">
                      {questionTitle}
                    </h3>
                    <text className="problem-solve-status">
                      Solved&nbsp;
                      <CheckCircleOutlined />
                    </text>
                  </div>
                  <div className="complexity-div">
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
                        complexity.charAt(0).toUpperCase() +
                          complexity.slice(1)}
                    </Tag>
                  </div>
                  <div id="tag-container" className="tag-container">
                    <text className="topic-label">Topics: </text>
                    {categories.map((category) => (
                      <Tag key={category}>{category}</Tag>
                    ))}
                  </div>
                  <div className="description-text">
                    <text>{description}</text>
                  </div>
                </div>
              </Row>
              <Row className="test-cases boxes">
                <div className="test-cases-div">
                  <div className="test-cases-top">
                    <h3 className="testcase-title">Testcases</h3>
                    <Button className="runtestcases-button">
                      Run testcases
                      <CaretRightOutlined />
                    </Button>
                  </div>
                  <div className="testcase-buttons">
                    <Button>Case 1</Button>
                    <Button>Case 2</Button>
                    <PlusCircleOutlined />
                  </div>
                  <div className="testcase-code-div">
                    <TextArea
                      className="testcase-code"
                      placeholder="Testcases code"
                    />
                  </div>
                </div>
              </Row>
            </Col>
            <Col className="col-boxes" span={11}>
              <Row className="code-editor boxes">
                <div className="code-editor-div">
                  <div className="code-editor-top">
                    <h3 className="code-editor-title">
                      <LeftOutlined />
                      <RightOutlined style={{ marginRight: "4px" }} />
                      Code
                    </h3>
                    <Button className="submit-solution-button">
                      Submit Solution
                      <CaretRightOutlined />
                    </Button>
                  </div>
                  <div className="language-select">
                    <div>
                      <text className="language-text">
                        Select Language:&nbsp;
                      </text>
                      <Select
                        className="select-language-button"
                        defaultValue={selectedItem}
                        onSelect={(val) => setSelectedItem(val)}
                        options={ProgrammingLanguageOptions}
                      />
                    </div>
                  </div>
                  <div className="code-editor-code-div">
                    <TextArea
                      className="code-editor-code"
                      placeholder="Insert code here"
                    ></TextArea>
                  </div>
                </div>
              </Row>
            </Col>
            <Col span={6} className="col-boxes">
              <Row className="session-details boxes">
                <div className="session-details-div">
                  <div className="session-details-top">
                    <h3 className="session-details-title">
                      <ClockCircleOutlined />
                      &nbsp;Session Details
                    </h3>
                    <Button className="end-session-button">End</Button>
                  </div>
                  <div className="session-details-text-div">
                    <div className="session-details-text">
                      <text className="session-headers">Start Time: </text>
                      01:23:45
                      <br />
                      <text className="session-headers">
                        Session Duration:{" "}
                      </text>
                      01:23:45
                      <br />
                      <text className="session-headers">Matched with: </text>
                      John Doe
                    </div>
                  </div>
                </div>
              </Row>
              <Row className="chat-box boxes">
                <div className="chat-box-div">
                  <div className="chat-box-top">
                    <h3 className="chat-box-title">
                      <CommentOutlined />
                      &nbsp;Chat
                    </h3>
                  </div>
                </div>
              </Row>
            </Col>
          </Row>
        </Content>
      </Layout>
    </div>
  );
}
