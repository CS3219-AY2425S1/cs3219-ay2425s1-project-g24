import {
  Avatar,
  Button,
  Divider,
  Dropdown,
  Menu,
  MenuProps,
  Space,
} from "antd";
import { Header as AntdHeader } from "antd/es/layout/layout";
import { useRouter } from "next/navigation";
import "./styles.scss";
import DropdownButton from "antd/es/dropdown/dropdown-button";
import { LogoutOutlined, UserOutlined } from "@ant-design/icons";

interface HeaderProps {
  selectedKey: string[] | undefined;
}
const Header = (props: HeaderProps): JSX.Element => {
  
  const { push } = useRouter();
  // Stores the details for the header buttons
  const items = [
    {
      key: 0,
      label: "Problems",
      onClick: () => push("/question"),
    },
    {
      key: 1,
      label: "Matching",
      onClick: () => push("/matching"),
    },
  ];

  const profileItems: MenuProps["items"] = [
    {
      key: 0,
      label: (
        <div className="profile-menu-items">
          <UserOutlined className="profile-menu-icon" /> Profile
        </div>
      ),
      onClick: () => push("/profile"),
    },
    {
      type: "divider",
    },
    {
      key: 2,
      label: (
        <div className="profile-menu-items">
          <LogoutOutlined className="profile-menu-icon" /> Logout
        </div>
      ),
      onClick: () => {
        // Clear away the previously stored jwt token in localstorage
        localStorage.clear();
        // Redirect user to login page
        push("/login");
      },
    },
  ];

  return (
    // Header Component
    <AntdHeader className="header">
      <div className="logo-container-special" onClick={() => push("/")}>
        <div className="logo1">Peer</div>
        <div className="logo2">Prep</div>
      </div>
      {/* Left Menu Pages Component */}
      <Menu
        mode="horizontal"
        defaultSelectedKeys={props.selectedKey}
        items={items}
        style={{ flex: 1, minWidth: 0 }}
      />
      <Dropdown
        menu={{ items: profileItems }}
        placement="bottom"
        overlayClassName="dropdown"
      >
        <a onClick={(e) => e.preventDefault()}>
          {/* Fetch and replace with first letter of user's name */}
          <Avatar icon={<UserOutlined />}></Avatar>
        </a>
      </Dropdown>
    </AntdHeader>
  );
};

export default Header;
