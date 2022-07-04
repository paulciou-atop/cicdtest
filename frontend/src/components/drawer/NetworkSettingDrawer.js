import React from "react";
import Drawer from "@mui/material/Drawer";
import {
  Alert,
  Button,
  Checkbox,
  Divider,
  FormControlLabel,
  Paper,
  Stack,
  TextField,
  Toolbar,
  Typography,
} from "@mui/material";
import { useDispatch, useSelector } from "react-redux";
import {
  closeNetworkSettingDrawer,
  singleNetworkSettingSelector,
} from "../../features/SingleNetworkSetting";
import { Box } from "@mui/system";

const drawerWidth = 400;

const inputItems = [
  { id: "IPAddress", label: "IP Address", valid: "10.6.50.7" },
  { id: "netmask", label: "Netmask ", valid: "255.255.255.0" },
  { id: "gateway", label: "Gateway", valid: "10.6.50.1" },
];

const SNMPonlyInputItem = [
  { id: "dns1", label: "Preferred DNS server", valid: "0.0.0.0" },
  { id: "dns2", label: "Alternate DNS server", valid: "0.0.0.0" },
];
const networkSettingTips =
  "Please make sure device username password setting and SNMP community is correct.";

const NetworkSettingDrawer = () => {
  const dispatch = useDispatch();
  const { visible } = useSelector(singleNetworkSettingSelector);

  const handleSaveClick = () => {
    dispatch(closeNetworkSettingDrawer());
  };
  const handleCancleClick = () => {
    dispatch(closeNetworkSettingDrawer());
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
            Network Setting
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
            <FormControlLabel control={<Checkbox />} label="DHCP" />
            {inputItems.map((item) => (
              <TextField
                key={item.id}
                label={item.label}
                fullWidth
                variant="outlined"
                margin="dense"
                size="small"
                value={item.valid}
              />
            ))}
            {SNMPonlyInputItem.map((item) => (
              <TextField
                key={item.id}
                label={item.label}
                fullWidth
                variant="outlined"
                margin="dense"
                size="small"
                value={item.valid}
              />
            ))}
            <TextField
              margin="dense"
              label="Hostname"
              fullWidth
              value="Switch 1"
              variant="outlined"
              size="small"
            />
            <Alert severity="error" color="warning">
              {networkSettingTips}
            </Alert>
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

export default NetworkSettingDrawer;
