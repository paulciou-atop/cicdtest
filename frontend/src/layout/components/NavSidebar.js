import React from "react";
import Box from "@mui/material/Box";
import Drawer from "@mui/material/Drawer";
import Toolbar from "@mui/material/Toolbar";
import NavSection from "./NavSection";
import navConfig from "./NavConfig";

const NavSidebar = (props) => {
  const { drawerWidth } = props;
  return (
    <Drawer
      elevation={1}
      variant="permanent"
      sx={{
        width: drawerWidth,
        flexShrink: 0,
        [`& .MuiDrawer-paper`]: {
          width: drawerWidth,
          boxSizing: "border-box",
        },
      }}
    >
      <Toolbar />
      <Box sx={{ overflow: "auto", mt: "3px" }}>
        <NavSection navConfig={navConfig} />
      </Box>
    </Drawer>
  );
};

export default NavSidebar;
