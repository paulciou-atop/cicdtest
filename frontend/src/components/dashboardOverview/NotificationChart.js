import { Paper } from "@mui/material";
import ReactApexChart from "react-apexcharts";
import { useTheme } from "@mui/material/styles";
import React, { useEffect, useState } from "react";

const NotificationChart = () => {
  const theme = useTheme();
  const [chartData, setChartData] = useState({
    series: [
      {
        name: "Information",
        data: [44, 55, 57, 56, 61, 58, 63],
      },
      {
        name: "Warning",
        data: [76, 85, 101, 98, 87, 105, 91],
      },
      {
        name: "Critical",
        data: [35, 41, 36, 26, 45, 48, 52],
      },
    ],
    options: {
      chart: {
        type: "bar",
        background: "transparant",
        height: 350,
      },
      //   grid: {
      //     padding: {
      //       top: 0,
      //       right: 0,
      //       bottom: 0,
      //       left: 0,
      //     },
      //   },
      colors: ["#00E396", "#FEB019", "#FF4560"],
      theme: {
        mode: theme.palette.mode,
      },
      title: {
        text: "Notification Chart",
      },
      plotOptions: {
        bar: {
          horizontal: false,
          columnWidth: "55%",
          endingShape: "rounded",
        },
      },
      dataLabels: {
        enabled: false,
      },
      stroke: {
        show: true,
        width: 2,
        colors: ["transparent"],
      },
      xaxis: {
        categories: [
          "14/05/22",
          "15/05/22",
          "16/05/22",
          "17/05/22",
          "18/05/22",
          "19/05/22",
          "20/05/22",
        ],
      },
      yaxis: {
        title: {
          text: "Notification Counts",
        },
      },
      fill: {
        opacity: 1,
      },
      // tooltip: {
      //   y: {
      //     formatter: function (val) {
      //       return "$ " + val + " thousands";
      //     },
      //   },
      // },
    },
  });

  useEffect(() => {
    setChartData((prev) => ({
      ...prev,
      options: { ...prev.options, theme: { mode: theme.palette.mode } },
    }));
  }, [theme.palette.mode]);

  return (
    <Paper sx={{ padding: "10px", minHeight: "385px" }}>
      <div id="notification-chart">
        <ReactApexChart
          options={chartData.options}
          series={chartData.series}
          type="bar"
          height={350}
        />
      </div>
    </Paper>
  );
};

export default NotificationChart;
