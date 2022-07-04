import { createAsyncThunk, createSlice } from "@reduxjs/toolkit";
import publicApis from "../utils/apis/publicApis";

export const GetSessionID = createAsyncThunk(
  "deviceDiscovery/GetSessionID",
  async (ipRange, thunkAPI) => {
    try {
      const response = await publicApis.post("/v1/scan/start", {
        range: ipRange,
        severIP: "10.6.50.90",
      });
      const data = await response.data;
      if (response.status === 200) {
        return data;
      } else {
        return thunkAPI.rejectWithValue(data.payload.reason);
      }
    } catch (e) {
      if (e.response) return thunkAPI.rejectWithValue(e.response.statusText);
      else return thunkAPI.rejectWithValue(e.message);
    }
  }
);

export const GetCheckStatus = createAsyncThunk(
  "deviceDiscovery/GetCheckStatus",
  async (_, thunkAPI) => {
    try {
      const { sessionId } = thunkAPI.getState().deviceDiscovery;
      const response = await publicApis.post("/v1/scan/check", {
        sessionId,
      });
      const data = await response.data;
      if (response.status === 200) {
        if (data?.info.status === "success") {
          thunkAPI.dispatch(GetDeviceData());
        }
        if (data?.info.status === "fail") {
          thunkAPI.dispatch(GetDeviceData());
          return thunkAPI.rejectWithValue("fail");
        }
        return data;
      } else {
        return thunkAPI.rejectWithValue(data.payload.reason);
      }
    } catch (e) {
      if (e.response) return thunkAPI.rejectWithValue(e.response.statusText);
      else return thunkAPI.rejectWithValue(e.message);
    }
  }
);

export const GetDeviceData = createAsyncThunk(
  "deviceDiscovery/GetDeviceData",
  async (_, thunkAPI) => {
    try {
      const { sessionId } = thunkAPI.getState().deviceDiscovery;
      const response = await publicApis.post("/v1/scan/result", {
        sessionId,
        page: 1,
        size: 20,
      });
      const data = await response.data;
      if (response.status === 200) {
        return data;
      } else {
        return thunkAPI.rejectWithValue(data.payload.reason);
      }
    } catch (e) {
      console.log("Error", e.message);
      if (e.response) return thunkAPI.rejectWithValue(e.response.statusText);
      else return thunkAPI.rejectWithValue(e.message);
    }
  }
);

export const GetlastScannedDeviceData = createAsyncThunk(
  "deviceDiscovery/GetlastScannedDeviceData",
  async (_, thunkAPI) => {
    try {
      const response = await publicApis.post("/v1/scan/result/last", {
        page: 1,
        size: 20,
      });
      const data = await response.data;
      if (response.status === 200) {
        return data;
      } else {
        return thunkAPI.rejectWithValue(data?.info.message);
      }
    } catch (e) {
      console.log("Error", e.message);
      if (e.response) return thunkAPI.rejectWithValue(e.response.statusText);
      else return thunkAPI.rejectWithValue(e.message);
    }
  }
);

const DeviceDiscoverySlice = createSlice({
  name: "deviceDiscovery",
  initialState: {
    sessionId: "",
    titleSessionId: "",
    deviceData: [],
    sessionDetails: {
      statusStart: "loading",
      messageStart: "no data",
      statusCheck: "loading",
      messageCheck: "no data",
      statusResult: "loading",
      messageResult: "no data",
    },
  },
  //reducers: {},
  extraReducers: {
    [GetSessionID.pending]: (state, { payload }) => {
      return {
        ...state,
        sessionDetails: {
          statusStart: "loading",
          messageStart: "no data",
          statusCheck: "loading",
          messageCheck: "no data",
          statusResult: "loading",
          messageResult: "no data",
        },
      };
    },
    [GetSessionID.fulfilled]: (state, { payload }) => {
      return {
        ...state,
        sessionId: payload?.info.sessionId,
        sessionDetails: {
          ...state.sessionDetails,
          statusStart: "done",
          messageStart: payload?.info.sessionId,
        },
      };
    },
    [GetSessionID.rejected]: (state, { payload }) => {
      return {
        ...state,
        sessionId: "",
        sessionDetails: {
          ...state.sessionDetails,
          statusStart: "failed",
          messageStart: "Failed to get sessionid",
        },
      };
    },
    [GetCheckStatus.pending]: (state, { payload }) => {
      return {
        ...state,
        sessionDetails: {
          ...state.sessionDetails,
          statusCheck: "loading",
          messageCheck: "running",
          statusResult: "loading",
          messageResult: "no data",
        },
      };
    },
    [GetCheckStatus.fulfilled]: (state, { payload }) => {
      let status = payload?.info.status;
      return {
        ...state,
        sessionDetails: {
          ...state.sessionDetails,
          statusCheck: status === "running" ? "loading" : "done",
          messageCheck: status,
        },
      };
    },
    [GetCheckStatus.rejected]: (state, { payload }) => {
      return {
        ...state,
        sessionDetails: {
          ...state.sessionDetails,
          statusCheck: "failed",
          messageCheck: "failed",
        },
      };
    },

    [GetDeviceData.pending]: (state, { payload }) => {
      return {
        ...state,
        sessionDetails: {
          ...state.sessionDetails,
          statusResult: "loading",
          messageResult: "no data",
        },
      };
    },
    [GetDeviceData.fulfilled]: (state, { payload }) => {
      const titleSessionId =
        payload?.content && payload?.content.length > 0
          ? payload?.content[0].sessionId
          : state.titleSessionId;
      return {
        ...state,
        titleSessionId,
        deviceData:
          payload?.content.length > 0 ? payload?.content : state.deviceData,
        sessionDetails: {
          ...state.sessionDetails,
          statusResult: "done",
          messageResult: `Found ${payload?.content.length} device data`,
        },
      };
    },
    [GetDeviceData.rejected]: (state, { payload }) => {
      return {
        ...state,
        sessionDetails: {
          ...state.sessionDetails,
          statusResult: "failed",
          messageResult: "no data",
        },
      };
    },
    [GetlastScannedDeviceData.rejected]: (state, { payload }) => {
      return {
        ...state,
        deviceData: [],
      };
    },
    [GetlastScannedDeviceData.fulfilled]: (state, { payload }) => {
      const titleSessionId =
        payload?.content && payload?.content.length > 0
          ? payload?.content[0].sessionId
          : "";
      return {
        ...state,
        titleSessionId,
        deviceData: payload?.content,
      };
    },
  },
});

//export const {} = DeviceDiscoverySlice.actions;

export const deviceDiscoverySelector = (state) => state.deviceDiscovery;

export default DeviceDiscoverySlice;
