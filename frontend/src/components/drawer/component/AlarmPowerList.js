import {
  Checkbox,
  List,
  ListItem,
  Stack,
  Tooltip,
  Typography,
} from "@mui/material";
import React from "react";

const AlarmPowerList = ({ items }) => {
  const HeaderItemRender = () => (
    <Stack
      direction="row"
      justifyContent="space-between"
      alignItems="center"
      flexGrow={1}
    >
      <Typography variant="body1">Name</Typography>
      <Typography variant="body1">On</Typography>
      <Typography variant="body1">Off</Typography>
    </Stack>
  );
  const BodyItemRender = ({ item }) => (
    <Stack
      direction="row"
      justifyContent="space-between"
      alignItems="center"
      flexGrow={1}
    >
      <Typography variant="body1">{`${item.powerName}`}</Typography>
      <Tooltip title={`When ${item.powerName} on show alarm`}>
        <Checkbox />
      </Tooltip>
      <Tooltip title={`When ${item.powerName} off show alarm`}>
        <Checkbox />
      </Tooltip>
    </Stack>
  );

  return (
    <List sx={{ border: 1, borderColor: "divider", pt: 0, pb: 0 }}>
      <ListItem key="header">
        <HeaderItemRender />
      </ListItem>
      {items.map((item) => (
        <ListItem
          key={item.powerName}
          sx={{ pt: 0, pb: 0, borderTop: 1, borderTopColor: "divider" }}
        >
          <BodyItemRender item={item} />
        </ListItem>
      ))}
    </List>
  );
};

export default AlarmPowerList;
