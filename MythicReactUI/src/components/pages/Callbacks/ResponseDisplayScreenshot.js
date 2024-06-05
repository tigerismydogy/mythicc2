import React from 'react';
import {Button} from '@mui/material';
import {ResponseDisplayScreenshotModal} from './ResponseDisplayScreenshotModal';
import { tigerDialog } from '../../tigerComponents/tigerDialog';
import {tigerStyledTooltip} from "../../tigerComponents/tigerStyledTooltip";


export const ResponseDisplayScreenshot = (props) =>{
  const [openScreenshot, setOpenScreenshot] = React.useState(false);

  const now = (new Date()).toUTCString();
  const clickOpenScreenshot = () => {
    setOpenScreenshot(true);
  }
    const scrollContent = (node, isAppearing) => {
        // only auto-scroll if you issued the task
        document.getElementById(`scrolltotaskbottom${props.task.id}`)?.scrollIntoView({
            //behavior: "smooth",
            block: "end",
            inline: "nearest"
        })
    }
    React.useLayoutEffect( () => {
        scrollContent()
    }, []);
  return (
    <>
      {openScreenshot &&
      <tigerDialog fullWidth={true} maxWidth="xl" open={openScreenshot} 
          onClose={()=>{setOpenScreenshot(false);}} 
          innerDialog={<ResponseDisplayScreenshotModal images={props.agent_file_id} onClose={()=>{setOpenScreenshot(false);}} />}
      />
      }
      <pre style={{display: "inline-block"}}>
        {props?.plaintext || ""}
      </pre>
      <tigerStyledTooltip title={props?.hoverText || "View Screenshot"}  >
        <Button color="primary" variant={props.variant ? props.variant : "contained"} onClick={clickOpenScreenshot} style={{marginBottom: "10px"}}>{props.name}</Button>
      </tigerStyledTooltip>
    </>
  );   
}