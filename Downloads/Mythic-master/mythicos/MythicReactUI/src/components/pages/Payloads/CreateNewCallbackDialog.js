import React, {useEffect} from 'react';
import Button from '@mui/material/Button';
import DialogActions from '@mui/material/DialogActions';
import DialogContent from '@mui/material/DialogContent';
import DialogTitle from '@mui/material/DialogTitle';
import { snackActions } from '../../utilities/Snackbar';

import {gql, useMutation} from '@apollo/client';
import { Table, TableBody, TableContainer, TableRow, TableHead, Paper } from '@mui/material';
import tigerTableCell from '../../tigerComponents/tigerTableCell';
import tigerTextField from '../../tigerComponents/tigerTextField';

const createCallback = gql`
mutation createNewCallback($payloadUuid: String!, $callbackConfig: newCallbackConfig) {
  createCallback(payloadUuid: $payloadUuid, newCallback: $callbackConfig){
    status
    error
  }
}
`;
 
export function CreateNewCallbackDialog(props) {
    const [IP, setIP] = React.useState("");
    const [externalIP, setExternalIP] = React.useState("");
    const [host, setHost] = React.useState("");
    const [user, setUser] = React.useState("");
    const [domain, setDomain] = React.useState("");
    const [description, setDescription] = React.useState("");
    const [sleepInfo, setSleepInfo] = React.useState("");
    const [extraInfo, setExtraInfo] = React.useState("");
    const [processName, setProcessName] = React.useState("");
    const onChangeText = (name, value, error) => {
      switch (name) {
        case "IP...":
          setIP(value);
          break;
        case "External IP...":
          setExternalIP(value);
          break;
        case "Host":
          setHost(value);
          break;
        case "User":
          setUser(value);
          break;
        case "Domain":
          setDomain(value);
          break;
        case "Description":
          setDescription(value);
          break;
        case "Sleep Info":
          setSleepInfo(value);
          break;
        case "Extra Info":
          setExtraInfo(value);
          break;
        case "Process Name":
          setProcessName(value);
          break;
      }
    }
    const [createCallbackMutation] = useMutation(createCallback, {
      onCompleted: data => {
        console.log(data);
        if (data.createCallback.status === "success"){
          snackActions.success("Successfully create new callback");
        } else {
          snackActions.error(data.createCallback.error);
        }
        props.onClose();
      },
      onError: error => {
        console.log(error)
        props.onClose();
      }
    })
    const submit = () => {
      createCallbackMutation({variables: {payloadUuid: props.uuid, callbackConfig: {
        ip: IP,
        externalIp: externalIP,
        user: user,
        host: host,
        domain: domain,
        description: description,
        processName: processName,
        sleepInfo: sleepInfo,
        extraInfo: extraInfo
      }}})
      
    }
  return (
    <React.Fragment>
        <DialogTitle id="form-dialog-title">Manually Create Callback for payload {props.filename}</DialogTitle>
        <DialogContent dividers={true}>
          This will generate a new callback based on this payload, but will not trigger a payload execution (there will be no payload running to fetch commands).
          This is useful for webshells that don't reach out to tiger, but still need a callback in order to issue tasking. This is also useful for development and testing purposes.
          <TableContainer className="tigerElement">
            <Table size="small" style={{ "maxWidth": "100%", "overflow": "scroll"}}>
                <TableHead>
                    <TableRow>
                       <tigerTableCell style={{width: "4rem"}}>Attribute</tigerTableCell>
                        <tigerTableCell >Value</tigerTableCell>
                    </TableRow>
                </TableHead>
                <TableBody>
                  <TableRow hover>
                    <tigerTableCell>IP</tigerTableCell>
                    <tigerTableCell>
                      <tigerTextField name="IP..." value={IP} onChange={onChangeText}  />
                    </tigerTableCell>
                  </TableRow>
                  <TableRow hover>
                    <tigerTableCell>External IP</tigerTableCell>
                    <tigerTableCell>
                      <tigerTextField name="External IP..." value={externalIP} onChange={onChangeText}  />
                    </tigerTableCell>
                  </TableRow>
                  <TableRow hover>
                    <tigerTableCell>User</tigerTableCell>
                    <tigerTableCell>
                      <tigerTextField name="User" value={user} onChange={onChangeText}  />
                    </tigerTableCell>
                  </TableRow>
                  <TableRow hover>
                    <tigerTableCell>Host</tigerTableCell>
                    <tigerTableCell>
                      <tigerTextField name="Host" value={host} onChange={onChangeText}  />
                    </tigerTableCell>
                  </TableRow>
                  <TableRow hover>
                    <tigerTableCell>Domain</tigerTableCell>
                    <tigerTableCell>
                      <tigerTextField name="Domain" value={domain} onChange={onChangeText}  />
                    </tigerTableCell>
                  </TableRow>
                  <TableRow hover>
                    <tigerTableCell>Description</tigerTableCell>
                    <tigerTableCell>
                      <tigerTextField name="Description" value={description} onChange={onChangeText}  />
                    </tigerTableCell>
                  </TableRow>
                  <TableRow hover>
                    <tigerTableCell>Process Name</tigerTableCell>
                    <tigerTableCell>
                     <tigerTextField name="Process Name" value={processName} onChange={onChangeText}  />
                    </tigerTableCell>
                  </TableRow>
                  <TableRow hover>
                    <tigerTableCell>Sleep Info</tigerTableCell>
                    <tigerTableCell>
                      <tigerTextField name="Sleep Info" value={sleepInfo} onChange={onChangeText}  />
                    </tigerTableCell>
                  </TableRow>
                  <TableRow hover>
                    <tigerTableCell>Extra Info</tigerTableCell>
                    <tigerTableCell>
                      <tigerTextField name="Extra Info" value={extraInfo} onChange={onChangeText}  />
                    </tigerTableCell>
                  </TableRow>
                </TableBody>
            </Table>
            </TableContainer>
        </DialogContent>
        <DialogActions>
          <Button onClick={props.onClose} variant="contained" color="primary">
            Close
          </Button>
          <Button onClick={submit} variant="contained" color="success">
            Submit
          </Button>
        </DialogActions>
  </React.Fragment>
  );
}

