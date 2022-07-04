import { Visibility, VisibilityOff } from "@mui/icons-material";
import { IconButton, InputAdornment, TextField } from "@mui/material";
import React, { useState } from "react";

function Input(props) {
  const { name, label, value, onChange, error = null, type, fullWidth } = props;
  const [showPassword, setShowPassword] = useState(false);
  const handleClickShowPassword = () => {
    setShowPassword(!showPassword);
  };
  const handleMouseDownPassword = (event) => {
    event.preventDefault();
  };
  const EndAdornmentProps = () => {
    return (
      <InputAdornment position="end">
        <IconButton
          aria-label="toggle password visibility"
          onClick={handleClickShowPassword}
          onMouseDown={handleMouseDownPassword}
          edge="end"
        >
          {showPassword ? <VisibilityOff /> : <Visibility />}
        </IconButton>
      </InputAdornment>
    );
  };
  return (
    <TextField
      margin="normal"
      variant="outlined"
      type={showPassword ? "text" : type === "password" ? "password" : "text"}
      name={name}
      label={label}
      value={value}
      size="small"
      onChange={onChange}
      fullWidth={fullWidth || false}
      {...(error && { error: true, helperText: error })}
      {...(type === "password" && {
        InputProps: { ...{ endAdornment: <EndAdornmentProps /> } },
      })}
    ></TextField>
  );
}

export default Input;
