const DeviceListColumn = () => {
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
      key: "description",
      label: "Description",
    },
    {
      key: "kernel",
      label: "Kernel",
    },
    { key: "firmwareVer", label: "Firmware Ver" },
  ];
  return column;
};

export default DeviceListColumn;
