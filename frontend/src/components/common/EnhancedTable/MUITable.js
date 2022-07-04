import Table from "@mui/material/Table";
import TableContainer from "@mui/material/TableContainer";
import Paper from "@mui/material/Paper";
import React, { useState } from "react";
import TblHead from "./TblHead";
import TblBody from "./TblBody";
import TblPagination from "./TblPagination";
import TblToolbar from "./TblToolbar";

const MUITable = ({
  Columns,
  selectedCol,
  rowKey,
  title,
  DataSource,
  TotalRowCount,
  currentPage = 0,
  options: { sortable, selectable, filtrable },
  toolbarOptions: { exportData, createNew, deleteSelected, globalFilter, columnSelected },
  actions = [],
  size = "small",
}) => {
  const [page, setPage] = useState(currentPage - 1);
  const [rowsPerPage, setRowsPerPage] = useState(10);
  const [order, setorder] = useState();
  const [orderBy, setorderBy] = useState();
  const [selected, setSelected] = useState([]);
  const [inputSearch, setInputSearch] = useState("");

  const handleSetSelected = (value) => {
    setSelected(value);
  };
  const handleSelectAllClick = (event) => {
    if (event.target.checked) {
      const newSelecteds = DataSource.map((n) => n[rowKey]);
      setSelected(newSelecteds);
      return;
    }
    setSelected([]);
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

  //console.log(recordAfterfiltering(DataSource));
  return (
    <React.Fragment>
      <TableContainer component={Paper}>
        <TblToolbar
          selected={selected}
          createNew={createNew}
          deleteSelected={deleteSelected}
          exportData={exportData}
          title={title}
          selectable={selectable}
          Columns={Columns}
          DataSource={DataSource}
          globalFilter={globalFilter}
          inputSearch={inputSearch}
          setInputSearch={setInputSearch}
          columnSelected={columnSelected}
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
            selectedCol={selectedCol}
            columnSelected={columnSelected}
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
            selectedCol={selectedCol}
            columnSelected={columnSelected}
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

export default MUITable;

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
