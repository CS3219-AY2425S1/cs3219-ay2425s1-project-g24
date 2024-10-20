"use client";
import Header from "@/components/Header/header";
import {
  Button,
  Col,
  Input,
  Layout,
  message,
  PaginationProps,
  Row,
  Select,
  Table,
  TableProps,
  Tag,
  Modal,
  Form,
} from "antd";
import { Content } from "antd/es/layout/layout";
import {
  DeleteOutlined,
  EditOutlined,
  PlusCircleOutlined,
  SearchOutlined,
} from "@ant-design/icons";
import "./styles.scss";
import { useEffect, useState, useLayoutEffect } from "react";
import {
  DeleteQuestion as DeleteQuestionByDocref,
  GetQuestions,
  Question,
  CreateQuestion,
  NewQuestion,
  EditQuestion,
} from "../services/question";
import {
  CategoriesOption,
  DifficultyOption,
  OrderOption,
} from "../../utils/SelectOptions";
import Link from "next/link";
import TextArea from "antd/es/input/TextArea";
import { ValidateUser, VerifyTokenResponseType } from "../services/user";
import { useRouter } from "next/navigation";

/**
 * defines the State of the page whe a user is deleing an object. Has 3 general states:
 * - {}: user is not deleting anything. The page's normal state
 * - {index: docref string, deleteConfirmed: false}: Modal popup asking whether to delete the question, pending user's decision to confirm or cancel
 * - {index: docref string, deleteConfirmed: true}: Currently deleting the question and reloading the database
 */
type DeletionStage = {} | { index: Question; deleteConfirmed: boolean };

function DeleteModal({
  isDeleting,
  questionTitle,
  okHandler,
  cancelHandler,
}: {
  questionTitle: string;
  okHandler: () => void;
  cancelHandler: () => void;
  isDeleting: boolean;
}) {
  const title: string = `Delete Question \"${questionTitle}\"?`;
  const text: string = "This action is irreversible(?)!";

  return (
    <Modal
      open={true}
      title={title}
      onOk={okHandler}
      onCancel={cancelHandler}
      confirmLoading={isDeleting}
      okButtonProps={{ danger: true }}
      cancelButtonProps={{ disabled: isDeleting }}
    >
      <p>{text}</p>
    </Modal>
  );
}

