import React from 'react';
import TableRow from '@mui/material/TableRow';
import Table from '@mui/material/Table';
import TableBody from '@mui/material/TableBody';
import TableHead from '@mui/material/TableHead';
import {useQuery, gql, useMutation} from '@apollo/client';
import FormControl from '@mui/material/FormControl';
import Select from '@mui/material/Select';
import DialogActions from '@mui/material/DialogActions';
import DialogContent from '@mui/material/DialogContent';
import DialogTitle from '@mui/material/DialogTitle';
import {Button, IconButton} from '@mui/material';
import AddCircleIcon from '@mui/icons-material/AddCircle';
import DeleteIcon from '@mui/icons-material/Delete';
import { snackActions } from '../../utilities/Snackbar';
import tigerStyledTableCell from "../../tigerComponents/tigerTableCell";
import tigerTextField from "../../tigerComponents/tigerTextField";
import MenuItem from '@mui/material/MenuItem';
import {tigerDialog} from "../../tigerComponents/tigerDialog";
import {ViewAllCallbacktigerTreeGroupsDialog} from "./ViewCallbacktigerTreeGroupsDialog";
import {tigerStyledTooltip} from "../../tigerComponents/tigerStyledTooltip";
import { useTheme } from '@mui/material/styles';
import LayersIcon from '@mui/icons-material/Layers';


