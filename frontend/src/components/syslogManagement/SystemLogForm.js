import { Grid, Stack } from "@mui/material";
import React from "react";
import { Form, useForm } from "../common/form/useForm";
import Controls from "../common/form/controls/Controls";

const initialFValues = {
  id: "0",
  flash: "",
  level: "",
  server: "",
  ipserver: "",
  serviceport: "",
};

const items = [
  { id: "ERR", title: "0: (LOG_ERR)" },
  { id: "EMERG", title: "1: (LOG_EMERG" },
  { id: "ALERT", title: "2: (LOG_ALERT" },
  { id: "CRIT", title: "3: (LOG_CRIT)" },
  { id: "WARNING", title: "4: (LOG_WARNING)" },
  { id: "NOTICE", title: "5: (LOG_NOTICE)" },
  { id: "INFO", title: "6: (LOG_INFO)" },
  { id: "DEBUG", title: "7: (LOG_DEBUG)" },
];

const SystemLogForm = ({ AddorEdit }) => {
  const validate = (fieldValues = values) => {
    let temp = { ...errors };
    if ("server" in fieldValues)
      temp.server = fieldValues.server ? "" : "This field is required";
    if ("ipserver" in fieldValues) {
      temp.ipserver = fieldValues.ipserver ? "" : "This field is required";
    }
    if ("serviceport" in fieldValues) {
     temp.serviceport = fieldValues.serviceport ? "" : "This field is required";
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
          <Controls.Switch
            name="flash"
            label="Log to Flash"
            onChange={handleInputChange}
          />
        </Grid>

        <Grid item sm={6}>
          <Controls.Switch
            name="server"
            label="Log to Server"
            onChange={handleInputChange}
          />
        </Grid>
        
         <Grid item sm={6}>
          <Controls.Select
            name="level"
            label="Log Level"
            value={values.level}
            onChange={handleInputChange}
            displayEmpty
            options={items}
            error={errors.level}
          />
        </Grid>

        <Grid item sm={6}>
          <Controls.Input
            required
            name="ipserver"
            label="Server IP"
            value={values.ipserver}
            onChange={handleInputChange}
            error={errors.ipserver}
          />
        </Grid>

        <Grid item sm={6}>
          <Controls.Input
            required
            name="serviceport"
            label="Server Port"
            value={values.serviceport}
            onChange={handleInputChange}
            error={errors.serviceport}
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

export default SystemLogForm;
