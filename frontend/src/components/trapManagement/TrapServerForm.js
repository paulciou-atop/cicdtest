import { Grid, Stack } from "@mui/material";
import React from "react";
import { Form, useForm } from "../common/form/useForm";
import Controls from "../common/form/controls/Controls";

const initialFValues = {
  id: "0",
  trapipserver: "",
  serviceport: "",
  trapcomm: "",
};

const TrapServerForm = ({ AddorEdit }) => {
  const validate = (fieldValues = values) => {
    let temp = { ...errors };
    if ("trapipserver" in fieldValues)
      temp.trapipserver = fieldValues.trapipserver ? "" : "This field is required";
    if ("serviceport" in fieldValues) {
      temp.serviceport = fieldValues.serviceport ? "" : "This field is required";
    }
    if ("trapcomm" in fieldValues) {
     temp.trapcomm = fieldValues.trapcomm ? "" : "This field is required";
    }

    seterrors({ ...temp });

    if (fieldValues === values)
      return Object.values(temp).every((x) => x === "");
  };

   const { values, errors, seterrors, handleInputChange, resetForm } =
    useForm(initialFValues, true, validate);

  const handleSubmit = (e) => {
    e.preventDefault();
    if (validate()) {
      AddorEdit(values, resetForm);
    }
  };

  return (
    <Form onSubmit={handleSubmit}>
      <Grid container columnSpacing={2}>
        <Grid item sm={6}>
          <Controls.Input
            required
            name="trapipserver"
            label="Trap Server IP"
            value={values.trapipserver}
            onChange={handleInputChange}
            error={errors.trapipserver}
          />
        </Grid>

        <Grid item sm={6}>
          <Controls.Input
            required
            name="serviceport"
            label="Trap Server Port"
            value={values.serviceport}
            onChange={handleInputChange}
            error={errors.serviceport}
          />
        </Grid>

        <Grid item sm={6}>
          <Controls.Input
            required
            name="trapcomm"
            label="Trap Comm String"
            value={values.trapcomm}
            onChange={handleInputChange}
            error={errors.trapcomm}
          />
        </Grid>

        <Grid item sm={6} pt={2}>
          <Stack direction="row" spacing={2}>
            <Controls.Button type="submit" text="Start" />
            <Controls.Button
              color="secondary"
              varient="outlined"
              text="Reset"
              onClick={resetForm}
            />
          </Stack>
          </Grid>
      </Grid>
    </Form>
  );
};

export default TrapServerForm;
