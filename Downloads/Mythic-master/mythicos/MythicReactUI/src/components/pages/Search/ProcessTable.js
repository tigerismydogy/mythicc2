import React, { useEffect } from 'react';
import {Button, IconButton, Typography, Link} from '@mui/material';
import Table from '@mui/material/Table';
import TableBody from '@mui/material/TableBody';
import TableCell from '@mui/material/TableCell';
import TableContainer from '@mui/material/TableContainer';
import TableHead from '@mui/material/TableHead';
import TableRow from '@mui/material/TableRow';
import Paper from '@mui/material/Paper';
import { tigerDialog, tigerModifyStringDialog, tigerViewJSONAsTableDialog } from '../../tigerComponents/tigerDialog';
import PlaylistAddCheckIcon from '@mui/icons-material/PlaylistAddCheck';
import { gql, useMutation } from '@apollo/client';
import {snackActions} from '../../utilities/Snackbar';
import EditIcon from '@mui/icons-material/Edit';
import tigerStyledTableCell from '../../tigerComponents/tigerTableCell';
import {TagsDisplay, ViewEditTags} from '../../tigerComponents/tigerTag';

const updateFileComment = gql`
mutation updateCommentMutation($tigertree_id: Int!, $comment: String!){
    update_tigertree_by_pk(pk_columns: {id: $tigertree_id}, _set: {comment: $comment}) {
        comment
        id
    }
}
`;

export function ProcessTable(props){
    const [files, setFiles] = React.useState([]);
    useEffect( () => {
        setFiles([...props.processes]);
    }, [props.processes]);
    const onEditComment = ({id, comment}) => {
        const updates = files.map( (file) => {
            if(file.id === id){
                return {...file, comment}
            }else{
                return {...file}
            }
        });
        setFiles(updates);
    }
    return (
        <TableContainer component={Paper} className="tigerElement" >
            <Table stickyHeader size="small" style={{"tableLayout": "fixed", "maxWidth": "100%", "overflow": "scroll"}}>
                <TableHead>
                    <TableRow>
                        <TableCell style={{width: "5rem"}}>Metadata</TableCell>
                        <TableCell style={{width: "5rem"}}> PID </TableCell>
                        <TableCell >Info</TableCell>
                        <TableCell> Name</TableCell>
                        <TableCell style={{width: "15rem"}}>Comment</TableCell>
                        <TableCell style={{width: "10rem"}}>Tags</TableCell>

                    </TableRow>
                </TableHead>
                <TableBody>
                
                {files.map( (op) => (
                    <ProcessTableRow
                        key={"process" + op.id}
                        me={props.me}
                        onEditComment={onEditComment}
                        {...op}
                    />
                ))}
                </TableBody>
            </Table>
        </TableContainer>
    )
}
function ProcessTableRow(props){
    const me = props.me;
    const [viewPermissionsDialogOpen, setViewPermissionsDialogOpen] = React.useState(false);
    const [editCommentDialogOpen, setEditCommentDialogOpen] = React.useState(false);
    const [updateComment] = useMutation(updateFileComment, {
        onCompleted: (data) => {
            snackActions.success("updated comment");
            props.onEditComment(data.update_tigertree_by_pk)
        }
    });
    const onSubmitUpdatedComment = (comment) => {
        updateComment({variables: {tigertree_id: props.id, comment: comment}})
    }
    return (
        <React.Fragment>
            <TableRow hover>
                {viewPermissionsDialogOpen && <tigerDialog fullWidth={true} maxWidth="md" open={viewPermissionsDialogOpen}
                    onClose={()=>{setViewPermissionsDialogOpen(false);}} 
                    innerDialog={<tigerViewJSONAsTableDialog title="View Permissions Data" leftColumn="Permission" rightColumn="Value" value={props.metadata} onClose={()=>{setViewPermissionsDialogOpen(false);}} />}
                    />
                }
                {editCommentDialogOpen && <tigerDialog fullWidth={true} maxWidth="md" open={editCommentDialogOpen}
                    onClose={()=>{setEditCommentDialogOpen(false);}} 
                    innerDialog={<tigerModifyStringDialog title="Edit File Browser Comment" onSubmit={onSubmitUpdatedComment} value={props.comment} onClose={()=>{setEditCommentDialogOpen(false);}} />}
                />
                }
                <tigerStyledTableCell>
                    <Button color="primary" variant="contained" onClick={() => setViewPermissionsDialogOpen(true)}><PlaylistAddCheckIcon /></Button>
                </tigerStyledTableCell>
                <tigerStyledTableCell>
                    <Typography variant="body2" style={{wordBreak: "break-all", textDecoration: props.deleted ? "strike-through" : ""}}>{props.full_path_text}</Typography>
                </tigerStyledTableCell>
                <tigerStyledTableCell>
                    <Typography variant="body2" style={{wordBreak: "break-all"}}><b>Host: </b> {props.host}</Typography>
                    {props.callback?.tigertree_groups.length > 0 ? (
                        <Typography variant="body2" style={{whiteSpace: "pre"}}>
                            <b>Groups: </b>{props?.callback.tigertree_groups.join(", ")}<br/>
                            <b>Callback: </b>{
                                <Link style={{wordBreak: "break-all"}} color="textPrimary" underline="always" target="_blank"
                                      href={"/new/callbacks/" + props.callback.display_id}>
                                    {props.callback.display_id}
                                </Link>
                            }<br/>
                            <b>Task: </b>{
                            <Link style={{wordBreak: "break-all"}} color="textPrimary" underline="always" target="_blank"
                                  href={"/new/task/" + props.task.display_id}>
                                {props.task.display_id}
                            </Link>
                        }
                        </Typography>
                    ) : null}
                </tigerStyledTableCell>

                <tigerStyledTableCell>
                    <Typography variant="body2" style={{wordBreak: "break-all"}}>{props.name_text}</Typography>
                </tigerStyledTableCell>
                <tigerStyledTableCell>
                    <IconButton onClick={() => setEditCommentDialogOpen(true)} size="small" style={{display: "inline-block"}}><EditIcon /></IconButton>
                    <Typography variant="body2" style={{wordBreak: "break-all", display: "inline-block"}}>{props.comment}</Typography>
                    </tigerStyledTableCell>
                <tigerStyledTableCell>
                    <ViewEditTags target_object={"tigertree_id"} target_object_id={props.id} me={me} />
                    <TagsDisplay tags={props.tags} />
                </tigerStyledTableCell>

            </TableRow>
        </React.Fragment>
    )
}

