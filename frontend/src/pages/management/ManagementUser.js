import { Add, Delete, Edit } from "@mui/icons-material";
import { Paper } from "@mui/material";
import React, { useState } from "react";
import MUITable from "../../components/common/EnhancedTable/MUITable";
import Popup from "../../components/common/Popup";
import UserForm from "../../components/userManagement/UserForm";
import usersData from "../../utils/data/usersData.json";

const ManagementUser = () => {
  const [openPopup, setopenPopup] = useState(false);
  const [recordsForEdit, setrecordsForEdit] = useState(null);

  const AddorEdit = (user, resetForm) => {
    if (user.id === "0") console.log("Insert Record", user);
    else console.log("update Record", user);
    resetForm();
    setrecordsForEdit(null);
    setopenPopup(false);
  };

  const openInPopup = (item) => {
    setrecordsForEdit(item);
    setopenPopup(true);
  };

  const column = [
    { key: "name", label: "Name" },
    { key: "username", label: "Username" },
    { key: "email", label: "Email" },
    { key: "role", label: "Role" },
    { key: "createdBy", label: "Created By" },
    { key: "createdAt", label: "Created At", disableSort: true },
  ];
  const actions = [
    {
      key: "edit",
      icon: <Edit />,
      color: "primary",
      tooltip: "Edit User",
      onClick: (event, rowData) => openInPopup(rowData),
    },
  ];

  return (
    <Paper>
      <MUITable
        Columns={column}
        rowKey="id"
        title="User Management"
        DataSource={usersData}
        TotalRowCount={usersData.length}
        currentPage={1}
        options={{ sortable: true, selectable: true, filtrable: true }}
        toolbarOptions={{
          exportData: true,
          createNew: {
            enable: true,
            label: "Add New",
            icon: <Add />,
            onClick: (event) => {
              setopenPopup(true);
              setrecordsForEdit(null);
            },
          },
          deleteSelected: {
            enable: true,
            label: "Delete",
            icon: <Delete />,
            onClick: (event, selectedRow) =>
              console.log("Add button clicked", selectedRow),
          },
          globalFilter: true,
        }}
        actions={actions}
      />
      <Popup
        title="User Form"
        openPopup={openPopup}
        setOpenPopup={setopenPopup}
      >
        <UserForm AddorEdit={AddorEdit} recordsForEdit={recordsForEdit} />
      </Popup>
    </Paper>
  );
};

export default ManagementUser;
