import { Box, Grid, Typography } from "@mui/material";
import { ThemeProvider, createTheme, useTheme } from "@mui/material/styles";
import CssBaseline from "@mui/material/CssBaseline";
import logo from "../assets/images/logo-new.svg";
import React from "react";
import { Form, useForm } from "../components/common/form/useForm";
import Controls from "../components/common/form/controls/Controls";
import { useNavigate } from "react-router-dom";

const initialFValues = {
  username: "",
  password: "",
  rememberMe: false,
};

const LoginPage = () => {
  const navigate = useNavigate();
  const theme = useTheme();

  const validate = (fieldValues = values) => {
    let temp = { ...errors };
    if ("username" in fieldValues)
      temp.username = fieldValues.username ? "" : "This field is required";
    if ("password" in fieldValues)
      temp.password = fieldValues.password ? "" : "This field is required";

    seterrors({ ...temp });

    if (fieldValues === values)
      return Object.values(temp).every((x) => x === "");
  };

  const { values, errors, seterrors, handleInputChange } = useForm(
    initialFValues,
    true,
    validate
  );

  const handleSubmit = (e) => {
    e.preventDefault();
    if (validate()) {
      localStorage.setItem("token", true);
      setTimeout(() => {
        navigate("/");
      }, 100);
    }
  };
  return (
    <ThemeProvider
      theme={createTheme({
        palette: {
          mode: "light",
          primary: {
            main: theme.palette.primary.main,
            contrastText: theme.palette.primary.contrastText,
          },
          secondary: {
            main: theme.palette.secondary.main,
            contrastText: theme.palette.secondary.contrastText,
          },
        },
      })}
    >
      <CssBaseline />
      <Box
        style={{
          height: "100vh",
          display: "flex",
          alignItems: "center",
          justifyContent: "center",
          background: "#f7f7f7",
        }}
      >
        <Box
          sx={{
            display: "flex",
            flexDirection: "column",
            alignItems: "center",
            justifyContent: "center",
            maxWidth: "450px",
          }}
        >
          <Box
            component="img"
            src={logo}
            alt="logo"
            width={260}
            sx={{ marginBottom: "10px" }}
          />
          <Typography component="h1" variant="h5">
            Sign in
          </Typography>
          <Form onSubmit={handleSubmit}>
            <Grid container>
              <Grid item sm={12}>
                <Controls.Input
                  name="username"
                  label="Username"
                  value={values.username}
                  onChange={handleInputChange}
                  error={errors.username}
                />
                <Controls.Input
                  name="password"
                  label="Password"
                  type="password"
                  value={values.password}
                  onChange={handleInputChange}
                  error={errors.password}
                />
                <Controls.CheckBox
                  name="rememberMe"
                  label="Remember me"
                  value={values.rememberMe}
                  onChange={handleInputChange}
                />
                <div>
                  <Controls.Button type="submit" text="Submit" fullWidth />
                </div>
              </Grid>
            </Grid>
          </Form>
        </Box>
      </Box>
    </ThemeProvider>
  );
};

export default LoginPage;
