import React from 'react';
import {Button} from '@mui/material';
import Table from '@mui/material/Table';
import TableBody from '@mui/material/TableBody';
import TableCell from '@mui/material/TableCell';
import TableContainer from '@mui/material/TableContainer';
import TableHead from '@mui/material/TableHead';
import TableRow from '@mui/material/TableRow';
import Paper from '@mui/material/Paper';
import { OperationTableRow } from './OperationTableRow';
import Typography from '@mui/material/Typography';
import AddCircleOutlineOutlinedIcon from '@mui/icons-material/AddCircleOutlineOutlined';
import { tigerDialog } from '../../tigerComponents/tigerDialog';
import {useTheme} from '@mui/material/styles';
import {SettingsOperatorDialog} from '../Settings/SettingsOperatorDialog';
import {snackActions} from '../../utilities/Snackbar';
import {useMutation, gql} from '@apollo/client';
import {tigerModifyStringDialog} from '../../tigerComponents/tigerDialog';
import VisibilityOffIcon from '@mui/icons-material/VisibilityOff';
import VisibilityIcon from '@mui/icons-material/Visibility';
import { tigerStyledTooltip } from '../../tigerComponents/tigerStyledTooltip';
import { IconButton } from '@mui/material';

const newOperatorMutation = gql`
mutation NewOperator($username: String!, $password: String!) {
  createOperator(input: {password: $password, username: $username}) {
    error
    id
    status
  }
}
`;
const Update_Operation = gql`
mutation MyMutation($operation_id: Int!, $channel: String!, $complete: Boolean!, $name: String!, $webhook: String!) {
    updateOperation(operation_id: $operation_id, channel: $channel, complete: $complete, name: $name, webhook: $webhook) {
        status
        error
        name
        complete
        id
    }
}
`;
const newOperationMutation = gql`
mutation newOperationMutation($name: String){
    createOperation(name: $name){
        status
        error
        operation_id
        operation_name
        
    }
}
`;

