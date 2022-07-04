import { Close, Search } from "@mui/icons-material";
import {
  IconButton,
  InputAdornment,
  Stack,
  TextField,
  Toolbar,
  Tooltip,
  Typography,
} from "@mui/material";
import React from "react";

const TblToolbar = ({
  title,
  globalFilter,
  inputSearch,
  setInputSearch,
  toolbarAction,
  selected,
}) => {
  const handleMouseDownPassword = (event) => {
    event.preventDefault();
  };
  return (
    <>
      <Toolbar
        sx={{
          pl: { sm: 2 },
          pr: { xs: 1, sm: 1 },
        }}
      >
        <Typography
          sx={{ flexGrow: "1" }}
          variant="h6"
          id="tableTitle"
          component="div"
        >
          {title}
        </Typography>
        <Stack direction="row" spacing={2}>
          {globalFilter && (
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
        </Stack>
      </Toolbar>
      {toolbarAction.length > 0 && (
        <Toolbar
          sx={{
            pl: { sm: 2 },
            pr: { xs: 1, sm: 1 },
            justifyContent: "center",
          }}
        >
          <Stack direction="row" spacing={1}>
            {toolbarAction.map((toolbarItem) => {
              return (
                <Tooltip title={toolbarItem.title} key={toolbarItem.name}>
                  <IconButton
                    aria-label={toolbarItem.name}
                    size="large"
                    onClick={(e) => toolbarItem.onClick(e, selected)}
                  >
                    {toolbarItem.icon}
                  </IconButton>
                </Tooltip>
              );
            })}
          </Stack>
        </Toolbar>
      )}
    </>
  );
};

export default TblToolbar;
