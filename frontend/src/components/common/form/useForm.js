import { makeStyles } from "@mui/styles";
import React, { useState } from "react";

export function useForm(initialFValues, validateOnChange = false, validate) {
  const [values, setvalues] = useState(initialFValues);
  const [errors, seterrors] = useState({});
  const handleInputChange = (e) => {
    const { name, value } = e.target;
    setvalues({ ...values, [name]: value });
    if (validateOnChange) validate({ [name]: value });
  };
  const resetForm = () => {
    setvalues(initialFValues);
    seterrors({});
  };
  return { values, setvalues, errors, seterrors, handleInputChange, resetForm };
}

const useStyles = makeStyles((theme) => ({
  root: {
    "& .MuiFormControl-root": {
      width: "100%",
    },
  },
}));
export function Form(props) {
  const classes = useStyles();
  const { children, ...others } = props;
  return (
    <form className={classes.root} autoComplete="off" {...others}>
      {children}
    </form>
  );
}
