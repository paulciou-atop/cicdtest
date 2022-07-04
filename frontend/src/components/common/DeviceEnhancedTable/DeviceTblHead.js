import React from "react";
import TableHead from "@mui/material/TableHead";
import TableRow from "@mui/material/TableRow";
import TableCell from "@mui/material/TableCell";
import TableSortLabel from "@mui/material/TableSortLabel";
import { Checkbox } from "@mui/material";

import { makeStyles } from "@mui/styles";

const useStyles = makeStyles((theme) => ({
  tableHeaderCell: {
    fontWeight: "bold",
    paddingTop: "8px",
    paddingBottom: "8px",
    backgroundColor: theme.palette.primary.dark,
    color: theme.palette.getContrastText(theme.palette.primary.dark),
    "& .MuiCheckbox-root": {
      color: theme.palette.getContrastText(theme.palette.primary.dark),
    },
  },
}));

const TblHead = ({
  Columns,
  order,
  orderBy,
  handleSortRequest,
  sortable,
  isActions,
  numSelected,
  onSelectAllClick,
  rowCount,
  selectable,
}) => {
  const classes = useStyles();
  return (
    <TableHead>
      <TableRow>
        {selectable && (
          <TableCell padding="checkbox" className={classes.tableHeaderCell}>
            <Checkbox
              indeterminate={numSelected > 0 && numSelected < rowCount}
              checked={rowCount > 0 && numSelected === rowCount}
              onChange={onSelectAllClick}
            />
          </TableCell>
        )}
        {Columns.map((headCell) => (
          <TableCell
            className={classes.tableHeaderCell}
            key={headCell.key}
            sortDirection={orderBy === headCell.id ? order : false}
          >
            {headCell.disableSort !== true && sortable ? (
              <TableSortLabel
                active={orderBy === headCell.key}
                onClick={() => handleSortRequest(headCell.key)}
                direction={orderBy === headCell.key ? order : "asc"}
              >
                {headCell.label}
              </TableSortLabel>
            ) : (
              headCell.label
            )}
          </TableCell>
        ))}
        {isActions && (
          <TableCell className={classes.tableHeaderCell}>Actions</TableCell>
        )}
      </TableRow>
    </TableHead>
  );
};

export default TblHead;
