import { Dialog, DialogContent, DialogTitle, Typography } from "@mui/material";
import { makeStyles } from "@mui/styles";
import Controls from "./form/controls/Controls";
import React from "react";
import { Close } from "@mui/icons-material";

const useStyles = makeStyles((theme) => ({
  dialogWraper: {
    padding: theme.spacing(2),
    position: "absolute",
    top: theme.spacing(5),
  },
}));

const Popup = ({ title, children, openPopup, setOpenPopup }) => {
  const classes = useStyles();
  return (
    <Dialog
      open={openPopup}
      maxWidth="md"
      fullWidth
      classes={{ paper: classes.dialogWraper }}
    >
      <DialogTitle>
        <div style={{ display: "flex", alignItems: "center" }}>
          <Typography variant="h6" component="div" style={{ flexGrow: 1 }}>
            {title}
          </Typography>
          <Controls.ActionButtons
            color="secondary"
            onClick={() => setOpenPopup(false)}
          >
            <Close></Close>
          </Controls.ActionButtons>
        </div>
      </DialogTitle>
      <DialogContent dividers>{children}</DialogContent>
    </Dialog>
  );
};

export default Popup;
