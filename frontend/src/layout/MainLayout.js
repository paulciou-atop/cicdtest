import React from "react";
import { Outlet } from "react-router-dom";
import Box from "@mui/material/Box";
import Toolbar from "@mui/material/Toolbar";
import NavSidebar from "./components/NavSidebar";
import NavHeader from "./components/NavHeader";
import BreadcrumItemLayout from "./components/BreadcrumItemLayout";
import NetworkSettingDrawer from "../components/drawer/NetworkSettingDrawer";
import AdvancedSettingDrawer from "../components/drawer/AdvancedSettingDrawer";

const drawerWidth = 240;

const MainLayout = () => {
  return (
    <Box sx={{ display: "flex" }}>
      <NavHeader />
      <NavSidebar drawerWidth={drawerWidth} />
      <Box component="main" sx={{ flexGrow: 1, p: 3, pt: 0 }}>
        <Toolbar />
        <BreadcrumItemLayout />
        <Outlet />
        <NetworkSettingDrawer />
        <AdvancedSettingDrawer />
      </Box>
    </Box>
  );
};

export default MainLayout;
