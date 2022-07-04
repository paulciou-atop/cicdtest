import { Paper } from "@mui/material";
import ReactApexChart from "react-apexcharts";
import { useTheme } from "@mui/material/styles";
import React, { useEffect, useState } from "react";

const SyslogChart = () => {
  const theme = useTheme();
  const [chartData, setChartData] = useState({
    series: [
      {
        name: "Syslog",
        data: [44, 55, 57, 56, 61, 58, 63],
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
      colors: ["#008FFB"],
      theme: {
        mode: theme.palette.mode,
      },
      title: {
        text: "Syslog Chart",
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
          text: "Syslog Counts",
        },
      },
      fill: {
        opacity: 1,
      },
      legend: {
        show: true,
        showForSingleSeries: true,
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
      <div id="syslog-chart">
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

export default SyslogChart;
