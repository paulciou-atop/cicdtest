import { Paper, Stack, Typography } from "@mui/material";
import React from "react";
import SummaryCard from "./SummaryCard";

const NotificationSummary = () => {
  return (
    <Paper sx={{ padding: "10px" }}>
      <Typography variant="h6" component="div" textAlign="center">
        Today's notifications count
      </Typography>
      <Stack direction="row" spacing={2} justifyContent="space-evenly">
        <SummaryCard
          hbcolor="#00E396"
          bbcolor="#e6fff6"
          hlabel="Information"
          blabel={54}
        />
        <SummaryCard
          hbcolor="#FEB019"
          bbcolor="#fff5e1"
          hlabel="Warning"
          blabel={87}
        />
        <SummaryCard
          hbcolor="#FF4560"
          bbcolor="#fde2e6"
          hlabel="Critical"
          blabel={87}
        />
      </Stack>
    </Paper>
  );
};

export default NotificationSummary;