export default function Home() {
  // State of Deletion
  const [deletionStage, setDeletionStage] = useState<DeletionStage>({});

  // Table States
  const [questions, setQuestions] = useState<Question[] | undefined>(undefined); // Store the questions
  const [totalCount, setTotalCount] = useState<number | undefined>(undefined); // Store the total count of questions
  const [totalPages, setTotalPages] = useState<number | undefined>(undefined); // Store the total number of pages
  const [currentPage, setCurrentPage] = useState<number>(1); // Store the current page
  const [limit, setLimit] = useState<number>(10); // Store the quantity of questions to be displayed
  const [isLoading, setIsLoading] = useState<boolean>(true); // Store the states related to table's loading

  // Filtering States
  const [search, setSearch] = useState<string | undefined>(undefined); // Store the search
  const [delayedSearch, setDelayedSearch] = useState<string | undefined>(
    undefined
  ); // Store the delayed search value
  const [categories, setCategories] = useState<string[]>([]); // Store the selected filter categories
  const [difficulty, setDifficulty] = useState<string[]>([]); // Store the selected difficulty level
  const [sortBy, setSortBy] = useState<string | undefined>(undefined); // Store the selected sorting parameter

  // Message States
  const [messageApi, contextHolder] = message.useMessage();

  // States for Create New Problem Modal
  const [form] = Form.useForm();
  const [isNewProblemModalOpen, setIsNewProblemModelOpen] = useState(false);

  // States for Edit Existing Problem Modal
  const [editForm] = Form.useForm();
  const [isEditModalOpen, setIsEditModalOpen] = useState<boolean[] | undefined>(
    undefined
  );

  // State for refreshing data
  const [refresh, setRefresh] = useState(false);

  // used to check if user JWT is verified

  const [userId, setUserId] = useState<string | undefined>(undefined);
  const [email, setEmail] = useState<string | undefined>(undefined);
  const [username, setUsername] = useState<string | undefined>(undefined);
  const [isAdmin, setIsAdmin] = useState<boolean | undefined>(undefined);

  const router = useRouter();

  useLayoutEffect(() => {
    var isAuth = false;

    ValidateUser()
      .then((data: VerifyTokenResponseType) => {
        setUserId(data.data.id);
        setEmail(data.data.email);
        setUsername(data.data.username);
        setIsAdmin(data.data.isAdmin);
        isAuth = true;
      })
      .finally(() => {
        if (!isAuth) {
          // cannot verify
          router.push("/login"); // Client-side redirect using router.push
        }
      });
  }, [router]);

  const handleEditClick = (index: number, question: Question) => {
    // Open the modal for the specific question
    const updatedModals =
      isEditModalOpen && isEditModalOpen.map((_, idx) => idx === index);
    setIsEditModalOpen(updatedModals); // Only the selected question's modal is open

    // Set the form value
    editForm.setFieldsValue({
      title: question.title,
      description: question.description,
      complexity: question.complexity,
      categories: question.categories,
    });
  };

  // Function to handle modal close
  const handleModalClose = (index: number) => {
    if (isEditModalOpen) {
      const updatedModals = [...isEditModalOpen];
      updatedModals[index] = false; // Close the specific modal
      setIsEditModalOpen(updatedModals);
    }
  };

  const handleEditQuestion = async (
    values: NewQuestion,
    index: number,
    docRefId: string
  ) => {
    try {
      const editedQuestion = await EditQuestion(values, docRefId);
      // Reset form or update UI as needed
      handleModalClose(index);
      editForm.resetFields();
      success("Problem Updated!");
      setRefresh(!refresh);
    } catch (err: any) {
      error(err.message);
    }
  };

  const handleCreateQuestion = async (values: NewQuestion) => {
    try {
      const createdQuestion = await CreateQuestion(values);
      // Reset form or update UI as needed
      setIsNewProblemModelOpen(false);
      form.resetFields();
      success("New Problem Created!");
      setRefresh(!refresh);
    } catch (err: any) {
      error(err.message);
    }
  };

  const showNewProblemModal = () => {
    setIsNewProblemModelOpen(true);
  };

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

  async function loadQuestions() {
    if (!isLoading) {
      setIsLoading(true);
    }

    let data = await GetQuestions(
      currentPage,
      limit,
      sortBy,
      difficulty,
      categories,
      delayedSearch
    );
    setQuestions(data.questions);
    setTotalCount(data.totalCount);
    setTotalPages(data.totalPages);
    setCurrentPage(data.currentPage);
    setLimit(data.limit);
    setIsLoading(false);
    setIsEditModalOpen(Array(data.questions.length).fill(false));
  }

  useEffect(() => {
    loadQuestions();
  }, [
    limit,
    currentPage,
    sortBy,
    difficulty,
    categories,
    delayedSearch,
    refresh,
  ]);

  // Delay the fetching of data only after user stops typing for awhile
  useEffect(() => {
    const timeout = setTimeout(() => {
      setDelayedSearch(search);
      setCurrentPage(1); // Reset the current page
    }, 500);
    return () => clearTimeout(timeout);
  }, [search]);

  // Table column specification
  var columns: TableProps<Question>["columns"];
  if (isAdmin) {
    columns = [
      {
        title: "Id",
        dataIndex: "id",
        key: "id",
        // render: (id: number) => <div>{id}</div>,
      },
      {
        title: "Title",
        dataIndex: "title",
        key: "title",
        render: (text: string, question: Question) => (
          <Link
            href={{
              pathname: `/question/${question.id}`,
              query: { data: question.docRefId }, // the data
            }}
          >
            <Button type="link">{text}</Button>
          </Link>
        ),
      },
      {
        title: "Categories",
        dataIndex: "categories",
        key: "categories",
        render: (categories: string[]) =>
          categories.map((category) => <Tag key={category}>{category}</Tag>),
      },
      {
        title: "Difficulty",
        dataIndex: "complexity",
        key: "complexity",
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
        title: "Actions",
        key: "actions",
        dataIndex: "id",
        render: (_: number, question: Question, index: number) => (
          <div>
            <Modal
              title="Edit Problem"
              open={isEditModalOpen && isEditModalOpen[index]}
              onCancel={() => handleModalClose(index)}
              footer={null}
              width={600}
            >
              <Form
                name="edit-form"
                {...layout}
                form={editForm}
                onFinish={(values) => {
                  handleEditQuestion(values, index, question.docRefId);
                }}
              >
                <Form.Item
                  name="title"
                  label="Title"
                  rules={[
                    {
                      required: true,
                      message: "Please enter question title!",
                    },
                  ]}
                >
                  <Input name="title" />
                </Form.Item>
                <Form.Item
                  name="description"
                  label="Description"
                  rules={[
                    {
                      required: true,
                      message: "Please enter question description!",
                    },
                  ]}
                >
                  <TextArea name="description" />
                </Form.Item>
                <Form.Item
                  name="complexity"
                  label="Complexity"
                  rules={[
                    {
                      required: true,
                      message: "Please select a complexity!",
                    },
                  ]}
                >
                  <Select
                    options={[
                      {
                        label: "Easy",
                        value: "easy",
                      },
                      {
                        label: "Medium",
                        value: "medium",
                      },
                      {
                        label: "Hard",
                        value: "hard",
                      },
                    ]}
                    onChange={(value) =>
                      form.setFieldValue("complexity", value)
                    }
                    allowClear
                  />
                </Form.Item>
                <Form.Item
                  name="categories"
                  label="Categories"
                  rules={[
                    {
                      required: true,
                      message: "Please select the relevant categories!",
                    },
                  ]}
                >
                  <Select
                    mode="multiple"
                    options={CategoriesOption}
                    onChange={(value) =>
                      form.setFieldValue("categories", value)
                    }
                    allowClear
                  />
                </Form.Item>
                <Form.Item
                  style={{ display: "flex", justifyContent: "flex-end" }}
                >
                  <Button type="primary" htmlType="submit">
                    Save
                  </Button>
                </Form.Item>
              </Form>
            </Modal>
            <Button
              className="edit-button"
              icon={<EditOutlined />}
              onClick={() => handleEditClick(index, question)}
            ></Button>
            <Button
              className="delete-button"
              danger
              icon={<DeleteOutlined />}
              onClick={() => {
                setDeletionStage({ index: question, deleteConfirmed: false });
              }}
            ></Button>
          </div>
        ),
      },
    ];
  } else {
    columns = [
      {
        title: "Id",
        dataIndex: "id",
        key: "id",
        // render: (id: number) => <div>{id}</div>,
      },
      {
        title: "Title",
        dataIndex: "title",
        key: "title",
        render: (text: string, question: Question) => (
          <Link
            href={{
              pathname: `/question/${question.id}`,
              query: { data: question.docRefId }, // the data
            }}
          >
            <Button type="link">{text}</Button>
          </Link>
        ),
      },
      {
        title: "Categories",
        dataIndex: "categories",
        key: "categories",
        render: (categories: string[]) =>
          categories.map((category) => <Tag key={category}>{category}</Tag>),
      },
      {
        title: "Difficulty",
        dataIndex: "complexity",
        key: "complexity",
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
    ];
  }

  // Handler for change in multi-select categories option
  const handleCategoriesChange = (value: string[]) => {
    setCategories(value);
    setCurrentPage(1); // Reset the current page
  };

  // Handler for clearing the filtering options
  const handleClear = () => {
    setCategories([]);
    setDifficulty([]);
    setSearch(undefined);
    setSortBy(undefined);
  };

  // Handler for filtering (TODO)
  // const handleFilter = async () => {
  //   success("Filtered Successfully!");
  // };

  // Handler for show size change for pagination
  const onShowSizeChange: PaginationProps["onShowSizeChange"] = (
    current,
    pageSize
  ) => {
    setCurrentPage(current);
    setLimit(pageSize);
  };

  // Handler for change in page jumper
  const onPageJump: PaginationProps["onChange"] = (pageNumber) => {
    setCurrentPage(pageNumber);
  };

  const confirmDeleteHandler = async () => {
    if (!("index" in deletionStage) || deletionStage.deleteConfirmed) {
      error("Cannot delete: invalid deletionStage");
      return;
    }
    if (questions == undefined) {
      error("Cannot delete: questions does not exist");
      return;
    }

    setDeletionStage({ index: deletionStage.index, deleteConfirmed: true });
    await DeleteQuestionByDocref(deletionStage.index.docRefId);
    if (questions.length == 1 && currentPage > 1) {
      setCurrentPage(currentPage - 1);
      success("Question deleted successfully");
    } else {
      try {
        await loadQuestions();
        success("Question deleted successfully");
      } catch (err) {
        if (typeof err == "string") {
          error(err);
        }
      }
    }
    setDeletionStage({});
  };

  const layout = {
    labelCol: { span: 4 },
    wrapperCol: { span: 20 },
  };

  return (
    <div>
      {contextHolder}
      <Layout className="layout">
        <Header selectedKey={["0"]} />
        <Content className="content">
          <div className="content-card">
            <div className="content-row-1">
              <div className="content-title">Problems</div>
              {isAdmin && (
                <div className="create-button">
                  <Button
                    type="primary"
                    icon={<PlusCircleOutlined />}
                    onClick={showNewProblemModal}
                  >
                    Create New Problem
                  </Button>
                  <Modal
                    title="Create New Problem"
                    open={isNewProblemModalOpen}
                    // onOk={() => setIsNewProblemModelOpen(false)} // Replace with handleSubmit
                    onCancel={() => setIsNewProblemModelOpen(false)}
                    footer={null}
                    width={600}
                  >
                    <Form
                      name="create-form"
                      {...layout}
                      form={form}
                      onFinish={(values) => {
                        handleCreateQuestion(values);
                      }}
                    >
                      <Form.Item
                        name="title"
                        label="Title"
                        rules={[
                          {
                            required: true,
                            message: "Please enter question title!",
                          },
                        ]}
                      >
                        <Input name="title" />
                      </Form.Item>
                      <Form.Item
                        name="description"
                        label="Description"
                        rules={[
                          {
                            required: true,
                            message: "Please enter question description!",
                          },
                        ]}
                      >
                        <TextArea name="description" />
                      </Form.Item>
                      <Form.Item
                        name="complexity"
                        label="Complexity"
                        rules={[
                          {
                            required: true,
                            message: "Please select a complexity!",
                          },
                        ]}
                      >
                        <Select
                          options={[
                            {
                              label: "Easy",
                              value: "easy",
                            },
                            {
                              label: "Medium",
                              value: "medium",
                            },
                            {
                              label: "Hard",
                              value: "hard",
                            },
                          ]}
                          onChange={(value) =>
                            form.setFieldValue("complexity", value)
                          }
                          allowClear
                        />
                      </Form.Item>
                      <Form.Item
                        name="categories"
                        label="Categories"
                        rules={[
                          {
                            required: true,
                            message: "Please select the relevant categories!",
                          },
                        ]}
                      >
                        <Select
                          mode="multiple"
                          options={CategoriesOption}
                          onChange={(value) =>
                            form.setFieldValue("categories", value)
                          }
                          allowClear
                        />
                      </Form.Item>
                      <Form.Item
                        style={{ display: "flex", justifyContent: "flex-end" }}
                      >
                        <Button type="primary" htmlType="submit">
                          Create
                        </Button>
                      </Form.Item>
                    </Form>
                  </Modal>
                </div>
              )}
            </div>
            <div className="content-filter">
              <Row gutter={8}>
                <Col span={6}>
                  <Input
                    placeholder="Search Question Title"
                    prefix={<SearchOutlined />}
                    onChange={(e) => setSearch(e.target.value)}
                    value={search}
                    allowClear
                  />
                </Col>
                <Col span={6}>
                  <Select
                    mode="multiple"
                    allowClear
                    placeholder="Categories"
                    onChange={handleCategoriesChange}
                    options={CategoriesOption}
                    className="categories-multi-select"
                    value={categories}
                  />
                </Col>
                <Col span={6}>
                  <Select
                    mode="multiple"
                    allowClear
                    placeholder="Difficulty"
                    onChange={(value: string[]) => {
                      setDifficulty(value);
                      setCurrentPage(1); //Reset the currentpage since filter params changed
                    }}
                    options={DifficultyOption}
                    className="difficulty-select"
                    value={difficulty}
                  />
                </Col>
                <Col span={4}>
                  <Select
                    allowClear
                    placeholder="Sort By"
                    onChange={(value: string) => setSortBy(value)}
                    options={OrderOption}
                    className="order-select"
                    value={sortBy}
                  />
                </Col>
                <Col span={2}>
                  <Button onClick={handleClear} className="clear-button">
                    Clear
                  </Button>
                  {/* <Button
                    type="primary"
                    className="filter-button"
                    onClick={handleFilter}
                  >
                    Filter
                  </Button> */}
                </Col>
              </Row>
            </div>
            <div className="content-table">
              <Table
                rowKey="id"
                dataSource={questions}
                columns={columns}
                loading={isLoading}
                pagination={{
                  current: currentPage,
                  total: totalCount,
                  showSizeChanger: true,
                  onShowSizeChange: onShowSizeChange,
                  // showQuickJumper: true,
                  onChange: onPageJump,
                }}
              />
            </div>
          </div>
        </Content>
      </Layout>
      {"index" in deletionStage && questions != undefined && (
        <DeleteModal
          okHandler={confirmDeleteHandler}
          cancelHandler={() => setDeletionStage({})}
          questionTitle={deletionStage.index.title}
          isDeleting={deletionStage.deleteConfirmed}
        />
      )}
    </div>
  );
}
