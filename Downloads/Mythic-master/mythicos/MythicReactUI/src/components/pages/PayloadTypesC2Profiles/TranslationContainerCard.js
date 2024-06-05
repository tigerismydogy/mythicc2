import React from 'react';
import { styled } from '@mui/material/styles';
import Card from '@mui/material/Card';
import Typography from '@mui/material/Typography';
import {FontAwesomeIcon} from '@fortawesome/react-fontawesome';
import { faLanguage } from '@fortawesome/free-solid-svg-icons';
import {useTheme} from '@mui/material/styles';
import IconButton from '@mui/material/IconButton';
import {useMutation, gql} from '@apollo/client';
import {snackActions} from '../../utilities/Snackbar';
import MenuBookIcon from '@mui/icons-material/MenuBook';
import {tigerConfirmDialog} from '../../tigerComponents/tigerConfirmDialog';
import DeleteIcon from '@mui/icons-material/Delete';
import RestoreFromTrashOutlinedIcon from '@mui/icons-material/RestoreFromTrashOutlined';
import TableRow from '@mui/material/TableRow';
import tigerTableCell from "../../tigerComponents/tigerTableCell";
import {tigerStyledTooltip} from "../../tigerComponents/tigerStyledTooltip";

const PREFIX = 'TranslationContainerCard';

const classes = {
  root: `${PREFIX}-root`,
  expand: `${PREFIX}-expand`,
  expandOpen: `${PREFIX}-expandOpen`,
  running: `${PREFIX}-running`,
  notrunning: `${PREFIX}-notrunning`
};


const toggleDeleteStatus = gql`
mutation toggleC2ProfileDeleteStatus($translationcontainer_id: Int!, $deleted: Boolean!){
  update_translationcontainer_by_pk(pk_columns: {id: $translationcontainer_id}, _set: {deleted: $deleted}) {
    id
  }
}
`;

export function TranslationContainerRow({service, showDeleted}) {
  const theme = useTheme();

  const [openDelete, setOpenDeleteDialog] = React.useState(false);
  const [updateDeleted] = useMutation(toggleDeleteStatus, {
      onCompleted: data => {
      },
      onError: error => {
          if(service.deleted){
              snackActions.error("Failed to restore translation profile");
          } else {
              snackActions.error("Failed to mark translation profile as deleted");
          }
      }
    });
    const onAcceptDelete = () => {
      updateDeleted({variables: {translationcontainer_id: service.id, deleted: !service.deleted}})
      setOpenDeleteDialog(false);
    }
  return (

        <TableRow hover>
            <tigerTableCell>
                {service.deleted ? (
                    <IconButton size="small" onClick={()=>{setOpenDeleteDialog(true);}} color="success" variant="contained"><RestoreFromTrashOutlinedIcon/></IconButton>
                ) : (
                    <IconButton size="small" onClick={()=>{setOpenDeleteDialog(true);}} color="error" variant="contained"><DeleteIcon/></IconButton>
                )}
                {openDelete &&
                    <tigerConfirmDialog onClose={() => {setOpenDeleteDialog(false);}} onSubmit={onAcceptDelete}
                                         open={openDelete}
                                         acceptText={service.deleted ? "Restore" : "Remove"}
                                         acceptColor={service.deleted ? "success": "error"} />
                }
            </tigerTableCell>
            <tigerTableCell>
                <FontAwesomeIcon icon={faLanguage} style={{width: "80px", height: "80px"}} />
            </tigerTableCell>
            <tigerTableCell>
                {service.name}
            </tigerTableCell>
            <tigerTableCell>
                Translation
            </tigerTableCell>
            <tigerTableCell>
                <Typography variant="body1" component="p">
                    <b>Author:</b> {service.author}
                </Typography>
                <Typography variant="body1" component="p">
                    <b>Supported Agents:</b> {service.payloadtypes.filter(pt => !pt.deleted).map( (pt) => pt.name).join(", ")}
                </Typography>
                <Typography variant="body2" component="p">
                    <b>Description: </b>{service.description}
                </Typography>
            </tigerTableCell>
            <tigerTableCell>
                <Typography variant="body2" component="p" color={service.container_running ? theme.palette.success.main : theme.palette.error.main} >
                    <b>{service.container_running ? "Online" : "Offline"}</b>
                </Typography>
            </tigerTableCell>
            <tigerTableCell>
                <tigerStyledTooltip title={"Documentation"}>
                    <IconButton
                        color={"secondary"}
                        href={"/docs/c2-profiles/" + service.name.toLowerCase()}
                        target="_blank"
                        size="large">
                        <MenuBookIcon />
                    </IconButton>
                </tigerStyledTooltip>
            </tigerTableCell>
        </TableRow>

  );
}
