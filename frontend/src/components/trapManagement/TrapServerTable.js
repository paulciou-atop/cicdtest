import { Add, Delete } from "@mui/icons-material";
import { Paper } from "@mui/material";
import React from "react";
import MUITable from "../common/EnhancedTable/MUITable";
import devicedata from "../../utils/data/deviceDataConfig.json";

const TrapServerTable = ({
  selectedDevice, statusValue}) => {

let resultArr = selectedDevice.map(i => devicedata[i])
resultArr.forEach(object => {
  object.status = statusValue;
});

const column = [
    { key: "model", label: "Model" },
    { key: "macAddress", label: "MAC Address" },
    { key: "ipAddress", label: "IP Address" },
    { key: "status", label: "Status" },
  ];

  return (
    <Paper>
        <MUITable
        Columns={column}
        rowKey="id"
        title="Devices List"
        DataSource={resultArr}
        TotalRowCount={resultArr.length}
        currentPage={1}
        options={{ sortable: false, selectable: false, filtrable: false }}
        toolbarOptions={{
          exportData: false,
          createNew: {
            enable: false,
            icon: <Add />,
            onClick: (event) => {
            },
          },
          deleteSelected: {
            enable: false,
            label: "Delete",
            icon: <Delete />,
            onClick: (event, selectedRow) =>
              console.log("Add button clicked", selectedRow),
          },
          globalFilter: false,
        }}
      />

    </Paper>
  );
};

export default TrapServerTable;
