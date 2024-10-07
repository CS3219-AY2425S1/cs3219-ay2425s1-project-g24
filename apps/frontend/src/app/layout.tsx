import React from "react";
import { AntdRegistry } from "@ant-design/nextjs-registry";
import { ConfigProvider } from "antd";

const RootLayout = ({
  children,
}: React.PropsWithChildren<React.ReactNode | undefined>) => (
  <html lang="en">
    <body>
      <ConfigProvider
        theme={{
          components: {
            Button: {
              colorPrimary: "#463F3A",
              algorithm: true, // Enable algorithm
            },
            Input: {
              colorBgContainer: "#FAFAFA",
            },
            Select: {
              colorBgContainer: "#FAFAFA",
            },
          },
        }}
      >
        <AntdRegistry>{children}</AntdRegistry>
      </ConfigProvider>
    </body>
  </html>
);

export default RootLayout;
