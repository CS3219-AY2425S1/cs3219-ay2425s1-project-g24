"use client";
import Header from "@/components/Header/header";
import {
  Button,
  Col,
  Input,
  Layout,
  message,
  Pagination,
  PaginationProps,
  Row,
  Select,
  Table,
  TableProps,
  Tabs,
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
import { useEffect, useState } from "react";
import Link from "next/link";
import TextArea from "antd/es/input/TextArea";
import { createUser } from '@/app/services/user'

type InputFields = {
  username?: string
  email?: string
  password?: string
  confirmPassword?: string
}

export default function Home() {
  function submitDetails({username, email, password}: InputFields): void {
    if (!username || username === "" ||
      !email || email === "" ||
      !password || password === ""
    ) {
      return;
    }
    createUser(username, email, password).catch(err => {
      console.error(err)
    })
  }
  return (
    <div>
      <Layout>
        <Header/>
        <Content>
          <div className="login-card">
            <h1>Register</h1>

            <Form
              name="basic"
              style={{ margin: "auto" }}
              onFinish={submitDetails}
            >
              
              <Form.Item<InputFields>
                name="username"  
                rules={[{required: true}]}
              >
                <Input
                  placeholder="Username"
                />
              </Form.Item>

              <Form.Item<InputFields>
                name="email"  
                rules={[{required: true}]}
              >
                <Input
                  placeholder="Email"
                />
              </Form.Item>

              <Form.Item<InputFields>
                name="password"  
                rules={[{required: true}]}
              >
                <Input.Password
                  placeholder="Password"
                />
              </Form.Item>

              <Form.Item<InputFields>
                name="confirmPassword"
                rules={[{required: true}]}
              >
                <Input.Password
                  placeholder="Confirm Password"
                />
              </Form.Item>

              <Form.Item>
                <Button type="primary" htmlType="submit">
                  Register
                </Button>
              </Form.Item>
            </Form>
            <p>
              Let me <Link
                href="/login"
              >login</Link>
            </p>
          </div>
        </Content>
      </Layout>
    </div>
  );
}