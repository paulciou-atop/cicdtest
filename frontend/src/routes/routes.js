import React from "react";
import { Navigate, useRoutes } from "react-router-dom";
//Layout import
import MainLayout from "../layout/MainLayout";
//Pages import
import AccountUserProfile from "../pages/account/AccountUserProfile";
import DashboardDevices from "../pages/dashboard/DashboardDevices";
import DashboardOverview from "../pages/dashboard/DashboardOverview";
import LoginPage from "../pages/LoginPage";
import DeviceConfiguration from "../pages/management/DeviceConfiguration";
import ManagementDevice from "../pages/management/ManagementDevice";
import ManagementUser from "../pages/management/ManagementUser";
import { InventoryDevice } from "../pages/inventory/InventoryDevice";
import NotFoundPage from "../pages/NotFoundPage";
import PrivateRoute from "./PrivateRoute";


const Router = () => {
  return useRoutes([
    { path: "/", element: <Navigate to="/dashboard/overview" /> },
    {
      path: "/dashboard",
      element: <MainLayout />,
      children: [
        { index: true, element: <Navigate to="/dashboard/overview" /> },
        {
          path: "overview",
          element: (
            <PrivateRoute>
              <DashboardOverview />
            </PrivateRoute>
          ),
        },
        {
          path: "devices",
          element: (
            <PrivateRoute>
              <DashboardDevices />
            </PrivateRoute>
          ),
        },
      ],
    },
    {
      path: "/management",
      element: <MainLayout />,
      children: [
        { index: true, element: <Navigate to="/management/users" /> },
        {
          path: "users",
          element: (
            <PrivateRoute>
              <ManagementUser />
            </PrivateRoute>
          ),
        },
        {
          path: "devices",
          element: (
            <PrivateRoute>
              <ManagementDevice />
            </PrivateRoute>
          ),
        },
        {
          path: "deviceConfig",
          element: (
            <PrivateRoute>
              <DeviceConfiguration />
            </PrivateRoute>
          ),
        },
      ],
    },
    {
      path: "/inventory",
      element: <MainLayout />,
      children: [
        { index: true, element: <Navigate to="/inventory/devices" /> },
        {
          path: "devices",
          element: (
            <PrivateRoute>
              <InventoryDevice />
            </PrivateRoute>
          ),
        },
      ],
    },
    {
      path: "/account",
      element: <MainLayout />,
      children: [
        {
          index: true,
          element: (
            <PrivateRoute>
              <AccountUserProfile />
            </PrivateRoute>
          ),
        },
      ],
    },
    { path: "/login", element: <LoginPage /> },
    { path: "/404", element: <NotFoundPage /> },
    { path: "*", element: <Navigate to="/404" /> },
  ]);
};

export default Router;
