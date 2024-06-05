import React from 'react';
import {Button} from '@mui/material';
import TableCell from '@mui/material/TableCell';
import TableRow from '@mui/material/TableRow';
import Switch from '@mui/material/Switch';
import Box from '@mui/material/Box';
import Collapse from '@mui/material/Collapse';
import IconButton from '@mui/material/IconButton';
import Table from '@mui/material/Table';
import TableBody from '@mui/material/TableBody';
import TableHead from '@mui/material/TableHead';
import Typography from '@mui/material/Typography';
import KeyboardArrowDownIcon from '@mui/icons-material/KeyboardArrowDown';
import KeyboardArrowUpIcon from '@mui/icons-material/KeyboardArrowUp';
import DeleteIcon from '@mui/icons-material/Delete';
import AddCircleOutlineOutlinedIcon from '@mui/icons-material/AddCircleOutlineOutlined';
import { APITokenRow } from './SettingsOperatorAPIToken';
import { SettingsOperatorDialog } from './SettingsOperatorDialog';
import { SettingsOperatorDeleteDialog } from './SettingsOperatorDeleteDialog';
import { tigerDialog } from '../../tigerComponents/tigerDialog';
import { toLocalTime } from '../../utilities/Time';
import tigerStyledTableCell from '../../tigerComponents/tigerTableCell';
import {SettingsOperatorUIConfigDialog} from './SettingsOperatorUIConfigDialog';
import SettingsIcon from '@mui/icons-material/Settings';
import { snackActions } from '../../utilities/Snackbar';
import RestoreFromTrashIcon from '@mui/icons-material/RestoreFromTrash';
import PasswordIcon from '@mui/icons-material/Password';
import {SettingsOperatorExperimentalUIConfigDialog} from "./SettingsOperatorExperimentalUIConfigDialog";
import RocketLaunchIcon from '@mui/icons-material/RocketLaunch';
import {tigerStyledTooltip} from "../../tigerComponents/tigerStyledTooltip";
import VpnKeyIcon from '@mui/icons-material/VpnKey';
import {SettingsOperatorSecretsConfigDialog} from "./SettingsOperatorSecretsConfigDialog";

