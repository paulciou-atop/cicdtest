import { Breadcrumbs, Link, Typography } from "@mui/material";
import { Link as RouterLink, useLocation } from "react-router-dom";
import navConfig from "./NavConfig";
import React from "react";

const BreadcrumItemLayout = () => {
  let breadcrumMap = {};
  const breadcrumbNameMap = (dataSource) => {
    dataSource.forEach((menu) => {
      if (menu.children) {
        breadcrumbNameMap(menu.children);
      }
      breadcrumMap = { ...breadcrumMap, [menu.path]: menu.title };
    });
    return breadcrumMap;
  };

  const LinkRouter = (props) => (
    <Link
      {...props}
      component={RouterLink}
      sx={{ textTransform: "capitalize" }}
    />
  );

  const BreadCrumbItem = () => {
    const location = useLocation();
    const pathnames = location.pathname.split("/").filter((x) => x);
    return (
      <Breadcrumbs aria-label="breadcrumb" sx={{ m: "10px 0" }}>
        <LinkRouter underline="hover" color="inherit" to="/">
          Home
        </LinkRouter>
        {pathnames.map((value, index) => {
          const last = index === pathnames.length - 1;
          const to = `/${pathnames.slice(0, index + 1).join("/")}`;
          return last ? (
            <Typography
              color="text.primary"
              key={to}
              sx={{ textTransform: "capitalize" }}
            >
              {breadcrumbNameMap(navConfig)[to]}
            </Typography>
          ) : (
            <LinkRouter underline="hover" color="inherit" to={to} key={to}>
              {breadcrumbNameMap(navConfig)[to]}
            </LinkRouter>
          );
        })}
      </Breadcrumbs>
    );
  };

  return <BreadCrumbItem />;
};

export default BreadcrumItemLayout;
