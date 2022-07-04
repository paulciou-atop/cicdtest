import React from "react";
import "./App.css";
import Router from "./routes/routes";
import { ThemeProvider, createTheme } from "@mui/material/styles";
import CssBaseline from "@mui/material/CssBaseline";
import { useSelector } from "react-redux";
import { themeSelector } from "./features/ThemeSlice";

function App() {
  const theme = createTheme(useSelector(themeSelector));
  return (
    <ThemeProvider theme={theme}>
      <CssBaseline />
      <Router />
    </ThemeProvider>
  );
}

export default App;
