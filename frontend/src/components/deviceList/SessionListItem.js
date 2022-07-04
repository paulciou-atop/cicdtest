import { Cancel, CheckCircle } from "@mui/icons-material";
import {
  CircularProgress,
  ListItem,
  ListItemIcon,
  ListItemText,
  Typography,
} from "@mui/material";
import React from "react";

const SessionListItem = ({ status, message, label }) => {
  return (
    <ListItem alignItems="flex-start">
      <ListItemIcon>
        {status === "loading" && <CircularProgress size={28} color="info" />}
        {status === "done" && <CheckCircle color="success" fontSize="large" />}
        {status === "failed" && <Cancel color="error" fontSize="large" />}
      </ListItemIcon>
      <ListItemText
        primary={label}
        secondary={
          <Typography
            sx={{ display: "inline" }}
            component="span"
            variant="body2"
            color="text.secondary"
          >
            {message}
          </Typography>
        }
      />
    </ListItem>
  );
};

export default SessionListItem;
