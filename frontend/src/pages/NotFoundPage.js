import { Box, Button, Typography } from "@mui/material";
import React from "react";
import { useNavigate } from "react-router-dom";

const NotFoundPage = () => {
  const navigate = useNavigate();
  return (
    <Box sx={{ textAlign: "center" }}>
      <Typography variant="h2" gutterBottom component="div">
        Page you have visited not fond !
      </Typography>
      <Button variant="contained" onClick={() => navigate("/")}>
        go to home
      </Button>
    </Box>
  );
};

export default NotFoundPage;
