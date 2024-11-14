"use client";
import Header from "@/components/Header/header";
import { Button, Layout } from "antd";
import { Content } from "antd/es/layout/layout";
import "./styles.scss";
import { useRouter } from "next/navigation";
import Image from "next/image";
import { useState } from "react";
import MatchingModal from "./matching/MatchingModal";

const HomePage = () => {
  const { push } = useRouter();
  const [matchingModalOpen, setMatchingModalOpen] = useState(false);

  const closeMatchingModal = () => {
    setMatchingModalOpen(false);
  };

  return (
    <Layout className="layout">
      <Header selectedKey={undefined} />
      <Content className="content">
        <MatchingModal isOpen={matchingModalOpen} close={closeMatchingModal} />
        <div className="home-content-card">
          <div className="left-panel">
            <div className="logo-container">
              <div className="logo1-title">Peer</div>
              <div className="logo2-title">Prep</div>
            </div>
            <div className="slogan">
              <span className="slogan-1">
                A better way to prepare for coding interviews with
              </span>
              <span className="slogan-2"> peers</span>
            </div>
            <div className="button-container">
              <Button
                className="match-button"
                onClick={() => setMatchingModalOpen(true)}
              >
                Start Matching
              </Button>
            </div>
          </div>
          <div className="right-panel">
            <Image
              src="/assets/coding.gif"
              alt="my gif"
              height={420}
              width={566}
              unoptimized
            />
          </div>
        </div>
      </Content>
    </Layout>
  );
};

export default HomePage;
