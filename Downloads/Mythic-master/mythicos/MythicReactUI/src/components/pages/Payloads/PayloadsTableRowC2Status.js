import React from 'react';
import Typography from '@mui/material/Typography';
import CancelIcon from '@mui/icons-material/Cancel';
import CheckCircleIcon from '@mui/icons-material/CheckCircle';
import {tigerStyledTooltip} from '../../tigerComponents/tigerStyledTooltip';
import PermScanWifiIcon from '@mui/icons-material/PermScanWifi';

export function PayloadsTableRowC2Status(props){
    return (
        <React.Fragment>
            {
                props.payloadc2profiles.map( (c2) => (
                    <Typography key={c2.c2profile.name + props.uuid} style={{display: "flex"}}> 
                        {c2.c2profile.is_p2p ?
                            ( c2.c2profile.container_running ? 
                                <tigerStyledTooltip title="C2 Container online">
                                    <CheckCircleIcon color="success"/>
                                </tigerStyledTooltip>: 
                                <tigerStyledTooltip title="C2 Container offline">
                                    <CancelIcon color="error"/>
                                </tigerStyledTooltip> )
                            :
                        ( c2.c2profile.running ? 
                            <tigerStyledTooltip title="C2 Internal Server Running">
                                <CheckCircleIcon color="success"/>
                            </tigerStyledTooltip> : 
                            (c2.c2profile.container_running ? (
                                <tigerStyledTooltip title="C2 Internal Server Not Running, but Container Online">
                                    <PermScanWifiIcon color="warning"/> 
                                </tigerStyledTooltip>
                            ) : (
                                <tigerStyledTooltip title="C2 Container offline">
                                    <CancelIcon color="error"/> 
                                </tigerStyledTooltip>
                            ))
                            )
                        } - {c2.c2profile.name}
                    </Typography>
                )) 
            }
                
        </React.Fragment>
        )
}
