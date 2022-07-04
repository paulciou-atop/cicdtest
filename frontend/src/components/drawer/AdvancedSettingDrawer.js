import {
  Button,
  Divider,
  Drawer,
  Paper,
  Stack,
  Toolbar,
  Typography,
} from "@mui/material";
import Box from "@mui/material/Box";
import Tab from "@mui/material/Tab";
import TabContext from "@mui/lab/TabContext";
import TabList from "@mui/lab/TabList";
import TabPanel from "@mui/lab/TabPanel";
import React, { useState } from "react";
import { useDispatch, useSelector } from "react-redux";
import {
  closeAdvancedSettingDrawer,
  deviceAdvancedSettingSelector,
} from "../../features/DeviceAdvancedSetting";
import Controls from "../common/form/controls/Controls";
import AlarmPortList from "./component/AlarmPortList";
import AlarmPowerList from "./component/AlarmPowerList";

const drawerWidth = 400;
const portData = [
  { portName: "Port 1" },
  { portName: "Port 2" },
  { portName: "Port 3" },
  { portName: "Port 4" },
  { portName: "Port 5" },
  { portName: "Port 6" },
  { portName: "Port 7" },
  { portName: "Port 8" },
];
const powerData = [{ powerName: "Power 1" }, { powerName: "Power 2" }];

const AdvancedSettingDrawer = () => {
  const dispatch = useDispatch();
  const { visible } = useSelector(deviceAdvancedSettingSelector);

  const [value, setValue] = useState("1");

  const handleChange = (event, newValue) => {
    setValue(newValue);
  };

  const handleSaveClick = () => {
    dispatch(closeAdvancedSettingDrawer());
  };
  const handleCancleClick = () => {
    dispatch(closeAdvancedSettingDrawer());
  };

  return (
    <Drawer
      anchor="right"
      open={visible}
      sx={{
        width: drawerWidth,
        "& .MuiDrawer-paper": {
          width: drawerWidth,
          boxSizing: "border-box",
        },
      }}
    >
      <Toolbar />
      <Stack direction="column" justifyContent="space-between" flexGrow={1}>
        <Box component="div">
          <Typography
            variant="h5"
            component="div"
            textAlign="center"
            color="primary.main"
            sx={{ mb: 1, mt: 1 }}
          >
            Device Advanced Setting
          </Typography>
          <Divider />
          <Paper
            elevation={0}
            sx={{ pl: 3, pr: 3, mt: 2, backgroundColor: "inherit" }}
          >
            <Typography
              variant="subtitle1"
              textAlign="center"
            >{`EHG7608-4PoE-4SFP (00:60:E9:21:2C:BE)`}</Typography>
            <TabContext value={value}>
              <Box sx={{ borderBottom: 1, borderColor: "divider" }}>
                <TabList onChange={handleChange}>
                  <Tab label="Authentication" value="1" />
                  <Tab label="Alarm" value="2" />
                </TabList>
              </Box>
              <TabPanel value="1" sx={{ pl: 0, pr: 0, pt: 0 }}>
                <Divider
                  textAlign="left"
                  sx={{ color: "info.main", mb: 1, mt: 2 }}
                >
                  General
                </Divider>
                <Controls.Input
                  name="username"
                  label="Username"
                  value="admin"
                  fullWidth={true}
                />
                <Controls.Input
                  name="password"
                  label="Password"
                  type="password"
                  value="default"
                  fullWidth={true}
                />
                <Divider
                  textAlign="left"
                  sx={{ color: "info.main", mb: 1, mt: 2 }}
                >
                  SNMP
                </Divider>
                <Controls.Select
                  name="snmpVersion"
                  label="Snmp Version"
                  value="v2c"
                  options={[
                    { id: "v1", title: "V1" },
                    { id: "v2c", title: "V2C" },
                  ]}
                  fullWidth={true}
                />
                <Controls.Input
                  name="readComunity"
                  label="Read Community"
                  value="public"
                  fullWidth={true}
                />
                <Controls.Input
                  name="writeCommunity"
                  label="Write Community"
                  value="private"
                  fullWidth={true}
                />
              </TabPanel>
              <TabPanel value="2" sx={{ pl: 0, pr: 0, pt: 0 }}>
                <Divider
                  textAlign="left"
                  sx={{ color: "info.main", mb: 1, mt: 2 }}
                >
                  Port
                </Divider>
                <AlarmPortList items={portData} />
                <Divider
                  textAlign="left"
                  sx={{ color: "info.main", mb: 1, mt: 2 }}
                >
                  Power
                </Divider>
                <AlarmPowerList items={powerData} />
              </TabPanel>
            </TabContext>
          </Paper>
        </Box>
        <Toolbar>
          <Stack direction="row" justifyContent="end" flexGrow={1} spacing={2}>
            <Button
              color="warning"
              variant="outlined"
              onClick={handleCancleClick}
            >
              Cancel
            </Button>
            <Button
              color="primary"
              variant="contained"
              onClick={handleSaveClick}
            >
              Save
            </Button>
          </Stack>
        </Toolbar>
      </Stack>
    </Drawer>
  );
};

export default AdvancedSettingDrawer;
