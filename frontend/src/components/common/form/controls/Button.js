import React from "react";
import { Button as MuiButton } from "@mui/material";
import { makeStyles } from "@mui/styles";

const useStyles = makeStyles((theme) => ({
  label: {
    textTransform: "none",
  },
}));

function Button(props) {
  const { text, size, varient, color, onClick, ...other } = props;
  const classes = useStyles();
  return (
    <MuiButton
      variant={varient || "contained"}
      color={color || "primary"}
      size={size || "large"}
      onClick={onClick}
      {...other}
      classes={{ label: classes.label }}
    >
      {text}
    </MuiButton>
  );
}

export default Button;
