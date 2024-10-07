"use client";
import Header from "@/components/Header/header";
import { Button, Input, Layout, message, Form } from "antd";
import { Content } from "antd/es/layout/layout";
import "./styles.scss";
import { useState } from "react";
import Link from "next/link";
import { loginUser } from "@/app/services/user";
import { setToken } from "@/app/services/login-store";
import { useRouter } from "next/navigation";

type InputFields = {
  email: string;
  password: string;
};

export default function Home() {
  const [isLoginFailed, setIsLoginFailed] = useState(false);
  const router = useRouter();
  const [messageApi, contextHolder] = message.useMessage();

  const successMessage = (message: string) => {
    messageApi.open({
      type: "success",
      content: message,
    });
  };

  function submitDetails({ email, password }: InputFields): void {
    loginUser(email, password)
      .then((jwt) => {
        successMessage("Login successful");
        setToken(jwt);
        router.push("/");
      })
      .catch((err) => {
        console.error(err);
        setIsLoginFailed(true);
      });
  }

  return (
    <>
      {contextHolder}
      <Layout>
        <Header selectedKey={undefined} />
        <Content>
          <div className="login-card">
            <h1>Login</h1>
            <Form name="basic" onFinish={submitDetails}>
              <Form.Item<InputFields>
                name="email"
                rules={[
                  {
                    required: true,
                    message: "Please provide your email.",
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
                    message: "Please enter your password.",
                  },
                ]}
              >
                <Input.Password placeholder="Password" />
              </Form.Item>

              <div
                style={{ visibility: isLoginFailed ? "visible" : "hidden" }}
                className="registration-failed-text"
              >
                This email/username is not valid
              </div>

              <Form.Item>
                <Button type="primary" htmlType="submit">
                  Login
                </Button>
              </Form.Item>
            </Form>
            <p>
              Let me <Link href="/register">register</Link>
            </p>
          </div>
        </Content>
      </Layout>
    </>
  );
}
