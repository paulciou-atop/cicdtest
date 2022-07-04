import { configureStore } from "@reduxjs/toolkit";
import DeviceAdvancedSetting from "../features/DeviceAdvancedSetting";
import DeviceDiscoverySlice from "../features/DeviceDiscoverySlice";
import SingleNetworkSetting from "../features/SingleNetworkSetting";
import { ThemeSlice } from "../features/ThemeSlice";

export const store = configureStore({
  reducer: {
    theme: ThemeSlice.reducer,
    deviceDiscovery: DeviceDiscoverySlice.reducer,
    singleNetworkSetting: SingleNetworkSetting.reducer,
    deviceAdvancedSetting: DeviceAdvancedSetting.reducer,
  },
});
