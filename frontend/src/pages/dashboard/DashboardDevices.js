import { Add, Delete } from "@mui/icons-material";
import { Stack } from "@mui/material";
import React, { useEffect, useState } from "react";
import { useDispatch, useSelector } from "react-redux";
import MUITable from "../../components/common/EnhancedTable/MUITable";
import column from "../../components/deviceList/DeviceListColumn";
import SessionDetails from "../../components/deviceList/SessionDetails";
import {
  deviceDiscoverySelector,
  GetCheckStatus,
  GetlastScannedDeviceData,
  GetSessionID,
} from "../../features/DeviceDiscoverySlice";
//import deviceDatas from "../../utils/data/deviceData.json";

const DashboardDevices = () => {
  const dispatch = useDispatch();
  const { sessionDetails, deviceData, titleSessionId } = useSelector(
    deviceDiscoverySelector
  );
  const [ipRange, setIpRange] = useState("10.6.50.1/24");
  useEffect(() => {
    dispatch(GetlastScannedDeviceData());
  }, []); // eslint-disable-line react-hooks/exhaustive-deps

  const handleStartClick = () => {
    dispatch(GetSessionID(ipRange));
    let statuscheck = "";
    setInterval(() => {
      if (statuscheck === "" || statuscheck === "running") {
        dispatch(GetCheckStatus())
          .unwrap()
          .then((result) => {
            statuscheck = result.info.status;
          })
          .catch((error) => {
            statuscheck = error;
          });
      }
    }, 5000);
  };

  return (
    <Stack spacing={2}>
      <SessionDetails
        ipRange={ipRange}
        setIpRange={setIpRange}
        sessionDetails={sessionDetails}
        onStartClick={handleStartClick}
      />
      {deviceData.length >= 0 && (
        <MUITable
          Columns={column()}
          rowKey="macAddress"
          title={`Device List (session: ${titleSessionId})`}
          DataSource={deviceData}
          TotalRowCount={deviceData.length}
          currentPage={1}
          options={{ sortable: true, selectable: false, filtrable: true }}
          toolbarOptions={{
            exportData: true,
            createNew: {
              enable: false,
              label: "Add New",
              icon: <Add />,
              onClick: (event) => console.log("Add button clicked"),
            },
            deleteSelected: {
              enable: false,
              label: "Delete",
              icon: <Delete />,
              onClick: (event, selectedRow) =>
                console.log("Delete button clicked", selectedRow),
            },
            globalFilter: true,
          }}
        />
      )}
    </Stack>
  );
};

export default DashboardDevices;
