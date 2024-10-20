"use client";
import { Button, Input, Layout, message, Form } from "antd";
import { Content } from "antd/es/layout/layout";
import "./styles.scss";
import { useState } from "react";
import Link from "next/link";
import { createUser } from "@/app/services/user";
import { useRouter } from "next/navigation";
import NoAuthHeader from "@/components/NoAuthHeader/NoAuthHeader";
import { UserOutlined } from "@ant-design/icons";

type InputFields = {
  username: string;
  email: string;
  password: string;
  confirmPassword: string;
};

export default function RegisterPage() {
  const [isRegistrationFailed, setIsRegistrationFailed] = useState(false);
  const router = useRouter();
  const [messageApi, contextHolder] = message.useMessage();

  const successMessage = (message: string) => {
    messageApi.open({
      type: "success",
      content: message,
    });
  };

  function submitDetails({ username, email, password }: InputFields): void {
    createUser(username, email, password)
      .then(() => {
        successMessage("Account successfully created");
        router.push("/login");
      })
      .catch((err) => {
        console.log(err);
        setIsRegistrationFailed(true);
      });
  }
  return (
    <>
      {contextHolder}
      <Layout className="layout">
        <NoAuthHeader />
        <Content className="content">
          <div className="register-card">
            <h1>Register</h1>

            <Form
              name="basic"
              style={{ margin: "auto" }}
              onFinish={submitDetails}
            >
              <Form.Item<InputFields>
                name="username"
                rules={[
                  {
                    required: true,
                    message: "You must provide a username.",
                  },
                  {
                    pattern: /^\S+$/,
                    message: "Please provide a valid username.",
                  },
                ]}
              >
                <Input placeholder="Username" />
              </Form.Item>

              <Form.Item<InputFields>
                name="email"
                rules={[
                  {
                    required: true,
                    message: "You must provide an email.",
                  },
                  {
                    type: "email",
                    message: "Please provide a valid email address.",
                  },
                ]}
              >
                <Input type="email" placeholder="Email" />
              </Form.Item>

              <Form.Item<InputFields>
                name="password"
                rules={[
                  {
                    required: true,
                    message: "You must provide a password.",
                  },
                ]}
              >
                <Input.Password placeholder="Password" />
              </Form.Item>

              <Form.Item<InputFields>
                name="confirmPassword"
                rules={[
                  {
                    required: true,
                    message: "Please confirm your password.",
                  },
                  ({ getFieldValue }) => ({
                    validator: async (r, confirmPassword) => {
                      if (
                        !!confirmPassword &&
                        getFieldValue("password") !== confirmPassword
                      ) {
                        throw new Error("Passwords do not match");
                      }
                    },
                  }),
                ]}
              >
                <Input.Password placeholder="Confirm Password" />
              </Form.Item>

              <div
                style={{
                  visibility: isRegistrationFailed ? "visible" : "hidden",
                }}
                className="registration-failed-text"
              >
                A user with the same email/username already exists
              </div>

              <Form.Item>
                <Button
                  type="primary"
                  htmlType="submit"
                  className="register-button"
                >
                  Register
                </Button>
              </Form.Item>
            </Form>
            <div>
              <p className="login-text">
                Existing user?{" "}
                <Link href="/login" className="login-account-link">
                  Login Here
                </Link>
              </p>
            </div>
          </div>
        </Content>
      </Layout>
    </>
  );
}
