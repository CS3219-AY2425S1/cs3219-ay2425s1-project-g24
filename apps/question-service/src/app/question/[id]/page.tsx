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
import {
  CategoriesOption,
  DifficultyOption,
  OrderOption,
} from "../utils/SelectOptions";
import React from "react";
import TextArea from "antd/es/input/TextArea";
import { useSearchParams } from "next/navigation";

export default function Home() {
  // Table States
  const [currentPage, setCurrentPage] = useState<number | undefined>(1); // Store the current page
  const [isLoading, setIsLoading] = useState<boolean>(true); // Store the states related to table's loading

  // Filtering States
  const [search, setSearch] = useState<string | undefined>(undefined); // Store the search
  const [difficulty, setDifficulty] = useState<string[]>([]); // Store the selected difficulty level
  const [sortBy, setSortBy] = useState<string | undefined>(undefined); // Store the selected sorting parameter

  // Message States
  const [messageApi, contextHolder] = message.useMessage();

  // Code Editor States
  const [questionTitle, setQuestionTitle] = useState<string | undefined>(undefined);
  const [complexity, setComplexity] = useState<string | undefined>(undefined);
  const [categories, setCategories] = useState<string[]>([]); // Store the selected filter categories
  const [description, setDescription] = useState<string | undefined>(undefined);

  const success = (message: string) => {
    messageApi.open({
      type: "success",
      content: message,
    });
  };

  const error = (message: string) => {
    messageApi.open({
      type: "error",
      content: message,
    });
  };

  const warning = (message: string) => {
    messageApi.open({
      type: "warning",
      content: message,
    });
  };

  // Include States for Create/Edit Modal (TODO: Sean)
  const searchParams = useSearchParams();
  const docRefId = searchParams.get("data")?.toString();

  // When code editor page is initialised, fetch the particular question, and display in code editor
  useEffect(() => {
    if (!isLoading) {
      setIsLoading(true);
    }

    GetSingleQuestion(docRefId).then(
      (data: any) => {
        setQuestionTitle(data.title);
        setComplexity(data.complexity);
        setCategories(data.categories);
        setDescription(data.description);
      }
    );
  }, [docRefId])

  const menuItems = [
    {
      key: '1',
      label: 'Python',
    },
    {
      key: '2',
      label: 'Java',
    },
    {
      key: '3',
      label: 'C++',
    },
  ];

  // State to hold the selected menu item
  const [selectedItem, setSelectedItem] = useState('Python');

  // Function to handle the item selection
  const handleMenuClick = (e: any) => {
    const selected = menuItems.find((item) => item.key === e.key);
    if (selected) {
      setSelectedItem(selected.label);
    }
  };

    // Menu component
    const menu = (
      <Menu onClick={handleMenuClick}>
        {menuItems.map((item) => (
          <Menu.Item key={item.key}>{item.label}</Menu.Item>
        ))}
      </Menu>
    );

  return (
    <div>
      {contextHolder}
      <Layout className="code-editor-layout">
        <Header />
        <Content className="code-editor-content">
            <Row className="entire-page">
                <Col className="col-boxes" span={8}>
                    <Row className="problem-description boxes">
                      <text className="problem-description-info">
                        <div className="problem-description-top">
                          <h3 className="problem-description-title">{questionTitle}</h3>
                          <text className="problem-solve-status">Solved&nbsp;<CheckCircleOutlined/></text>
                        </div>
                        <Tag>{complexity}</Tag>
                        <div id="tag-container">
                          <text>Topics: </text>
                          {categories.map(category => (
                              <Tag key={category}>{category}</Tag>
                          ))}
                        </div>
                        <text>{description}</text>
                      </text>
                    </Row>
                    <Row className="test-cases boxes">
                      <div className="test-cases-div">
                        <div className="test-cases-top">
                          <h3 className="testcase-title">Testcases</h3>
                          <Button className="runtestcases-button">Run testcases<CaretRightOutlined/></Button>
                        </div>
                        <div className="testcase-buttons">
                          <Button>Case 1</Button>
                          <Button>Case 2</Button>
                          <PlusCircleOutlined/>
                        </div>
                        <div className="testcase-code-div">
                          <TextArea className="testcase-code" placeholder="Testcases code"/>
                        </div>
                      </div>
                    </Row>
                </Col>
                <Col className="code-editor-box boxes col-boxes" span={11}>
                  <div className="code-editor-div">
                    <div className="code-editor-top">
                      <h3 className="code-editor-title"><LeftOutlined/><RightOutlined/>Code</h3>
                      <Button className="submit-solution-button">Submit Solution<CaretRightOutlined/></Button>
                    </div>
                    <div className="language-select">
                      <text>Select Language:&nbsp;
                        <Dropdown className="select-language-button" overlay={menu} trigger={['click']}>
                          <Button>{selectedItem}</Button>
                        </Dropdown>
                      </text>
                    </div>
                    <div className="code-editor-code-div">
                      <TextArea className="code-editor-code" placeholder="Insert code here"></TextArea>
                    </div>
                  </div>
                </Col>
                <Col span={5} className="col-boxes">
                    <Row className="session-details boxes">
                      <div className="session-details-div">
                        <div className="session-details-top">
                          <h3 className="session-details-title"><ClockCircleOutlined/>&nbsp;Session Details</h3>
                          <Button className="end-session-button">End</Button>
                        </div>
                        <div className="session-details-text-div">
                          <text className="session-details-text">
                            Start Time: 01:23:45<br/>
                            Session Duration: 01:23:45<br/>
                            Matched with John Doe
                          </text>
                        </div>
                      </div>
                    </Row>
                    <Row className="chat-box boxes">
                      <div className="chat-box-div">
                        <div className="chat-box-top">
                          <h3 className="chat-box-title"><CommentOutlined/>&nbsp;Chat</h3>
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
