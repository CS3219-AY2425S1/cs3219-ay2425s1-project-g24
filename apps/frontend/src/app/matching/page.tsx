"use client";

import Header from "@/components/Header/header";
import { 
  Layout, 
  message,
 } from "antd";
import { Content } from "antd/es/layout/layout";
import "./styles.scss";
import HomePage from "@/app/page";
import MatchingModal from "./MatchingModal";

const MatchingPage = () => {
  const [messageApi, contextHolder] = message.useMessage();

  return (
    <div>
      <HomePage />
    </div>
  );
};

export default MatchingPage;