export function SettingsOperatorTableRow(props){
    const [open, setOpen] = React.useState(false);
    const [openUpdate, setOpenUpdateDialog] = React.useState(false);
    const [openDelete, setOpenDeleteDialog] = React.useState(false);
    const [openUIConfig, setOpenUIConfig] = React.useState(false);
    const [openExperimentalUIConfig, setOpenExperimentalUIConfig] = React.useState(false);
    const [openSecretsConfig, setOpenSecretsConfig] = React.useState(false);
    const me = props.me;
    const isMe = ( me?.user?.user_id || 0 ) === props.id;
    const onViewUTCChanged = (evt) => {
        const {id} = props;
        props.onViewUTCChanged(id, !props[evt.target.name]);
    }
    const onAdminChanged = (evt) => {
        const {id} = props;
        props.onAdminChanged(id, !props[evt.target.name]);
    }
    const onActiveChanged = (evt) => {
        const {id} = props;
        props.onActiveChanged(id, !props[evt.target.name]);
    }
    const onAccept = (id, username, passwordOld, passwordNew) => {
        if(username !== props.username){
          props.onUsernameChanged(id, username);
        }
        if(passwordNew.length > 0){
          props.onPasswordChanged({user_id: id, old_password: passwordOld, new_password: passwordNew});
          
        } else if (passwordOld.length > 0) {
          snackActions.warning("Old password set but not new password");
          return;
        }
        setOpenUpdateDialog(false);
    }
    const onAcceptDelete = (id) => {
        props.onDeleteOperator(id, !props.deleted);
        setOpenDeleteDialog(false);
    }
    return (
        <React.Fragment>
            <TableRow key={props.id}>
                <tigerStyledTableCell >
                    <IconButton size="large" onClick={()=>{setOpenDeleteDialog(true);}}
                              disabled={(isMe || !props.userIsAdmin)} color={props.deleted ? "success": "error"}
                              variant="contained">
                        {props.deleted ? <RestoreFromTrashIcon /> : <DeleteIcon/>}
                    </IconButton>
                  {openDelete && 
                      <tigerDialog open={openDelete} 
                      onClose={()=>{setOpenDeleteDialog(false);}} 
                      innerDialog={<SettingsOperatorDeleteDialog onClose={()=>{setOpenDeleteDialog(false);}}  onAccept={onAcceptDelete} {...props} />}
                  />
                  }
                  
                </tigerStyledTableCell>
                <tigerStyledTableCell>{props.username}</tigerStyledTableCell>
                <tigerStyledTableCell>
                    <IconButton size="large"
                            disabled={!(isMe || props.userIsAdmin)}
                           onClick={()=>{setOpenUpdateDialog(true);}}
                           color="error" ><PasswordIcon /></IconButton>
                  {openUpdate &&
                    <tigerDialog open={openUpdate} 
                     onClose={()=>{setOpenUpdateDialog(false);}} 
                    innerDialog={<SettingsOperatorDialog onAccept={onAccept} handleClose={()=>{setOpenUpdateDialog(false);}} title="Update Operator Username/Password"  {...props}/>}
                />
                  }
                    
                </tigerStyledTableCell>
                <tigerStyledTableCell>
                    <Switch
                        checked={props.view_utc_time}
                        disabled={!isMe}
                        onChange={onViewUTCChanged}
                        color="info"
                        inputProps={{ 'aria-label': 'primary checkbox' }}
                        name="view_utc_time"
                      />
                </tigerStyledTableCell>
                <tigerStyledTableCell>
                  {isMe && 
                      <>
                        <IconButton style={{display: "inline-block"}} size="large" onClick={()=>{setOpenUIConfig(true);}} color="info" variant='contained'>
                            <SettingsIcon />
                        </IconButton>
                        {openUIConfig &&
                          <tigerDialog open={openUIConfig} onClose={()=>{setOpenUIConfig(false)}} maxWidth={"md"} fullWidth
                          innerDialog={<SettingsOperatorUIConfigDialog  onClose={()=>{setOpenUIConfig(false);}} {...props} />}
                          />
                        }
                          <tigerStyledTooltip title={"Secrets"} >
                              <IconButton size="large" onClick={()=>{setOpenSecretsConfig(true);}} color="error" variant='contained'>
                                  <VpnKeyIcon />
                              </IconButton>
                          </tigerStyledTooltip>
                          {openSecretsConfig &&
                              <tigerDialog open={openSecretsConfig} onClose={()=>{setOpenSecretsConfig(false)}} maxWidth={"xl"} fullWidth
                                            innerDialog={<SettingsOperatorSecretsConfigDialog  onClose={()=>{setOpenSecretsConfig(false);}} {...props} />}
                              />
                          }
                        <tigerStyledTooltip title={"Experimental UI Settings"} >
                            <IconButton size="large" onClick={()=>{setOpenExperimentalUIConfig(true);}} color="warning" variant='contained'>
                                <RocketLaunchIcon />
                            </IconButton>
                        </tigerStyledTooltip>
                          {openExperimentalUIConfig &&
                              <tigerDialog open={openExperimentalUIConfig} onClose={()=>{setOpenExperimentalUIConfig(false)}} maxWidth={"md"} fullWidth
                                            innerDialog={<SettingsOperatorExperimentalUIConfigDialog  onClose={()=>{setOpenExperimentalUIConfig(false);}} {...props} />}
                              />
                          }
                      </>
                  }
                </tigerStyledTableCell>
                <tigerStyledTableCell>
                    <Switch
                        color={ isMe || !props.userIsAdmin ? "secondary" : "info"}
                        checked={props.active}
                        disabled={isMe || !props.userIsAdmin}
                        onChange={onActiveChanged}
                        inputProps={{ 'aria-label': 'primary checkbox' }}
                        name="active"
                      />
                </tigerStyledTableCell>
                <tigerStyledTableCell>{toLocalTime(props.last_login, me?.user?.view_utc_time )}</tigerStyledTableCell>
                <tigerStyledTableCell>{toLocalTime(props.creation_time, me?.user?.view_utc_time )}</tigerStyledTableCell>
                <tigerStyledTableCell>
                    <Switch
                        color={ isMe || !props.userIsAdmin ? "secondary" : "info"}
                        checked={props.admin}
                        disabled={isMe || !props.userIsAdmin}
                        onChange={onAdminChanged}
                        inputProps={{ 'aria-label': 'primary checkbox' }}
                        name="admin"
                      />
                </tigerStyledTableCell>
                <tigerStyledTableCell>
                  {props.id === me.user.id && 
                    <IconButton aria-label="expand row" size="small" onClick={() => setOpen(!open)}>
                      {open ? <KeyboardArrowUpIcon className="tigerElement"/> : <KeyboardArrowDownIcon className="tigerElement"/>}
                    </IconButton>
                  }
                    
                </tigerStyledTableCell>
            </TableRow>
            <TableRow>
              {props.id === me.user.id &&
                <tigerStyledTableCell style={{ paddingBottom: 0, paddingTop: 0 }} colSpan={10}>
                  <Collapse in={open} timeout="auto" unmountOnExit>
                    <Box margin={1}>
                      <Typography variant="h6" gutterBottom component="div" style={{display: "inline-block"}}>
                        API Tokens
                      </Typography>
                      <Button size="small" onClick={props.onCreateAPIToken} style={{float: "right"}} startIcon={<AddCircleOutlineOutlinedIcon/>} color="success" variant="contained">New</Button>
                      <Table size="small" aria-label="tokens" style={{"tableLayout": "fixed", "overflowWrap": "break-word"}}>
                        <TableHead>
                          <TableRow>
                            <TableCell style={{width: "10rem"}}>Delete</TableCell>
                            <TableCell>Token</TableCell>
                          </TableRow>
                        </TableHead>
                        <TableBody>
                          {
                            props.apitokens.map((token) => (<APITokenRow {...token} key={"token" + token.id} onDeleteAPIToken={props.onDeleteAPIToken} />))
                          }
                        </TableBody>
                      </Table>
                    </Box>
                  </Collapse>
                </tigerStyledTableCell>
              }
            
          </TableRow>
        </React.Fragment>
        )
}

