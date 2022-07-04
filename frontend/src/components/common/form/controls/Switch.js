import {
  FormControl,
  FormControlLabel,
  Switch as MuiSwitch,
} from "@mui/material";
import React from "react";

function Switch(props) {
  const { name, label, value, onChange } = props;

  const convertDefaultEventPara = (name, value) => ({
    target: {
      name,
      value,
    },
  });

  return (
    <FormControl>
      <FormControlLabel
        control={
          <MuiSwitch
            name={name}
            checked={value}
            onChange={(e) =>
              onChange(convertDefaultEventPara(name, e.target.checked))
            }
          />
        }
        label={label}
      />
    </FormControl>
  );
}

export default Switch;
