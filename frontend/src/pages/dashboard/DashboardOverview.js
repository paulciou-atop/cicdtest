import { Grid, Paper, Stack, Typography } from "@mui/material";
import React from "react";
import DeviceSummary from "../../components/dashboardOverview/DeviceSummary";
import DiskspceUtilization from "../../components/dashboardOverview/DiskspceUtilization";
import EventListCard from "../../components/dashboardOverview/EventListCard";
import NotificationChart from "../../components/dashboardOverview/NotificationChart";
import NotificationSummary from "../../components/dashboardOverview/NotificationSummary";
import SyslogChart from "../../components/dashboardOverview/SyslogChart";
import TrapChart from "../../components/dashboardOverview/TrapChart";

const DashboardOverview = () => {
  return (
    <Grid container spacing={[2, 2]}>
      <Grid item xs={8}>
        <Grid container spacing={[2, 2]}>
          <Grid item xs={4}>
            <DeviceSummary />
          </Grid>
          <Grid item xs={8}>
            <NotificationSummary />
          </Grid>
          <Grid item xs={4}>
            <DiskspceUtilization />
          </Grid>
          <Grid item xs={8}>
            <NotificationChart />
          </Grid>
          <Grid item xs={6}>
            <SyslogChart />
          </Grid>
          <Grid item xs={6}>
            <TrapChart />
          </Grid>
        </Grid>
      </Grid>
      <Grid item xs={4}>
        <div>
          <Paper sx={{ mb: "10px" }}>
            <Typography variant="h6" component="div" textAlign="center">
              Today's event list
            </Typography>
          </Paper>
          <Stack spacing={1}>
            <EventListCard ledColor="#00E396" />
            <EventListCard ledColor="#FEB019" />
            <EventListCard ledColor="#FF4560" />
            <EventListCard ledColor="#00E396" />
            <EventListCard ledColor="#FEB019" />
            <EventListCard ledColor="#FF4560" />
            <EventListCard ledColor="#00E396" />
            <EventListCard ledColor="#FEB019" />
            <EventListCard ledColor="#FF4560" />
          </Stack>
        </div>
      </Grid>
    </Grid>
  );
};

export default DashboardOverview;
