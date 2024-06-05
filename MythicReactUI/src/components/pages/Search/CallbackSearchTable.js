import React, { useEffect } from 'react';
import {IconButton, Typography, Link} from '@mui/material';
import { tigerDialog } from '../../tigerComponents/tigerDialog';
import Table from '@mui/material/Table';
import TableBody from '@mui/material/TableBody';
import TableCell from '@mui/material/TableCell';
import TableContainer from '@mui/material/TableContainer';
import TableHead from '@mui/material/TableHead';
import TableRow from '@mui/material/TableRow';
import Paper from '@mui/material/Paper';
import {tigerConfirmDialog} from '../../tigerComponents/tigerConfirmDialog';
import {  useMutation } from '@apollo/client';
import {snackActions} from '../../utilities/Snackbar';
import {useTheme} from '@mui/material/styles';
import VisibilityOffIcon from '@mui/icons-material/VisibilityOff';
import VisibilityIcon from '@mui/icons-material/Visibility';
import {toggleHideCallbackMutations} from '../Callbacks/CallbackMutations';
import { tigerStyledTooltip } from '../../tigerComponents/tigerStyledTooltip';
import {DetailedCallbackTable} from '../Callbacks/DetailedCallbackTable';
import InfoIcon from '@mui/icons-material/Info';
import tigerStyledTableCell from '../../tigerComponents/tigerTableCell';



export function CallbackSearchTable(props){
    const [callbacks, setCallbacks] = React.useState([]);
    useEffect( () => {
        setCallbacks([...props.callbacks]);
    }, [props.callbacks]);

    const onEditDeleted = ({id, active}) => {
        const updates = callbacks.map( (cred) => {
            if(cred.id === id){
                return {...cred, active}
            }else{
                return {...cred}
            }
        });
        setCallbacks(updates);
    }

    return (
        <TableContainer component={Paper} className="tigerElement">
            <Table stickyHeader size="small" style={{"maxWidth": "100%", "overflow": "auto"}}>
                <TableHead>
                    <TableRow>
                        <TableCell style={{width: "2rem"}}>View</TableCell>
                        <TableCell >User</TableCell>
                        <TableCell >Domain</TableCell>
                        <TableCell >Host</TableCell>
                        <TableCell >Groups</TableCell>
                        <TableCell >Description</TableCell>
                        <TableCell >IP</TableCell>
                        <TableCell >ID</TableCell>
                        <TableCell style={{width: "40px"}}>Agent</TableCell>
                        <TableCell style={{width: "2rem"}}>Details</TableCell>
                    </TableRow>
                </TableHead>
                <TableBody>
                
                {callbacks.map( (op) => (
                    <CallbackSearchTableRow
                        key={"cred" + op.id}
                        onEditDeleted={onEditDeleted}
                        {...op}
                    />
                ))}
                </TableBody>
            </Table>
        </TableContainer>
    )
}

function CallbackSearchTableRow(props){
    const theme = useTheme();
    const [openDeleteDialog, setOpenDeleteDialog] = React.useState(false);
    const [openMetaDialog, setOpenMetaDialog] = React.useState(false);
    const [updateDeleted] = useMutation(toggleHideCallbackMutations, {
        onCompleted: (data) => {
            snackActions.success("Updated active status");
            props.onEditDeleted({id: props.id, active: !props.active});
        },
        onError: (data) => {
            snackActions.error("Operation not allowed");
        }
    });
    const ips = JSON.parse(props.ip);
    const onAcceptDelete = () => {
        updateDeleted({variables: {callback_display_id: props.display_id, active: !props.active}})
    }
    return (
        <React.Fragment>
            <TableRow hover>
                <tigerConfirmDialog onClose={() => {setOpenDeleteDialog(false);}} onSubmit={onAcceptDelete} open={openDeleteDialog} acceptText={props.active ? "Hide" : "Restore" }/>
                
                <tigerStyledTableCell>{!props.active ? (
                    <tigerStyledTooltip title="Restore Callback for Tasking">
                        <IconButton size="small" onClick={()=>{setOpenDeleteDialog(true);}} style={{color: theme.palette.error.main}} variant="contained"><VisibilityOffIcon/></IconButton>
                    </tigerStyledTooltip>
                ) : (
                    <tigerStyledTooltip title="Hide Callback so it can't be used in Tasking">
                        <IconButton size="small" onClick={()=>{setOpenDeleteDialog(true);}} style={{color: theme.palette.success.main}} variant="contained"><VisibilityIcon/></IconButton>
                    </tigerStyledTooltip>
                )} </tigerStyledTableCell>
                <tigerStyledTableCell>
                    <Typography variant="body2" style={{wordBreak: "break-all"}}>{props.user}</Typography>
                </tigerStyledTableCell>
                <tigerStyledTableCell >
                    <Typography variant="body2" style={{wordBreak: "break-all"}}>{props.domain}</Typography>
                </tigerStyledTableCell>
                <tigerStyledTableCell>{props.host}</tigerStyledTableCell>
                <tigerStyledTableCell style={{whiteSpace: "pre"}}>{props.tigertree_groups.join("\n")}</tigerStyledTableCell>
                <tigerStyledTableCell >
                    <Typography variant="body2" style={{wordBreak: "break-all", display: "inline-block"}}>{props.description}</Typography>
                </tigerStyledTableCell>
                <tigerStyledTableCell style={{whiteSpace: "pre"}}>
                    {ips.slice(0,5).join("\n")}
                    {ips.length > 5 ? "\n..." : null}
                </tigerStyledTableCell>
                <tigerStyledTableCell>
                <Link style={{wordBreak: "break-all"}} color="textPrimary" underline="always" target="_blank" 
                        href={"/new/callbacks/" + props.display_id}>
                            {props.display_id}
                    </Link>
                </tigerStyledTableCell>
                <tigerStyledTableCell>
                <tigerStyledTooltip title={props.payload.payloadtype.name}>
                    <img
                        style={{width: "35px", height: "35px"}}
                        src={"/static/" + props.payload.payloadtype.name + ".svg"}
                    />
                </tigerStyledTooltip>
                </tigerStyledTableCell>
                <tigerStyledTableCell>
                    <InfoIcon onClick={() => setOpenMetaDialog(true)} style={{color: theme.palette.info.main, cursor: "pointer"}}/>
                    {openMetaDialog && 
                        <tigerDialog fullWidth={true} maxWidth="lg" open={openMetaDialog}
                            onClose={()=>{setOpenMetaDialog(false);}} 
                            innerDialog={<DetailedCallbackTable onClose={()=>{setOpenMetaDialog(false);}} callback_id={props.id} />}
                        />
                    }
                </tigerStyledTableCell>
            </TableRow>
        </React.Fragment>
    )
}

