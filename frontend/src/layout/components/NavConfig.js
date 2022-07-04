import {
  DashboardCustomizeOutlined,
  DashboardOutlined,
  ManageAccountsOutlined,
  Inventory2Outlined,
} from "@mui/icons-material";

const navConfig = [
  {
    title: "dashboard",
    path: "/dashboard",
    icon: <DashboardOutlined />,
    children: [
      {
        title: "overview",
        path: "/dashboard/overview",
      },
      {
        title: "devices",
        path: "/dashboard/devices",
      },
    ],
  },
  {
    title: "management",
    path: "/management",
    icon: <DashboardCustomizeOutlined />,
    children: [
      {
        title: "devices",
        path: "/management/devices",
      },
      {
        title: "devices Configuration",
        path: "/management/deviceConfig",
      },
      {
        title: "users",
        path: "/management/users",
      },
    ],
  },
  {
    title: "inventory",
    path: "/inventory",
    icon: <Inventory2Outlined />,
    children: [
      {
        title: "devices",
        path: "/inventory/devices",
      },
    ]
  },
  {
    title: "account",
    path: "/account",
    icon: <ManageAccountsOutlined />,
  },
];

export default navConfig;
