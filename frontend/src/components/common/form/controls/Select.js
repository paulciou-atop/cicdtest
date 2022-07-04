import {
  FormControl,
  FormHelperText,
  InputLabel,
  MenuItem,
  Select as MuiSelect,
} from "@mui/material";
import React from "react";

const Select = ({
  name,
  label,
  value,
  onChange,
  options,
  error = null,
  fullWidth,
}) => {
  return (
    <FormControl
      variant="outlined"
      {...(error && { error: true })}
      sx={{ mt: 2, mb: 1 }}
      size="small"
      fullWidth={fullWidth || false}
    >
      <InputLabel>{label}</InputLabel>
      <MuiSelect label={label} name={name} value={value} onChange={onChange}>
        <MenuItem value="">None</MenuItem>
        {options.map((option) => (
          <MenuItem key={option.id} value={option.id}>
            {option.title}
          </MenuItem>
        ))}
      </MuiSelect>
      {error && <FormHelperText>{error}</FormHelperText>}
    </FormControl>
  );
};

export default Select;
