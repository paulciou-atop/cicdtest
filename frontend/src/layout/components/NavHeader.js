import React from "react";
import AppBar from "@mui/material/AppBar";
import Toolbar from "@mui/material/Toolbar";
import Typography from "@mui/material/Typography";
import {
  Avatar,
  Badge,
  Box,
  IconButton,
  Menu,
  MenuItem,
  Stack,
  Tooltip,
} from "@mui/material";
import { useTheme } from "@mui/material/styles";
import logo from "../../assets/images/logo-new.svg";
import avatarImg from "../../assets/images/81.jpg";
import { useDispatch } from "react-redux";
import { changeThemeMode } from "../../features/ThemeSlice";
import { Brightness4, Brightness7, Notifications } from "@mui/icons-material";
import { useNavigate } from "react-router-dom";

const NavHeader = () => {
  const navigate = useNavigate();
  const theme = useTheme();
  const dispatch = useDispatch();
  const [anchorElUser, setAnchorElUser] = React.useState(null);
  const handleOpenUserMenu = (event) => {
    setAnchorElUser(event.currentTarget);
  };
  const handleCloseUserMenu = () => {
    setAnchorElUser(null);
  };

  const handleLogout = () => {
    localStorage.setItem("token", false);
    navigate("/login");
    setAnchorElUser(null);
  };

  const handleThemeChangeMode = () => {
    const mode = theme.palette.mode === "dark" ? "light" : "dark";
    dispatch(changeThemeMode(mode));
  };

  return (
    <AppBar
      position="fixed"
      sx={{
        zIndex: (theme) => theme.zIndex.drawer + 1,
      }}
      color="inherit"
    >
      <Toolbar>
        <Box component="img" src={logo} alt="logo" width={200} />
        {/* <Typography variant="h6" noWrap component="div" sx={{ ml: 1.5 }}>
          Atech Solution
        </Typography> */}
        <Box sx={{ flexGrow: 1 }} />
        <Stack direction="row" alignItems="center">
          <IconButton size="large" color="inherit">
            <Badge badgeContent={17} color="error">
              <Notifications />
            </Badge>
          </IconButton>

          <IconButton
            size="large"
            onClick={handleThemeChangeMode}
            color="inherit"
          >
            {theme.palette.mode === "dark" ? <Brightness7 /> : <Brightness4 />}
          </IconButton>
          <Tooltip title="Open settings">
            <IconButton
              disableRipple
              onClick={handleOpenUserMenu}
              color="inherit"
            >
              <Avatar alt="Remy Sharp" src={avatarImg} />
              <Typography variant="h6" noWrap component="div" sx={{ ml: 1 }}>
                Admin
              </Typography>
            </IconButton>
          </Tooltip>
          <Menu
            sx={{ mt: "45px" }}
            id="menu-appbar"
            anchorEl={anchorElUser}
            anchorOrigin={{
              vertical: "top",
              horizontal: "right",
            }}
            keepMounted
            transformOrigin={{
              vertical: "top",
              horizontal: "right",
            }}
            open={Boolean(anchorElUser)}
            onClose={handleCloseUserMenu}
          >
            <MenuItem key="logout" onClick={handleLogout}>
              <Typography textAlign="center">Logout</Typography>
            </MenuItem>
          </Menu>
        </Stack>
      </Toolbar>
    </AppBar>
  );
};

export default NavHeader;
