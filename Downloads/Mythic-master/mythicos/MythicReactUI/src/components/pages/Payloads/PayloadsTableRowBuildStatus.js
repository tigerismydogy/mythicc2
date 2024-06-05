import IconButton from '@mui/material/IconButton';
import CircularProgress from '@mui/material/CircularProgress';
import React from 'react';
import GetAppIcon from '@mui/icons-material/GetApp';
import { snackActions } from '../../utilities/Snackbar';
import ReportProblemIcon from '@mui/icons-material/ReportProblem';
import {tigerStyledTooltip} from '../../tigerComponents/tigerStyledTooltip';
import { tigerDialog } from '../../tigerComponents/tigerDialog';
import {PayloadBuildMessageDialog} from './PayloadBuildMessageDialog';
import { Link } from '@mui/material';

export function PayloadsTableRowBuildStatus(props){
    const [openBuildMessage, setOpenBuildMessageDialog] = React.useState(false);
    const onErrorClick = () => {
        snackActions.warning("Payload failed to build, cannot download");
        setOpenBuildMessageDialog(true);
    }
    return (
        <React.Fragment>
            {props.build_phase === "success" ?
                ( <tigerStyledTooltip title="Download payload">
                    <a href={"/direct/download/" + props.filemetum.agent_file_id} >
                        <GetAppIcon color="success" style={{marginLeft: "12px"}} />
                    </a>
                  </tigerStyledTooltip>
                    
                )
                : 
                (props.build_phase === "building" ? 
                (<tigerStyledTooltip title="Payload still building">
                    <IconButton variant="contained" size="large"><CircularProgress size={20} thickness={4} color="info"/></IconButton>
                </tigerStyledTooltip>) : 
                (<>
                    <IconButton
                        variant="contained"
                        onClick={onErrorClick}
                        disableFocusRipple={true}
                        disableRipple={true}
                        size="large">
                        <ReportProblemIcon color="error" />
                    </IconButton>
                    {openBuildMessage ? (
                    <tigerDialog fullWidth={true} maxWidth="lg" open={openBuildMessage} 
                        onClose={()=>{setOpenBuildMessageDialog(false);}} 
                        innerDialog={<PayloadBuildMessageDialog payload_id={props.id} viewError={true} onClose={()=>{setOpenBuildMessageDialog(false);}} />}
                    />
                ): null }
                </>
                ) 
                )
            }
        </React.Fragment>
    );
}

