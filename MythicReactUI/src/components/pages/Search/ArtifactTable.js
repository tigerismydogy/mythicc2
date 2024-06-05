import React, { useEffect } from 'react';
import {Typography, Link} from '@mui/material';
import Table from '@mui/material/Table';
import TableBody from '@mui/material/TableBody';
import TableCell from '@mui/material/TableCell';
import TableContainer from '@mui/material/TableContainer';
import TableHead from '@mui/material/TableHead';
import TableRow from '@mui/material/TableRow';
import Paper from '@mui/material/Paper';
import tigerStyledTableCell from '../../tigerComponents/tigerTableCell';


export function ArtifactTable(props){
    const [artifacts, setArtifacts] = React.useState([]);
    useEffect( () => {
        setArtifacts([...props.artifacts]);
    }, [props.artifacts]);

    return (
        <TableContainer component={Paper} className="tigerElement" style={{}}>
            <Table stickyHeader size="small" style={{}}>
                <TableHead>
                    <TableRow>
                        <TableCell >Type</TableCell>
                        <TableCell >Command</TableCell>
                        <TableCell >Task</TableCell>
                        <TableCell >Callback</TableCell>
                        <TableCell >Operator</TableCell>
                        <TableCell >Host</TableCell>
                        <TableCell >Artifact</TableCell>
                    </TableRow>
                </TableHead>
                <TableBody>
                
                {artifacts.map( (op) => (
                    <ArtifactTableRow
                        key={"cred" + op.id}
                        {...op}
                    />
                ))}
                </TableBody>
            </Table>
        </TableContainer>
    )
}

function ArtifactTableRow(props){
    return (
        <React.Fragment>
            <TableRow hover>
                <tigerStyledTableCell>
                    <Typography variant="body2" >{props.base_artifact}</Typography>
                </tigerStyledTableCell>
                <tigerStyledTableCell >
                    <Typography variant="body2" >{props.task.command.cmd}</Typography>
                </tigerStyledTableCell>
                <tigerStyledTableCell>
                    <Link style={{wordBreak: "break-all"}} color="textPrimary" underline="always" target="_blank" 
                        href={"/new/task/" + props.task.display_id}>
                            {props.task.display_id}
                    </Link>
                </tigerStyledTableCell>
                <tigerStyledTableCell style={{wordBreak: "break-all"}}>
                    <Link style={{wordBreak: "break-all"}} color="textPrimary" underline="always" target="_blank" 
                        href={"/new/callbacks/" + props.task.callback.display_id}>
                            {props.task.callback.display_id}
                    </Link>
                    {props.task?.callback?.tigertree_groups.length > 0 ? (
                        <Typography variant="body2" style={{whiteSpace: "pre"}}>
                            <b>Groups: </b>{"\n" + props?.task?.callback.tigertree_groups.join("\n")}
                        </Typography>
                    ) : null}
                </tigerStyledTableCell>
                <tigerStyledTableCell>
                <Typography variant="body2" style={{ display: "inline-block"}}>{props?.task?.operator?.username || null}</Typography>
                </tigerStyledTableCell>
                <tigerStyledTableCell >
                    <Typography variant="body2" style={{ display: "inline-block"}}>{props.host}</Typography>
                </tigerStyledTableCell>
                <tigerStyledTableCell>
                <Typography variant="body2" style={{wordBreak: "break-all", display: "inline-block"}}>{props.artifact_text}</Typography>
                </tigerStyledTableCell>
              
            </TableRow>
        </React.Fragment>
    )
}

