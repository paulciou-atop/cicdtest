import { createSlice } from "@reduxjs/toolkit";

const DeviceAdvancedSetting = createSlice({
  name: "DeviceAdvancedSetting",
  initialState: {
    visible: false,
  },
  reducers: {
    openAdvancedSettingDrawer: (state, { payload }) => {
      state.visible = true;
    },
    closeAdvancedSettingDrawer: (state, { payload }) => {
      state.visible = false;
    },
  },
});

export const { openAdvancedSettingDrawer, closeAdvancedSettingDrawer } =
  DeviceAdvancedSetting.actions;

export const deviceAdvancedSettingSelector = (state) => {
  const { visible } = state.deviceAdvancedSetting;
  return { visible };
};

export default DeviceAdvancedSetting;
