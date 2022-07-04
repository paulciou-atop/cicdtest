import Table from "@mui/material/Table";
import TableContainer from "@mui/material/TableContainer";
import Paper from "@mui/material/Paper";
import React, { useState } from "react";
import TblHead from "./DeviceTblHead";
import TblBody from "./DeviceTblBody";
import TblPagination from "./DeviceTblPagination";
import TblToolbar from "./DeviceTblToolbar";

const DeviceMUITable = ({
  Columns,
  rowKey,
  title,
  DataSource,
  TotalRowCount,
  currentPage = 0,
  options: { sortable, selectable, contextMenu },
  toolbarOptions: { toolbarAction = [], globalFilter },
  handleOnSelect,
  actions = [],
  size = "medium",
  onContextMenuClick,
}) => {
  const [page, setPage] = useState(currentPage - 1);
  const [rowsPerPage, setRowsPerPage] = useState(10);
  const [order, setorder] = useState();
  const [orderBy, setorderBy] = useState();
  const [selected, setSelected] = useState([]);
  const [inputSearch, setInputSearch] = useState("");

  const handleSetSelected = (value) => {
    setSelected(value);
    handleOnSelect(value);
  };
  const handleSelectAllClick = (event) => {
    if (event.target.checked) {
      const newSelecteds = DataSource.map((n) => n[rowKey]);
      setSelected(newSelecteds);
      handleOnSelect(newSelecteds);
      return;
    }
    setSelected([]);
    handleOnSelect([]);
  };
  // Avoid a layout jump when reaching the last page with empty rows.
  const emptyRows =
    page > 0 ? Math.max(0, (1 + page) * rowsPerPage - DataSource.length) : 0;

  const isActions = actions.length > 0;

  const handleChangePage = (event, newPage) => {
    setPage(newPage);
  };

  const handleChangeRowsPerPage = (event) => {
    setRowsPerPage(parseInt(event.target.value, 10));
    setPage(0);
  };

  const handleSortRequest = (cellId) => {
    const isAsc = orderBy === cellId && order === "asc";
    setorder(isAsc ? "desc" : "asc");
    setorderBy(cellId);
  };

  const recordAfterfiltering = (dataSource) => {
    return dataSource.filter((row) => {
      let rec = Columns.map((element) =>
        row[element.key].includes(inputSearch)
      );
      return rec.includes(true);
    });
  };

  const recordAfterPaginationAndSorting = (dataSource) => {
    return stableSort(
      recordAfterfiltering(dataSource),
      getComparator(order, orderBy)
    ).slice(page * rowsPerPage, page * rowsPerPage + rowsPerPage);
  };

  return (
    <React.Fragment>
      <TableContainer component={Paper}>
        <TblToolbar
          selected={selected}
          title={title}
          globalFilter={globalFilter}
          inputSearch={inputSearch}
          setInputSearch={setInputSearch}
          toolbarAction={toolbarAction}
        />
        <Table size={size}>
          <TblHead
            Columns={Columns}
            order={order}
            orderBy={orderBy}
            handleSortRequest={handleSortRequest}
            sortable={sortable}
            isActions={isActions}
            rowCount={TotalRowCount}
            numSelected={selected.length}
            onSelectAllClick={handleSelectAllClick}
            selectable={selectable}
          />
          <TblBody
            Columns={Columns}
            Items={recordAfterPaginationAndSorting(DataSource)}
            emptyRows={emptyRows}
            actions={actions}
            rowKey={rowKey}
            selected={selected}
            handleSetSelected={handleSetSelected}
            selectable={selectable}
            contextMenus={contextMenu}
            onContextMenuClick={onContextMenuClick}
          />
          <TblPagination
            rowsPerPageOptions={[5, 10, 15]}
            count={TotalRowCount ? TotalRowCount : DataSource.length}
            rowsPerPage={rowsPerPage}
            page={page}
            onPageChange={handleChangePage}
            onRowsPerPageChange={handleChangeRowsPerPage}
          />
        </Table>
      </TableContainer>
    </React.Fragment>
  );
};

export default DeviceMUITable;

function getComparator(order, orderBy) {
  return order === "desc"
    ? (a, b) => descendingComparator(a, b, orderBy)
    : (a, b) => -descendingComparator(a, b, orderBy);
}

function stableSort(array, comparator) {
  const stabilizedThis = array.map((el, index) => [el, index]);
  stabilizedThis.sort((a, b) => {
    const order = comparator(a[0], b[0]);
    if (order !== 0) return order;
    return a[1] - b[1];
  });
  return stabilizedThis.map((el) => el[0]);
}

function descendingComparator(a, b, orderBy) {
  if (b[orderBy] < a[orderBy]) {
    return -1;
  }
  if (b[orderBy] > a[orderBy]) {
    return 1;
  }
  return 0;
}
