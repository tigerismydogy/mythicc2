import React from 'react';
import {Box, Button} from '@mui/material';
import TableRow from '@mui/material/TableRow';
import IconButton from '@mui/material/IconButton';
import DeleteIcon from '@mui/icons-material/Delete';
import { tigerDialog } from '../../tigerComponents/tigerDialog';
import {Typography} from '@mui/material';
import {tigerConfirmDialog} from '../../tigerComponents/tigerConfirmDialog';
import { useTheme, adaptV4Theme } from '@mui/material/styles';
import { tigerStyledTooltip } from '../../tigerComponents/tigerStyledTooltip';
import tigerStyledTableCell from '../../tigerComponents/tigerTableCell';
import { createTheme } from '@mui/material/styles';
import {NewTagtypesDialog} from './NewTagtypesDialog';
import Chip from '@mui/material/Chip';
import SettingsIcon from '@mui/icons-material/Settings';


export function TagtypesTableRow(props){
    const theme = useTheme();
    const [openUpdate, setOpenUpdateDialog] = React.useState(false);
    const [openDelete, setOpenDeleteDialog] = React.useState(false);
    const [lightColor, setLightColor] = React.useState(theme.palette.text.primary)
    const [darkColor, setDarkColor] = React.useState(theme.palette.text.primary);
    
    const onAcceptDelete = () => {
        props.onDeleteTagtype(props.id);
        setOpenDeleteDialog(false);
    }
    React.useEffect( () => {
      let lightTheme = createTheme(adaptV4Theme({palette: {mode: "light",}}));
      let darkTheme = createTheme(adaptV4Theme({palette: {mode: "dark",}}));
      setLightColor(lightTheme.palette.text.primary);
      setDarkColor(darkTheme.palette.text.primary);
    }, [])
    return (
      
        <React.Fragment>
            <TableRow key={"payload" + props.id} hover>
                <tigerStyledTableCell>

                  <tigerStyledTooltip title={"Delete the tag type and all associated tags"}>
                    <IconButton size="small" onClick={()=>{setOpenDeleteDialog(true);}} color="error" variant="contained"><DeleteIcon/></IconButton>
                  </tigerStyledTooltip>
                  
                  {openDelete && 
                    <tigerConfirmDialog onClose={() => {setOpenDeleteDialog(false);}} onSubmit={onAcceptDelete} open={openDelete}/>
                  }
                  
                </tigerStyledTableCell>
                <tigerStyledTableCell>
                    <IconButton size="small" onClick={()=>{setOpenUpdateDialog(true);}} color="info" variant="contained"><SettingsIcon color="warning" /></IconButton>
                  {openUpdate && 
                    <tigerDialog fullWidth={true} maxWidth="sm" open={openUpdate} 
                      onClose={()=>{setOpenUpdateDialog(false);}} 
                      innerDialog={<NewTagtypesDialog onClose={()=>{setOpenUpdateDialog(false);}} onSubmit={props.onUpdateTagtype} currentTag={props}/>}
                  />}
                </tigerStyledTableCell>
                <tigerStyledTableCell>
                  <Chip label={props.name} size="small" style={{backgroundColor:props.color}} />
                </tigerStyledTableCell>
                <tigerStyledTableCell>
                  {props.description}
                </tigerStyledTableCell>
                <tigerStyledTableCell>
                  {props.tags_aggregate.aggregate.count}
                </tigerStyledTableCell>
                <tigerStyledTableCell>
                  
                  <Box sx={{width: "100%", height: 25, backgroundColor: props.color}} >
                    <Typography style={{textAlign: "center", color: lightColor}} >
                      {"Sample Text Light Theme - "}{props.color}
                    </Typography>
                  </Box>
                  <Box sx={{width: "100%", height: 25, backgroundColor: props.color}} >
                    <Typography style={{textAlign: "center", color: darkColor}}>
                      {"Sample Text Dark Theme - "}{props.color}
                    </Typography>
                  </Box>
                </tigerStyledTableCell>
            </TableRow>
        </React.Fragment>
    )
}

