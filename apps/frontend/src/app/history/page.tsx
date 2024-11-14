"use client";
import Header from "@/components/Header/header";
import { Layout, message, PaginationProps, Row, Table, Tag } from "antd";
import { Content } from "antd/es/layout/layout";
import { HistoryOutlined } from "@ant-design/icons";
import "./styles.scss";
import { useEffect, useState } from "react";
import React from "react";
import { useRouter } from "next/navigation";
import { GetUserHistories, History } from "@/app/services/history";
import { ValidateUser, VerifyTokenResponseType } from "@/app/services/user";

interface TablePagination {
  totalCount: number;
  totalPages: number;
  currentPage: number;
  limit: number;
}

export default function QuestionPage() {
  // Message States
  const [messageApi, contextHolder] = message.useMessage();

  const error = (message: string) => {
    messageApi.open({
      type: "error",
      content: message,
    });
  };

  const router = useRouter();

  const [username, setUsername] = useState<string | undefined>(undefined);
  const [userQuestionHistories, setUserQuestionHistories] =
    useState<History[]>();
  const [isHistoryLoading, setIsHistoryLoading] = useState<boolean>(true);
  const [paginationParams, setPaginationParams] = useState<TablePagination>({
    totalCount: 0,
    totalPages: 0,
    currentPage: 1,
    limit: 10,
  });

  // Handler for change in page jumper
  const onPageJump: PaginationProps["onChange"] = (pageNumber) => {
    setPaginationParams((prev) => {
      loadQuestionHistories(pageNumber, prev.limit);
      return { ...paginationParams, currentPage: pageNumber };
    });
  };

  // Handler for show size change for pagination
  const onShowSizeChange: PaginationProps["onShowSizeChange"] = (
    current,
    pageSize
  ) => {
    setPaginationParams((prev) => {
      loadQuestionHistories(current, pageSize);
      return { ...paginationParams, currentPage: current, limit: pageSize };
    });
  };

  async function loadQuestionHistories(currentPage: number, limit: number) {
    if (username === undefined) return;
    setIsHistoryLoading(true);
    GetUserHistories(username, currentPage, limit)
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

  useEffect(() => {
    loadQuestionHistories(paginationParams.currentPage, paginationParams.limit);
  }, [username]);

  useEffect(() => {
    ValidateUser().then((data: VerifyTokenResponseType) => {
      setUsername(data.data.username);
    });
  }, []);

  const columns = [
    {
      title: "Question Title",
      dataIndex: "title",
      key: "title",
    },
    {
      title: "Categories",
      dataIndex: "questionTopics",
      key: "questionTopics",
      render: (categories: string[]) =>
        categories.map((category) => <Tag key={category}>{category}</Tag>),
    },
    {
      title: "Difficulty",
      dataIndex: "questionDifficulty",
      key: "questionDifficulty",
      render: (difficulty: string) => {
        let color = "";
        if (difficulty === "easy") {
          color = "#2DB55D";
        } else if (difficulty === "medium") {
          color = "orange";
        } else if (difficulty === "hard") {
          color = "red";
        }
        return (
          <div style={{ color }}>
            {difficulty.charAt(0).toUpperCase() + difficulty.slice(1)}
          </div>
        );
      },
    },
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

  const handleRowClick = (h: History) => {
    // Link to page
    // questionId is just read as "history", as only the doc ref id is involved in requests
    // If the question database is reloaded, then the questionDocRefId may not be correct
    router.push(
      `/question/history?data=${h.questionDocRefId}&history=${h.historyDocRefId}`
    );
  };

  return (
    <div>
      {contextHolder}
      <Layout className="layout">
        <Header selectedKey={["0"]} />
        <Content className="content">
          <div className="content-card">
            <div className="content-row-1" style={{ marginBottom: "20px" }}>
              <div className="content-title">Submission History</div>
            </div>
            <div className="content-table">
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
                  current: paginationParams.currentPage,
                  total: paginationParams.totalCount,
                  pageSize: paginationParams.limit,
                  onChange: onPageJump,
                  showSizeChanger: true,
                  onShowSizeChange: onShowSizeChange,
                }}
              />
            </div>
          </div>
        </Content>
      </Layout>
    </div>
  );
}
