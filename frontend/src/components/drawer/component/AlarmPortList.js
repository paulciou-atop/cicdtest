import {
  Checkbox,
  List,
  ListItem,
  Stack,
  Tooltip,
  Typography,
} from "@mui/material";
import React from "react";

const AlarmPortList = ({ items }) => {
  const HeaderItemRender = () => (
    <Stack
      direction="row"
      justifyContent="space-between"
      alignItems="center"
      flexGrow={1}
    >
      <Typography variant="body1">Name</Typography>
      <Typography variant="body1">LinkUp</Typography>
      <Typography variant="body1">LinkDown</Typography>
    </Stack>
  );
  const BodyItemRender = ({ item }) => (
    <Stack
      direction="row"
      justifyContent="space-between"
      alignItems="center"
      flexGrow={1}
    >
      <Typography variant="body1">{`${item.portName}`}</Typography>
      <Tooltip
        title={`When ${item.portName} link up show alarm.`}
        enterDelay={0.3}
      >
        <Checkbox />
      </Tooltip>
      <Tooltip
        title={`When ${item.portName} link down show alarm.`}
        enterDelay={0.3}
      >
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
          key={item.portName}
          sx={{ pt: 0, pb: 0, borderTop: 1, borderTopColor: "divider" }}
        >
          <BodyItemRender item={item} />
        </ListItem>
      ))}
    </List>
  );
};

export default AlarmPortList;
