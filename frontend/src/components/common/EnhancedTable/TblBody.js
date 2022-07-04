import React from "react";
import TableBody from "@mui/material/TableBody";
import TableRow from "@mui/material/TableRow";
import TableCell from "@mui/material/TableCell";
import IconButton from "@mui/material/IconButton";
import { Checkbox } from "@mui/material";

const TblBody = ({
  Columns,
  Items,
  emptyRows,
  actions,
  rowKey,
  selected,
  handleSetSelected,
  selectable,
  selectedCol,
  columnSelected,
}) => {
  const isSelected = (id) => selected.indexOf(id) !== -1;

  const handleClick = (event, id) => {
    const selectedIndex = selected.indexOf(id);
    let newSelected = [];

    if (selectedIndex === -1) {
      newSelected = newSelected.concat(selected, id);
    } else if (selectedIndex === 0) {
      newSelected = newSelected.concat(selected.slice(1));
    } else if (selectedIndex === selected.length - 1) {
      newSelected = newSelected.concat(selected.slice(0, -1));
    } else if (selectedIndex > 0) {
      newSelected = newSelected.concat(
        selected.slice(0, selectedIndex),
        selected.slice(selectedIndex + 1)
      );
    }

    handleSetSelected(newSelected);
  };

  return (
    <TableBody>
      {Items.length <= 0 && (
        <TableRow
          style={{
            height: 100,
          }}
        >
          <TableCell
            colSpan={
              actions.length > 0
                ? selectable
                  ? Columns.length + 2
                  : Columns.length + 1
                : Columns.length
            }
            align="center"
          >
            No data found
          </TableCell>
        </TableRow>
      )}
      {Items.map((item) => {
        const isItemSelected = isSelected(item[rowKey]);
        return (
          <TableRow key={item[rowKey]} hover selected={isItemSelected}>
            {selectable && (
              <TableCell padding="checkbox">
                <Checkbox
                  color="primary"
                  checked={isItemSelected}
                  onChange={(event) => handleClick(event, item[rowKey])}
                />
              </TableCell>
            )}
            {Columns.map((bodyCell) => (       
               columnSelected && columnSelected.enable === true ?(  selectedCol.includes(bodyCell.key) === true ? (
                bodyCell.render !== undefined ? (
                <TableCell key={bodyCell.key}>
                  {bodyCell.render(item)}
                </TableCell>
              ) : (
                <TableCell key={bodyCell.key}>{item[bodyCell.key]}</TableCell>
              )
             ): null 
               ) :(
                bodyCell.render !== undefined ? (
                  <TableCell key={bodyCell.key}>
                    {bodyCell.render(item)}
                  </TableCell>
                ) : (
                  <TableCell key={bodyCell.key}>{item[bodyCell.key]}</TableCell>
                )
               )
      ))}
            {actions.length > 0 && (
              <TableCell>
                {actions.map((action) => (
                  <IconButton
                    key={action.key}
                    color={action.color || "primary"}
                    onClick={(e) => action.onClick(e, item)}
                  >
                    {action.icon}
                  </IconButton>
                ))}
              </TableCell>
            )}
          </TableRow>
        );
      })}
      {emptyRows > 0 && (
        <TableRow
          style={{
            height: actions.length > 0 ? 73 * emptyRows : 53 * emptyRows,
          }}
        >
          <TableCell
            colSpan={
              actions.length > 0
                ? selectable
                  ? Columns.length + 2
                  : Columns.length + 1
                : Columns.length
            }
          />
        </TableRow>
      )}
    </TableBody>
  );
};

export default TblBody;

