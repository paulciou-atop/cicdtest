import { createSlice } from "@reduxjs/toolkit";

const initialState = {
  palette: {
    mode: "light",
    primary: {
      main: "#13c2c2",
      contrastText: "#fff",
    },
    secondary: {
      main: "#f50057",
      contrastText: "#fff",
    },
    background: {
      default: "#dfdfdf",
      paper: "#fff",
    },
  },
  components: {
    MuiAppBar: {
      styleOverrides: {
        root: {
          backgroundColor: "#001529",
          color: "#fff",
        },
        colorInherit: {
          color: "#fff",
        },
      },
    },
  },
};
const nmsTheme = JSON.parse(localStorage.getItem("nms-theme"));
const isThemeSaved = nmsTheme && Object.keys(nmsTheme).length !== 0;

export const ThemeSlice = createSlice({
  name: "theme",
  initialState: isThemeSaved ? nmsTheme : initialState,
  reducers: {
    changeThemeMode: (state, { payload }) => {
      state.palette.mode = payload;
      state.palette.background.default =
        payload === "dark" ? "#181818" : "#dfdfdf";
      state.palette.background.paper = payload === "dark" ? "#303030" : "#fff";
      localStorage.setItem("nms-theme", JSON.stringify(state));
    },
  },
});

export const { changeThemeMode } = ThemeSlice.actions;
export const themeSelector = (state) => state.theme;
