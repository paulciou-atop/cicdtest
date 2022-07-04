import {
  DeviceHub,
  Language,
  Settings,
  SettingsBackupRestore,
  WbIridescent,
} from "@mui/icons-material";
import {
  ListItemIcon,
  ListItemText,
  Menu,
  MenuItem,
  TableRow,
} from "@mui/material";
import React, { useState } from "react";

const DeviceRowContextMenu = ({
  children,
  isItemSelected,
  rowKey,
  data,
  contextMenus,
  onContextMenuClick,
}) => {
  const [contextMenu, setContextMenu] = useState(null);

  const handleContextMenu = (event) => {
    event.preventDefault();
    setContextMenu(
      contextMenu === null
        ? {
            mouseX: event.clientX + 2,
            mouseY: event.clientY - 6,
          }
        : // repeated contextmenu when it is already open closes it with Chrome 84 on Ubuntu
          // Other native context menus might behave different.
          // With this behavior we prevent contextmenu from the backdrop to re-locale existing context menus.
          null
    );
  };

  const handleClose = () => {
    setContextMenu(null);
  };

  const handleMenuItemClick = (event, data, name) => {
    onContextMenuClick(event, data, name);
    handleClose();
  };

  return (
    <TableRow
      style={{ cursor: "context-menu" }}
      key={rowKey}
      hover
      selected={isItemSelected}
      {...(contextMenus && { onContextMenu: handleContextMenu })}
    >
      {children}
      <Menu
        open={contextMenu !== null}
        onClose={handleClose}
        anchorReference="anchorPosition"
        anchorPosition={
          contextMenu !== null
            ? { top: contextMenu.mouseY, left: contextMenu.mouseX }
            : undefined
        }
      >
        <MenuItem
          onClick={(e) => {
            handleMenuItemClick(e, data, "openWeb");
          }}
          key="openWeb"
        >
          <ListItemIcon>
            <Language />
          </ListItemIcon>
          <ListItemText>Open in web</ListItemText>
        </MenuItem>
        <MenuItem
          onClick={(e) => {
            handleMenuItemClick(e, data, "beep");
          }}
          key="beep"
        >
          <ListItemIcon>
            <WbIridescent />
          </ListItemIcon>
          <ListItemText>Beep</ListItemText>
        </MenuItem>
        <MenuItem
          onClick={(e) => {
            handleMenuItemClick(e, data, "reboot");
          }}
          key="reboot"
        >
          <ListItemIcon>
            <SettingsBackupRestore />
          </ListItemIcon>
          <ListItemText>Reboot</ListItemText>
        </MenuItem>
        <MenuItem
          onClick={(e) => {
            handleMenuItemClick(e, data, "networkSetting");
          }}
          key="networkSetting"
        >
          <ListItemIcon>
            <DeviceHub />
          </ListItemIcon>
          <ListItemText>Network Setting</ListItemText>
        </MenuItem>
        <MenuItem
          onClick={(e) => {
            handleMenuItemClick(e, data, "deviceAdvancedSetting");
          }}
          key="deviceAdvancedSetting"
        >
          <ListItemIcon>
            <Settings />
          </ListItemIcon>
          <ListItemText>Device Advanced Setting</ListItemText>
        </MenuItem>
      </Menu>
    </TableRow>
  );
};

export default DeviceRowContextMenu;
