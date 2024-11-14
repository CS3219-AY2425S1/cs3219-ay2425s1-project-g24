import { Menu } from "antd";
import { Header as AntdHeader } from "antd/es/layout/layout";
import { useRouter } from "next/navigation";
import "./styles.scss";

const NoAuthHeader = (): JSX.Element => {
  const { push } = useRouter();
  // Disable the header for non-authenticated header
  const items = [
    {
      key: 0,
      label: "Problems",
      disabled: true,
    },
    {
      key: 1,
      label: "Matching",
      disabled: true,
    },
  ];

  return (
    // Header Component
    <AntdHeader className="noauth-header">
      <div className="logo-container">
        <div className="logo1">Peer</div>
        <div className="logo2">Prep</div>
      </div>
      {/* Left Menu Pages Component */}
      <Menu
        mode="horizontal"
        items={items}
        style={{ flex: 1, minWidth: 0 }}
        onClick={(info) => {
          push("/");
        }}
      />
    </AntdHeader>
  );
};

export default NoAuthHeader;
