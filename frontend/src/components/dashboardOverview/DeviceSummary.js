import { Paper, Stack, Typography } from "@mui/material";
import React from "react";
import SummaryCard from "./SummaryCard";

const DeviceSummary = () => {
  return (
    <Paper sx={{ padding: "10px" }}>
      <Typography variant="h6" component="div" textAlign="center">
        Online/Offline devices
      </Typography>
      <Stack direction="row" spacing={2} justifyContent="space-evenly">
        <SummaryCard
          hbcolor="#46b300"
          bbcolor="#E8F5E9"
          hlabel="Online"
          blabel={54}
        />
        <SummaryCard
          hbcolor="#D50000"
          bbcolor="#FFEBEE"
          hlabel="Offline"
          blabel={87}
        />
      </Stack>
    </Paper>
  );
};

export default DeviceSummary;
