import { Card, CardContent, Typography } from "@mui/material";
import React from "react";

const SummaryCard = ({ hbcolor, bbcolor, hlabel, blabel }) => {
  return (
    <Card elevation={0} sx={{ flexGrow: 1 }}>
      <CardContent sx={{ p: 0, pb: "0px !important" }}>
        <Typography
          variant="h6"
          component="div"
          textAlign="center"
          sx={{ backgroundColor: hbcolor, color: "#ffffff" }}
        >
          {hlabel}
        </Typography>
        <Typography
          variant="h5"
          color="text.secondary"
          component="div"
          textAlign="center"
          sx={{ backgroundColor: bbcolor, color: "#303030" }}
        >
          {blabel}
        </Typography>
      </CardContent>
    </Card>
  );
};

export default SummaryCard;
