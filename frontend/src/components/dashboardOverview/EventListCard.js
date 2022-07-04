import { Grid, Paper } from "@mui/material";
import React from "react";

const EventListCard = ({ ledColor }) => {
  return (
    <Paper sx={{ p: "10px" }}>
      <Grid container sx={{ mb: "5px" }}>
        <Grid item xs={6} textAlign="left">
          Switch 1
        </Grid>
        <Grid item xs={6} textAlign="right">
          2009-10-31T01:48:52Z
        </Grid>
      </Grid>
      <Grid container>
        <Grid item flex="auto">
          <Grid container>
            <Grid item xs={12} textAlign="left">
              10.6.50.7
            </Grid>
            <Grid item xs={12} textAlign="center">
              port link-up
            </Grid>
          </Grid>
        </Grid>
        <Grid item flex="60px" textAlign="center" sx={{ m: "auto" }}>
          <svg height="40" width="40">
            <circle
              cx="20"
              cy="20"
              r="17"
              stroke="inherit"
              strokeWidth="2"
              fill={ledColor}
            />
            Sorry, your browser does not support inline SVG.
          </svg>
        </Grid>
      </Grid>
    </Paper>
  );
};

export default EventListCard;
