import React from "react";
import { useDispatch } from "react-redux";
import DeviceMUITable from "../../components/common/DeviceEnhancedTable/DeviceMUITable";
import { openAdvancedSettingDrawer } from "../../features/DeviceAdvancedSetting";
import { openNetworkSettingDrawer } from "../../features/SingleNetworkSetting";
import devicedata from "../../utils/data/deviceData.json";

const ManagementDevice = () => {
  const dispatch = useDispatch();
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

  const handleContextMenuClicked = (event, data, name) => {
    switch (name) {
      case "openWeb":
        window.open(`http://10.6.50.7`, "_blank");
        break;
      case "networkSetting":
        dispatch(openNetworkSettingDrawer());
        break;
      case "deviceAdvancedSetting":
        dispatch(openAdvancedSettingDrawer());
        break;
      default:
        break;
    }
  };

  return (
    <DeviceMUITable
      Columns={column}
      DataSource={devicedata}
      rowKey="id"
      title="Device management"
      TotalRowCount={devicedata.length}
      currentPage={1}
      handleOnSelect={(value) => console.log(value)}
      options={{ sortable: true, selectable: true, contextMenu: true }}
      onContextMenuClick={(event, data, name) => {
        handleContextMenuClicked(event, data, name);
      }}
      toolbarOptions={{
        globalFilter: true,
      }}
    />
  );
};

export default ManagementDevice;
