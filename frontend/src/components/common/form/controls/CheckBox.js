import {
  FormControl,
  FormControlLabel,
  Checkbox as MuiCheckbox,
} from "@mui/material";
import React from "react";

function CheckBox(props) {
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
          <MuiCheckbox
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

export default CheckBox;
