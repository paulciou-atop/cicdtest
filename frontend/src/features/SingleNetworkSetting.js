import { createSlice } from "@reduxjs/toolkit";

const SingleNetworkSetting = createSlice({
  name: "SingleNetworkSetting",
  initialState: {
    visible: false,
  },
  reducers: {
    openNetworkSettingDrawer: (state, { payload }) => {
      state.visible = true;
    },
    closeNetworkSettingDrawer: (state, { payload }) => {
      state.visible = false;
    },
  },
});

export const { openNetworkSettingDrawer, closeNetworkSettingDrawer } =
  SingleNetworkSetting.actions;

export const singleNetworkSettingSelector = (state) => {
  const { visible } = state.singleNetworkSetting;
  return { visible };
};

export default SingleNetworkSetting;
