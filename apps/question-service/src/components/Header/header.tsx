import { Menu } from "antd";
import { Header as AntdHeader } from "antd/es/layout/layout";
import { useRouter } from "next/navigation";
import "./styles.scss";

const Header = (): JSX.Element => {
  const { push } = useRouter();
  // Stores the details for the header buttons
  const items = [
    {
      key: 0,
      label: "Problems",
    },
    {
      key: 1,
      label: "Matching",
      disabled: true,
    },
  ];

  return (
    // Header Component
    <AntdHeader className="header">
      <div className="logo-container">
        <div className="logo1">Peer</div>
        <div className="logo2">Prep</div>
      </div>
      <Menu
        mode="horizontal"
        defaultSelectedKeys={["0"]}
        items={items}
        style={{ flex: 1, minWidth: 0 }}
        onClick={(info) => {
          push("/");
        }}
      />
    </AntdHeader>
  );
};

export default Header;