export function OperationTable(props){
    const theme = useTheme();
    const [openNewOperator, setOpenNewOperatorDialog] = React.useState(false);
    const [openNewOperation, setOpenNewOperationDialog] = React.useState(false);
    const [showDeleted, setShowDeleted] = React.useState(false);
    const [newOperator] = useMutation(newOperatorMutation, {
        update: (cache, {data}) => {
            if(data.createOperator.status === "success"){
                snackActions.success("Created operator");
            }else{
                snackActions.error(data.createOperator.error);
            }
        },
        onError: (err) => {
          snackActions.warning("Unable to create new operator - Access Denied");
          console.log(err);
        }
    });
    const [updateOperation] = useMutation(Update_Operation, {
        onCompleted: (data) => {
            if(data.updateOperation.status === "success"){
                props.onUpdateOperation(data.updateOperation);
                snackActions.success("Successfully updated operation");
            } else {
                snackActions.error(data.updateOperation.error);
            }

        },
        onError: (data) => {
          snackActions.error("Failed to update operation");
          console.log("error updating operation", data);
        }
      });
    const [newOperation] = useMutation(newOperationMutation, {
        onCompleted: (data) => {
            //console.log(data);
            if(data.createOperation.status === "success"){
                snackActions.success("Successfully created operation!");
                props.onNewOperation({name: data.createOperation.operation_name, id: data.createOperation.operation_id});
            }else{
                snackActions.error(data.createOperation.error);
            }
        },
        onError: (data) => {
            snackActions.error("Unable to create new operation - Access Denied")
            console.log(data);
        }
    })
    const onUpdateOperation = ({operation_id, name, channel, webhook, complete}) => {
        updateOperation({variables:{
            operation_id,
            name,
            channel,
            webhook,
            complete
        }});
    }
    const onSubmitNewOperator = (id, username, passwordOld, passwordNew) => {
        if(passwordOld !== passwordNew){
            snackActions.error("Passwords don't match");
        }else if(passwordNew.length === 0){
            snackActions.error("Password must not be empty",);
        }else if(username.length === 0) {
            snackActions.error("Username must not be empty",);
        } else if(passwordNew.length < 12){
            snackActions.error("Password must be at least 12 characters long");
        }else{
            newOperator({variables:{username:username, password:passwordNew}})
            setOpenNewOperatorDialog(false);
        }
    }
    const onSubmitNewOperation = (operation_name) => {
        newOperation({variables: {name: operation_name}})
    }
    return (
        <React.Fragment>
        <Paper elevation={5} style={{backgroundColor: theme.pageHeader.main, color: theme.pageHeaderText.main, marginBottom: "5px", marginTop: "10px", marginRight: "5px"}} variant={"elevation"}>
            <Typography variant="h3" style={{textAlign: "left", display: "inline-block", marginLeft: "20px"}}>
                Operations
            </Typography>
            {showDeleted ? (
                <tigerStyledTooltip title={"Hide Deleted Operations"} style={{float: "right"}}>
                    <IconButton size="small" style={{float: "right", marginTop: "5px"}} variant="contained" onClick={() => setShowDeleted(!showDeleted)}><VisibilityIcon /></IconButton>
                </tigerStyledTooltip>
                
              ) : (
                <tigerStyledTooltip title={"Show Deleted Operations"} style={{float: "right"}}>
                  <IconButton size="small" style={{float: "right",  marginTop: "5px"}} variant="contained" onClick={() => setShowDeleted(!showDeleted)} ><VisibilityOffIcon /></IconButton>
                </tigerStyledTooltip>
              )}
            <Button size="small" onClick={() => {setOpenNewOperationDialog(true);}} style={{marginRight: "20px", float: "right", marginTop: "10px"}} startIcon={<AddCircleOutlineOutlinedIcon/>} color="success" variant="contained">New Operation</Button>
            <Button size="small" onClick={()=>{setOpenNewOperatorDialog(true);}} style={{marginRight: "20px", float: "right", marginTop: "10px"}} startIcon={<AddCircleOutlineOutlinedIcon/>} color="success" variant="contained">New Operator</Button>
            
            {openNewOperator &&
                <tigerDialog open={openNewOperator} 
                    onClose={()=>{setOpenNewOperatorDialog(false);}} 
                    innerDialog={<SettingsOperatorDialog title="New Operator" onAccept={onSubmitNewOperator} handleClose={()=>{setOpenNewOperatorDialog(false);}}  {...props}/>}
                />
            }
            {openNewOperation &&
                <tigerDialog 
                    fullWidth={true} 
                    open={openNewOperation}  
                    onClose={() => {setOpenNewOperationDialog(false);}}
                    innerDialog={
                        <tigerModifyStringDialog title={"New Operation's Name"} 
                            onClose={() => {setOpenNewOperationDialog(false);}} 
                            value={""} 
                            onSubmit={onSubmitNewOperation} 
                        />
                    }
                />
            }
            
        </Paper>
        <TableContainer component={Paper} className="tigerElement">   
            
            <Table  size="small" style={{"tableLayout": "fixed", "maxWidth": "calc(100vw)", "overflow": "scroll"}}>
                <TableHead>
                    <TableRow>
                        <TableCell style={{width: "8rem"}}></TableCell>
                        <TableCell style={{width: "8rem"}}>Configure</TableCell>
                        <TableCell style={{width: "8rem"}}>Operators</TableCell>
                        <TableCell>Operation Name</TableCell>
                        <TableCell>Operation Admin</TableCell>
                        <TableCell style={{width: "10rem"}}>Analysis</TableCell>
                        <TableCell style={{width: "12rem"}}>Operation Status</TableCell>
                    </TableRow>
                </TableHead>
                <TableBody>
                
                {props.operations.map( (op) => (
                    showDeleted || !op.deleted ? (
                        <OperationTableRow
                            me={props.me}
                            key={"operation" + op.id} onUpdateOperation={onUpdateOperation}
                            updateDeleted={props.updateDeleted}
                            {...op} operator={props.operator}
                        />
                    ) : (null)
                ))}
                </TableBody>
            </Table>
        </TableContainer>
    </React.Fragment>
    )
}

