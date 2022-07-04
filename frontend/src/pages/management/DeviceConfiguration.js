import React, { useState } from "react";
import { Button, Grid, Paper } from "@mui/material";
import {
  Publish,
  Event,
  AcUnit,
  DeviceHub,
  Backup,
  SettingsBackupRestore,
  Schedule,
} from "@mui/icons-material";
import devicedata from "../../utils/data/deviceDataConfig.json";
import DeviceMUITable from "../../components/common/DeviceEnhancedTable/DeviceMUITable";
import EventTips from "../../components/eventTips/EventTips";

const DeviceConfiguration = () => {
  const [selectedDevice, setSelectedDevice] = useState([]);
  const [selectedValue, setSelectedValue] = useState([]);
  const [selectableDevice, setSelectableDevice] = useState(false);
  const [configButtonClicked, setconfigButtonClicked] = useState("");
  const [ClickedValue, setClickedValue] = useState("");
  const [openPopup, setopenPopup] = useState(false);
  const column = [
    {
      key: "deviceType",
      label: "Device Type",
      disableSort: true,
    },
    {
      key: "model",
      label: "Model",
      disableSort: true,
    },
    {
      key: "ipAddress",
      label: "IP Address",
      disableSort: true,
    },
    {
      key: "macAddress",
      label: "MAC Address",
    },
    {
      key: "hostName",
      label: "Hostname",
    },
    {
      key: "kernel",
      label: "Kernel",
    },
    { key: "ap", label: "Firmware Ver" },
  ];

  const handleDeviceSelect = (value) => {
    setSelectedDevice(value);
    setSelectedValue(value)
  };

  const handleConfigButtonClicked = (value) => {
    setconfigButtonClicked(value);
    setClickedValue(value);
    setSelectableDevice(true);
  };

  const prepareData = () => {
    switch (configButtonClicked) {
      case "syslogSetting":
      case "trapSetting":
      case "backupAndRestore":
      case "resetToDefault":
        return devicedata.filter(
          (device) =>
            device.deviceType === "GWD/SNMP" || device.deviceType === "SNMP"
        );
      case "enableSnmp":
        return devicedata.filter((device) => device.deviceType === "GWD");
      case "":
        return [];

      default:
        return devicedata;
    }
  };

  var handleOkClick;
  if (selectedDevice.length===0){
  handleOkClick = () => {
    setconfigButtonClicked("");
    setSelectableDevice(false);
    setopenPopup(false)
        setSelectedDevice([])
  }
  }
  else{
    handleOkClick = () => {
    setconfigButtonClicked("");
    setSelectableDevice(false);
    setopenPopup(true)
    setSelectedDevice([])
  }
  }

  const handleCancelClick = () => {
    setconfigButtonClicked("");
    setSelectableDevice(false);
  };

  return (
    <Paper sx={{ p: 2 }}>
      <EventTips
        selectedValue={selectedValue}
        selectedDevice={selectedDevice}
        configButtonClicked={configButtonClicked}
        handleCancelClick={handleCancelClick}
        handleOkClick={handleOkClick}
        openPopup={openPopup}
        setopenPopup={setopenPopup}
        ClickedValue={ClickedValue}
      />
      <Grid container spacing={2}>
        <Grid item xs={6} md={3}>
          <Button
            variant="contained"
            color="primary"
            size="large"
            fullWidth
            startIcon={<Publish />}
            onClick={() => handleConfigButtonClicked("fwUpload")}
          >
            Firmware Upload
          </Button>
        </Grid>
        <Grid item xs={6} md={3}>
          <Button
            variant="contained"
            color="primary"
            size="large"
            fullWidth
            startIcon={<Event />}
            onClick={() => handleConfigButtonClicked("syslogSetting")}
          >
            Syslog Server Settings
          </Button>
        </Grid>
        <Grid item xs={6} md={3}>
          <Button
            variant="contained"
            color="primary"
            size="large"
            fullWidth
            startIcon={<AcUnit />}
            onClick={() => handleConfigButtonClicked("trapSetting")}
          >
            Trap Server Settings
          </Button>
        </Grid>

        <Grid item xs={6} md={3}>
          <Button
            variant="contained"
            color="primary"
            size="large"
            fullWidth
            startIcon={<DeviceHub />}
            onClick={() => handleConfigButtonClicked("networkSetting")}
          >
            Network Settings
          </Button>
        </Grid>
        <Grid item xs={6} md={3}>
          <Button
            variant="contained"
            color="primary"
            size="large"
            fullWidth
            startIcon={<Backup />}
            onClick={() => handleConfigButtonClicked("backupAndRestore")}
          >
            Backup and Restore
          </Button>
        </Grid>
        <Grid item xs={6} md={3}>
          <Button
            variant="contained"
            color="primary"
            size="large"
            fullWidth
            startIcon={<SettingsBackupRestore />}
            onClick={() => handleConfigButtonClicked("resetToDefault")}
          >
            Reset to default
          </Button>
        </Grid>
        <Grid item xs={6} md={3}>
          <Button
            variant="contained"
            color="primary"
            size="large"
            fullWidth
            startIcon={<Schedule />}
            onClick={() => handleConfigButtonClicked("")}
          >
            schedule backup
          </Button>
        </Grid>
        <Grid item xs={6} md={3}>
          <Button
            variant="contained"
            color="primary"
            size="large"
            fullWidth
            startIcon={<Schedule />}
            onClick={() => handleConfigButtonClicked("enableSnmp")}
          >
            Enable SNMP
          </Button>
        </Grid>
      </Grid>
      <Grid container>
        <Grid item xs={12} md={12} sx={{ pt: 2 }}>
          {prepareData().length > 0 && (
            <DeviceMUITable
              Columns={column}
              DataSource={prepareData()}
              rowKey="id"
              title="Select Device to Configure"
              TotalRowCount={prepareData().length}
              currentPage={1}
              handleOnSelect={(value) => {
                handleDeviceSelect(value);
              }}
              options={{
                sortable: true,
                selectable: selectableDevice,
                contextMenu: false,
              }}
              onContextMenuClick={(event, data, name) => {
                console.log("context menu clicked");
              }}
              toolbarOptions={{
                globalFilter: true,
              }}
            />
          )}
        </Grid>
      </Grid>
    </Paper>
  );
};

export default DeviceConfiguration;
