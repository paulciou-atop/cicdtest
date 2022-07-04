import { Alert, AlertTitle, Button,  CardHeader, Box, Stack } from "@mui/material";
import React, { useState } from "react";
import Popup from '../common/PopupDevice';
import SystemLogForm from "../syslogManagement/SystemLogForm";
import SystemLogTable from "../syslogManagement/SystemLogTable";
import TrapServerForm from "../trapManagement/TrapServerForm";
import TrapServerTable from "../trapManagement/TrapServerTable";

const TIPS = "(This feature only for device with SNMP support.)";
const messages = {
  fwUpload: "Firmware Upload",
  networkSetting: "Network Setting",
  resetToDefault: "Reset To Default",
  backupAndRestore: "Backup and Restore",
  syslogSetting: "Syslog Server Setting",
  trapSetting: "Trap Server Setting",
};
const popupTitle = {
  fwUpload: "Firmware Configuration Field",
  networkSetting: "Network Configuration Field",
  resetToDefault: "Reset Configuration Field",
  backupAndRestore: "Backup Configuration Field",
  syslogSetting: "Syslog Server Configuration Field",
  trapSetting: "Trap Configuration Field",
};

const EventTips = ({
  selectedDevice,
  configButtonClicked,
  handleOkClick,
  handleCancelClick,
  openPopup,
  setopenPopup,
  ClickedValue,
  selectedValue,
}) => {

  const [deviceStatus, setdeviceStatus] = useState("")
  const [deviceValue, setdeviceValue] = useState("")
  const AddorEdit = (device, resetForm) => {
    if (device.id === "0") console.log("Insert Record", device);
    else console.log("update Record", device);
    resetForm();
    setdeviceStatus(device)
    setdeviceValue(device)
  };

  var statusValue;
  if(deviceStatus==="" && openPopup === true){
    statusValue="WAITING"
  }
  else if(deviceStatus!=="" && openPopup === false){
    statusValue="WAITING"
    setdeviceStatus("")
  }
  else{
    statusValue="SUCCESS"
  }
  
 return (
    <Alert
      severity="info"
      className={configButtonClicked === "" ? "alert hide" : "alert"}
    >
      <AlertTitle>{messages[configButtonClicked]}</AlertTitle>
      <div>
        <div>
          Select devices and press
          <Button
            variant="text"
            size="medium"
            sx={{ minWidth: 5 }}
            disabled={selectedDevice.length === 0}
            onClick={handleOkClick}
          >
            ok
          </Button>
          <Popup
            title={messages[ClickedValue]}
            openPopup={openPopup}
            setOpenPopup={setopenPopup}
            >
          <Box sx={{ width: '100%' }}>
          <Stack spacing={3}>
          <CardHeader 
          subheader={popupTitle[ClickedValue]}
          />
          {popupTitle[ClickedValue] === 'Syslog Server Configuration Field'?  
          <Stack spacing={3}>
          <SystemLogForm AddorEdit={AddorEdit}  />
          <SystemLogTable  selectedDevice={selectedValue} statusValue={statusValue}/>
          </Stack>: null }

          {popupTitle[ClickedValue] === 'Trap Configuration Field'?  
          <Stack spacing={3}>
          <TrapServerForm AddorEdit={AddorEdit}  />
          <TrapServerTable selectedDevice={selectedValue} statusValue={statusValue}/>
          </Stack>: null }
           
          </Stack>
          </Box>
          </Popup>
            or
          <Button variant="text" size="medium" onClick={handleCancelClick}>
            cancel
          </Button>
          .
        </div>
        <div style={{ color: "red" }}>{TIPS}</div>
      </div>
    </Alert>
  );
};

export default EventTips;
