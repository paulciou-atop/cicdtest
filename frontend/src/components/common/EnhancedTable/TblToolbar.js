import { Close, Download, Search } from "@mui/icons-material";
import {
  alpha,
  Button,
  IconButton,
  InputAdornment,
  Menu,
  MenuItem,
  Stack,
  TextField,
  Toolbar,
  Typography,
} from "@mui/material";
import jsPDF from "jspdf";
import "jspdf-autotable";
import React, { useState } from "react";
import { CSVLink } from "react-csv";

const TblToolbar = ({
  selected,
  createNew,
  deleteSelected,
  exportData,
  title,
  Columns,
  DataSource,
  globalFilter,
  inputSearch,
  setInputSearch,
  columnSelected,
}) => {
  const [anchorEl, setAnchorEl] = useState(null);
  const [csvHeader, setCsvHeader] = useState([]);
  const [csvData, setCsvData] = useState([]);
  //const [inputSearch, setInputSearch] = useState("");

  const open = Boolean(anchorEl);
  const handleClick = (event) => {
    setAnchorEl(event.currentTarget);
  };

  const handleCsvDownload = (event, done) => {
    const headers = Columns.map((col) => col.label);
    const data = DataSource.map((item) => {
      return Columns.map((col) => item[col.key]);
    });
    setCsvHeader(headers);
    setCsvData(data);
    //setallowDownload(true);
    handleClose();
    done(true);
  };

  const handleMouseDownPassword = (event) => {
    event.preventDefault();
  };

  const handlePdfDownload = (event) => {
    const headers = [Columns.map((col) => col.label)];
    const data = DataSource.map((item) => {
      return Columns.map((col) => item[col.key]);
    });
    exportPDF(data, headers, title);
    handleClose();
  };

  const handleClose = () => {
    setAnchorEl(null);
  };
  return (
    <Toolbar
      sx={{
        pl: { sm: 2 },
        pr: { xs: 1, sm: 1 },
        ...(selected.length > 0 && {
          bgcolor: (theme) =>
            alpha(
              theme.palette.secondary.main,
              theme.palette.action.activatedOpacity
            ),
        }),
      }}
    >
      {selected.length > 0 ? (
        <Typography
          sx={{ flexGrow: "1" }}
          color="inherit"
          variant="subtitle1"
          component="div"
        >
          {selected.length} selected
        </Typography>
      ) : (
        <Typography
          sx={{ flexGrow: "1" }}
          variant="h6"
          id="tableTitle"
          component="div"
        >
          {title}
        </Typography>
      )}
      <Stack direction="row" spacing={2}>
        {selected.length <= 0 && globalFilter && (
          <TextField
            label="Search"
            value={inputSearch}
            onChange={(e) => setInputSearch(e.target.value)}
            size="small"
            InputProps={{
              startAdornment: (
                <InputAdornment position="start">
                  <Search />
                </InputAdornment>
              ),
              endAdornment: (
                <InputAdornment position="end">
                  <IconButton
                    onClick={() => setInputSearch("")}
                    onMouseDown={handleMouseDownPassword}
                    edge="end"
                  >
                    <Close />
                  </IconButton>
                </InputAdornment>
              ),
            }}
          />
        )}
        {selected.length <=0 && columnSelected && columnSelected.enable &&(
          <Button
          variant="outlined"
          startIcon={columnSelected.icon}
          onClick={(e) => columnSelected.onClick(e)}
        >
          {columnSelected.label}
        </Button>
        )
        }
        {selected.length <= 0 && createNew && createNew.enable && (
          <Button
            variant="outlined"
            startIcon={createNew.icon}
            onClick={(e) => createNew.onClick(e)}
          >
            {createNew.label}
          </Button>
        )}
        {selected.length > 0 && deleteSelected && deleteSelected.enable && (
          <Button
            variant="outlined"
            color="secondary"
            startIcon={deleteSelected.icon}
            onClick={(e) => deleteSelected.onClick(e, selected)}
          >
            {deleteSelected.label}
          </Button>
        )}
        {selected.length <= 0 && exportData && (
          <>
            <Button
              variant="outlined"
              startIcon={<Download />}
              onClick={handleClick}
              sx={{ mr: 2 }}
            >
              export
            </Button>
            <Menu
              anchorEl={anchorEl}
              open={open}
              onClose={handleClose}
              MenuListProps={{
                "aria-labelledby": "basic-button",
              }}
              transformOrigin={{ horizontal: "right", vertical: "top" }}
              anchorOrigin={{ horizontal: "right", vertical: "bottom" }}
            >
              <MenuItem onClick={handlePdfDownload}>PDF Export</MenuItem>
              <MenuItem>
                <CSVLink
                  headers={csvHeader}
                  data={csvData}
                  filename={getFilename(title)}
                  asyncOnClick={true}
                  onClick={handleCsvDownload}
                  style={{ textDecoration: "none", color: "inherit" }}
                >
                  CSV Export
                </CSVLink>
              </MenuItem>
            </Menu>
          </>
        )}
      </Stack>
    </Toolbar>
  );
};

export default TblToolbar;

const exportPDF = (data, headers, title) => {
  const unit = "pt";
  const size = "A4"; // Use A1, A2, A3 or A4
  const orientation = "landscape"; // portrait or landscape

  const marginLeft = 40;
  const doc = new jsPDF(orientation, unit, size);

  doc.setFontSize(15);

  let content = {
    startY: 50,
    head: headers,
    body: data,
  };

  doc.text(title, marginLeft, 40);
  doc.autoTable(content);
  doc.save(getFilename(title));
};

const getFilename = (title) => {
  // For todays date;
  // eslint-disable-next-line no-extend-native
  Date.prototype.today = function () {
    return (
      (this.getDate() < 10 ? "0" : "") +
      this.getDate() +
      (this.getMonth() + 1 < 10 ? "0" : "") +
      (this.getMonth() + 1) +
      this.getFullYear()
    );
  };

  // For the time now
  // eslint-disable-next-line no-extend-native
  Date.prototype.timeNow = function () {
    return (
      (this.getHours() < 10 ? "0" : "") +
      this.getHours() +
      (this.getMinutes() < 10 ? "0" : "") +
      this.getMinutes() +
      (this.getSeconds() < 10 ? "0" : "") +
      this.getSeconds()
    );
  };
  const currentdate = new Date();

  return `${title}_${currentdate.today()}_${currentdate.timeNow()}`;
};
