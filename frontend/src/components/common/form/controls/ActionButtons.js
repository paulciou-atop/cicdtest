import { Button } from "@mui/material";
import { makeStyles } from "@mui/styles";
import React from "react";

const useStyles = makeStyles((theme) => ({
  root: {
    minWidth: 0,
    margin: theme.spacing(0.5),
  },
  secondary: {
    backgroundColor: theme.palette.secondary.light,
    color: theme.palette.secondary.dark,
  },
  primary: {
    backgroundColor: theme.palette.primary.light,
    color: theme.palette.primary.dark,
  },
}));

const ActionButtons = ({ color, children, onClick, ...others }) => {
  const classes = useStyles();
  return (
    <Button
      onClick={onClick}
      className={`${classes.root} ${classes[color]}`}
      {...others}
    >
      {children}
    </Button>
  );
};

export default ActionButtons;