const getCallbacktigerTreeGroups = gql`
query getCallbacktigerTreeGroups($callback_id: Int!) {
  callback_by_pk(id: $callback_id){
    tigertree_groups
  }
  callback {
    tigertree_groups
  }
}
`;
const settigerTreeGroups = gql`
mutation settigerTreeGroups($callback_id: Int!, $tigertree_groups: [String]!){
  update_callback(where: {id: {_eq: $callback_id}}, _set: {tigertree_groups: $tigertree_groups}) {
    affected_rows
  }
}
`;
export function ModifyCallbacktigerTreeGroupsDialog(props){
    const theme = useTheme();
    const [groups, setGroups] = React.useState([]);
    const [otherGroups, setOtherGroups] = React.useState([]);
    const [selectedGroupDropdown, setSelectedGroupDropdown] = React.useState('');
    const [openViewAllCallbacksDialog, setOpenViewAllCallbacksDialog] = React.useState(false);
    const [setCallbackGroups] = useMutation(settigerTreeGroups, {
      onCompleted: data => {
            if(data.update_callback.affected_rows > 0){
                snackActions.success("Successfully updated callback groups.\nPlease close and reopen all process browser and file browser tabs.");
            } else {
                snackActions.error("Failed to update callback groups");
            }
      },
      onError: error => {
        snackActions.error(error.message);
      }
    })
    const { data } = useQuery(getCallbacktigerTreeGroups, {
        fetchPolicy: "no-cache",
        variables: {callback_id: props.callback_id},
        onCompleted: data => {
            setGroups(data.callback_by_pk.tigertree_groups);
            let otherGroupOptions = new Set([]);
            for(let i = 0; i < data.callback.length; i++){
                if(data.callback[i].tigertree_groups.length > 0){
                    data.callback[i].tigertree_groups.forEach( (e) => otherGroupOptions.add(e) );
                }
            }
            otherGroupOptions.delete("Default")
            let otherGroupArray = Array.from(otherGroupOptions).sort();
            otherGroupArray.unshift("Default");
            setOtherGroups(otherGroupArray);
            if( otherGroupArray.length > 0 ){
                setSelectedGroupDropdown(otherGroupArray[0]);
            }
        }
        });
    const submit = (event) => {
        props.onClose(event);
        setCallbackGroups({variables:{callback_id: props.callback_id, tigertree_groups: groups}});
    }
    const addArrayOption = () => {
        const newArray = [...groups, selectedGroupDropdown];
        setGroups(newArray);
    }
    const addNewArrayValue = () => {
        const newArray = [...groups, ""];
        setGroups(newArray);
    }
    const removeArrayValue = (index) => {
        let removed = [...groups];
        removed.splice(index, 1);
        setGroups(removed);
    }
    const onChangeArrayText = (value, error, index) => {
        let values = [...groups];
        if(value.includes("\n")){
            let new_values = value.split("\n");
            values = [...values, ...new_values.slice(1)];
            values[index] = new_values[0];
        }else{
            values[index] = value;
        }
        setGroups(values);
    }
    return (
        <React.Fragment>
          <DialogTitle id="form-dialog-title" style={{display: "flex", justifyContent: "space-between"}}>
              Updating Callback Groups for Callback {props.callback_id}
              <tigerStyledTooltip title="View all groups" >
                  <IconButton size="small" onClick={()=>{setOpenViewAllCallbacksDialog(true);}} style={{color: theme.palette.info.main}} variant="contained"><LayersIcon/></IconButton>
              </tigerStyledTooltip>
          </DialogTitle>
            <div style={{paddingLeft: "30px"}}>
                Group information from this callback and others when looking at the FileBrowser and ProcessBrowser trees. <br/>
                <b>Note:</b> Having <b>no</b> group entries will hide all information from this callback from your FileBrowser and ProcessBrowser views.
            </div>
          <DialogContent dividers={true}>
            <Table size="small" aria-label="details" style={{ "overflowWrap": "break-word"}}>
                <TableHead>
                </TableHead>
                <TableBody>
                    {otherGroups.length > 0 &&
                        <TableRow>
                            <tigerStyledTableCell style={{width: "20%"}}>Add an existing group to this callback</tigerStyledTableCell>
                            <tigerStyledTableCell>
                                <FormControl >
                                    <Select
                                        value={selectedGroupDropdown}
                                        onChange={evt => setSelectedGroupDropdown(evt.target.value)}
                                    >
                                        {
                                            otherGroups.map((opt) => (
                                                <MenuItem key={opt} value={opt}>{opt}</MenuItem>
                                            ))
                                        }
                                    </Select>
                                </FormControl>
                                <IconButton onClick={addArrayOption} size="large"> <AddCircleIcon color="success"  /> </IconButton>
                            </tigerStyledTableCell>
                        </TableRow>
                    }
                    {groups.map( (a, i) => (
                        <TableRow key={'array' + props.name + i} >
                            <tigerStyledTableCell style={{width: "2rem", paddingLeft:"0"}}>
                                <IconButton onClick={(e) => {removeArrayValue(i)}} size="large"><DeleteIcon color="error" /> </IconButton>
                            </tigerStyledTableCell>
                            <tigerStyledTableCell>
                                <tigerTextField required={props.required} fullWidth={true} placeholder={""} value={a} multiline={true} autoFocus={ i > 0}
                                                 onChange={(n,v,e) => onChangeArrayText(v, e, i)} display="inline-block" maxRows={5}
                                />
                            </tigerStyledTableCell>
                        </TableRow>
                    ))}
                    <TableRow >
                        <tigerStyledTableCell style={{width: "5rem", paddingLeft:"0"}}>
                            <IconButton onClick={addNewArrayValue} size="large"> <AddCircleIcon color="success"  /> </IconButton>
                        </tigerStyledTableCell>
                        <tigerStyledTableCell></tigerStyledTableCell>
                    </TableRow>

                </TableBody>
            </Table>
          </DialogContent>
          <DialogActions>
            <Button onClick={props.onClose} variant="contained" color="primary">
              Close
            </Button>
          <Button onClick={submit} variant="contained" color={"success"}>
              Update
          </Button>
        </DialogActions>
            {openViewAllCallbacksDialog &&
                <tigerDialog
                    fullWidth={true}
                    maxWidth={"lg"}
                    open={openViewAllCallbacksDialog}
                    onClose={() => {setOpenViewAllCallbacksDialog(false);}}
                    innerDialog={
                        <ViewAllCallbacktigerTreeGroupsDialog onClose={() => {setOpenViewAllCallbacksDialog(false);}} />
                    }
                />
            }
        </React.Fragment>
        )
}

