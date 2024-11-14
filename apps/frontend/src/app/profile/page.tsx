"use client";
import Header from "@/components/Header/header";
import {
  Avatar,
  Button,
  Col,
  Form,
  Input,
  Layout,
  message,
  Row,
  Tooltip,
} from "antd";
import { Content } from "antd/es/layout/layout";
import "./styles.scss";
import { EditOutlined, SaveOutlined } from "@ant-design/icons";
import { useEffect, useState } from "react";
import {
  UpdateUser,
  ValidateUser,
  VerifyTokenResponseType,
} from "../services/user";

interface ProfilePageProps {}

const ProfilePage = (props: ProfilePageProps): JSX.Element => {
  const [form] = Form.useForm();
  const [id, setId] = useState<string | undefined>(undefined);
  const [email, setEmail] = useState<string | undefined>(undefined);
  const [username, setUsername] = useState<string | undefined>(undefined);
  const [isEditable, setIsEditable] = useState<boolean>(false);
  const [refresh, setRefresh] = useState<boolean>(false);
  const [messageApi, contextHolder] = message.useMessage();

  useEffect(() => {
    ValidateUser().then((data: VerifyTokenResponseType) => {
      form.setFieldsValue({
        username: data.data.username,
        password: "",
        email: data.data.email,
      });
      setId(data.data.id);
      setEmail(data.data.email);
      setUsername(data.data.username);
    });
  }, [refresh]);

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

  return (
    <Layout className="layout">
      {contextHolder}
      <Header selectedKey={undefined} />
      <Content className="content">
        <div className="content-card1">
          <div className="gradient-header"></div>
          <div className="profile-header">
            <div className="left-header">
              <Avatar size={80} className="avatar">
                {username?.charAt(0).toUpperCase()}
              </Avatar>
              <div className="name-container">
                <div className="username">{username}</div>
                <div className="email">{email}</div>
              </div>
            </div>
            <div className="right-header">
              {!isEditable && (
                <Button
                  icon={<EditOutlined />}
                  className="edit-button"
                  onClick={() => setIsEditable(true)}
                >
                  Edit
                </Button>
              )}
            </div>
          </div>
          <div className="profile-form">
            <Form
              form={form}
              onFinish={(values) => {
                let data = values;
                if (!values.password || values.password.trim() === "") {
                  data = {
                    username: values.username,
                    email: values.email,
                  };
                }
                UpdateUser(data, id as string)
                  .then((value) => {
                    setIsEditable(false);
                    setRefresh(!refresh);
                    success("Profile Updated");
                  })
                  .catch((errors: Error) => {
                    error(errors.message);
                  });
              }}
              layout="vertical"
              disabled={!isEditable}
            >
              <Row gutter={16}>
                <Col span={12}>
                  <Form.Item
                    name="username"
                    label="Username"
                    tooltip="Unmodifable to prevent confusion in other users history when your username changes"
                    rules={[
                      {
                        message: "Please enter valid username!",
                      },
                    ]}
                  >
                    <Input name="username" disabled />
                  </Form.Item>
                </Col>
                <Col span={12}>
                  <Form.Item
                    name="email"
                    label="Email"
                    rules={[
                      {
                        required: true,
                        message: "Please enter valid email!",
                      },
                    ]}
                  >
                    <Input name="email" type="email" />
                  </Form.Item>
                </Col>
              </Row>
              <Row gutter={16}>
                <Col span={12}>
                  <Form.Item
                    name="password"
                    label="Password"
                    tooltip="Password is not updated if this field is left empty"
                  >
                    <Input.Password
                      name="password"
                      type="password"
                      placeholder="Enter new password"
                    />
                  </Form.Item>
                </Col>
              </Row>
              <Form.Item
                style={{ display: "flex", justifyContent: "flex-end" }}
              >
                {isEditable && (
                  <>
                    <Button
                      className="cancel-button"
                      onClick={() => {
                        form.setFieldValue("username", username);
                        form.setFieldValue("email", email);
                        form.setFieldValue("password", undefined);
                        setIsEditable(false);
                      }}
                    >
                      Cancel
                    </Button>
                    <Button
                      icon={<SaveOutlined />}
                      type="primary"
                      className="save-button"
                      htmlType="submit"
                    >
                      Save
                    </Button>
                  </>
                )}
              </Form.Item>
            </Form>
          </div>
        </div>
      </Content>
    </Layout>
  );
};

export default ProfilePage;
