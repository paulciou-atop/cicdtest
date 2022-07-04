import { Paper } from "@mui/material";
import ReactApexChart from "react-apexcharts";
import { useTheme } from "@mui/material/styles";
import React, { useEffect, useState } from "react";

const DiskspceUtilization = () => {
  const theme = useTheme();
  const [graphConfig, setGraphConfig] = useState({
    series: [278, 197],
    options: {
      chart: {
        type: "donut",
        background: "transparant",
      },
      colors: ["#00E396", "#FEB019"],
      theme: {
        mode: theme.palette.mode,
      },
      labels: ["Free Space", "Occupied Space"],
      legend: {
        position: "bottom",
      },
      title: {
        text: "Disk Space Utilization",
      },
      stroke: {
        show: false,
      },
      plotOptions: {
        pie: {
          expandOnClick: false,
          donut: {
            labels: {
              show: true,
              name: {
                show: true,
              },
              value: {
                show: true,
              },
              total: {
                showAlways: false,
                show: true,
              },
            },
          },
        },
      },
    },
  });

  useEffect(() => {
    setGraphConfig((prev) => ({
      ...prev,
      options: { ...prev.options, theme: { mode: theme.palette.mode } },
    }));
  }, [theme.palette.mode]);

  return (
    <Paper sx={{ padding: "10px", minHeight: "385px" }}>
      <div id="diskuti-chart">
        <ReactApexChart
          options={graphConfig.options}
          series={graphConfig.series}
          type="donut"
          height={350}
        />
      </div>
    </Paper>
  );
};

export default DiskspceUtilization;
