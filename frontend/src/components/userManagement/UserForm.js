import { Grid, Stack } from "@mui/material";
import React, { useEffect } from "react";
import { Form, useForm } from "../common/form/useForm";
import Controls from "../common/form/controls/Controls";

const initialFValues = {
  id: "0",
  name: "",
  username: "",
  email: "",
  role: "",
  password: "",
  confirmPassword: "",
};

const items = [
  { id: "admin", title: "Admin" },
  { id: "supervisor", title: "Supervisor" },
  { id: "operator", title: "Operator" },
];

const UserForm = ({ AddorEdit, recordsForEdit }) => {
  const validate = (fieldValues = values) => {
    let temp = { ...errors };
    if ("name" in fieldValues)
      temp.name = fieldValues.name ? "" : "This field is required";
    if ("username" in fieldValues)
      temp.username = fieldValues.username ? "" : "This field is required";
    if ("email" in fieldValues) {
      temp.email = fieldValues.email
        ? /^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9-]+(?:\.[a-zA-Z0-9-]+)*$/.test(
            fieldValues.email
          )
          ? ""
          : "Email is not valid"
        : "This field is required";
    }
    if ("role" in fieldValues) {
      temp.role = fieldValues.role.length !== 0 ? "" : "This field is required";
    }
    if ("password" in fieldValues) {
      temp.password = fieldValues.password ? "" : "This field is required";
      temp.confirmPassword = values.confirmPassword
        ? fieldValues.password === values.confirmPassword
          ? ""
          : "Password and confirm should match"
        : "";
    }
    if ("confirmPassword" in fieldValues) {
      temp.confirmPassword = fieldValues.confirmPassword
        ? values.password === fieldValues.confirmPassword
          ? ""
          : "Password and confirm should match"
        : "This field is required";
    }

    seterrors({ ...temp });

    if (fieldValues === values)
      return Object.values(temp).every((x) => x === "");
  };

  const { values, setvalues, errors, seterrors, handleInputChange, resetForm } =
    useForm(initialFValues, true, validate);

  const handleSubmit = (e) => {
    e.preventDefault();
    if (validate()) {
      AddorEdit(values, resetForm);
      console.log(values);
    }
  };

  useEffect(() => {
    if (recordsForEdit !== null)
      setvalues({
        ...recordsForEdit,
      });
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [recordsForEdit]);

  return (
    <Form onSubmit={handleSubmit}>
      <Grid container columnSpacing={2}>
        <Grid item sm={6}>
          <Controls.Input
            name="name"
            label="Full Name"
            value={values.name}
            onChange={handleInputChange}
            error={errors.name}
          />
        </Grid>
        <Grid item sm={6}>
          <Controls.Input
            name="email"
            label="Email"
            value={values.email}
            onChange={handleInputChange}
            error={errors.email}
          />
        </Grid>
        <Grid item sm={6}>
          <Controls.Input
            name="username"
            label="Username"
            value={values.username}
            onChange={handleInputChange}
            error={errors.username}
          />
        </Grid>
        <Grid item sm={6}>
          <Controls.Select
            name="role"
            label="Role"
            value={values.role}
            onChange={handleInputChange}
            options={items}
            error={errors.role}
          />
        </Grid>
        <Grid item sm={6}>
          {values.id === "0" && (
            <Controls.Input
              name="password"
              label="Password"
              type="password"
              value={values.password}
              onChange={handleInputChange}
              error={errors.password}
            />
          )}
        </Grid>
        <Grid item sm={6}>
          {values.id === "0" && (
            <Controls.Input
              name="confirmPassword"
              label="Confirm Password"
              type="password"
              value={values.confirmPassword}
              onChange={handleInputChange}
              error={errors.confirmPassword}
            />
          )}
          <Stack direction="row" spacing={2}>
            <Controls.Button type="submit" text="Submit" />
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

export default UserForm;
