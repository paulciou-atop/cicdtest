import {
  Box,
  Button,
  Card,
  CardContent,
  Checkbox,
  Divider,
  FormControlLabel,
  List,
  Paper,
  Stack,
  TextField,
  Typography,
} from "@mui/material";
import React, { useState } from "react";
import SessionListItem from "./SessionListItem";

const SessionDetails = ({
  ipRange,
  setIpRange,
  sessionDetails,
  onStartClick,
}) => {
  const [show, setShow] = useState(true);
  const [discoveryStarts, setDiscoveryStarts] = useState(false);

  const handleStartClick = () => {
    setDiscoveryStarts(true);
    onStartClick();
  };

  return (
    <Paper>
      <Card>
        <CardContent>
          <Typography variant="h6" gutterBottom component="div">
            Start Device Scan
          </Typography>
          <Stack spacing={2} direction="row" alignItems="center">
            <TextField
              label="IP Range"
              variant="outlined"
              value={ipRange}
              onChange={(e) => setIpRange(e.target.value)}
              size="small"
            />
            <Button
              variant="contained"
              onClick={handleStartClick}
              disabled={ipRange === ""}
              size="medium"
            >
              Start
            </Button>
            <FormControlLabel
              control={
                <Checkbox checked={show} onChange={() => setShow(!show)} />
              }
              label="Show scan details"
            />
          </Stack>
          <Box sx={{ mt: 1 }}>
            {discoveryStarts && show && (
              <List sx={{ width: "100%", bgcolor: "background.paper" }}>
                <SessionListItem
                  status={sessionDetails.statusStart}
                  message={sessionDetails.messageStart}
                  label="Session ID"
                />
                <Divider variant="inset" component="li" />
                <SessionListItem
                  status={sessionDetails.statusCheck}
                  message={sessionDetails.messageCheck}
                  label="Session Status"
                />
                <Divider variant="inset" component="li" />
                <SessionListItem
                  status={sessionDetails.statusResult}
                  message={sessionDetails.messageResult}
                  label="Session Data Result"
                />
              </List>
            )}
          </Box>
        </CardContent>
      </Card>
    </Paper>
  );
};

export default SessionDetails;
