import {tigerTabPanel, tigerTabLabel} from '../../tigerComponents/tigerTabPanel';
import React, {useEffect, useRef} from 'react';
import {gql, useQuery, useSubscription } from '@apollo/client';
import { tigerDialog } from '../../tigerComponents/tigerDialog';
import {useTheme} from '@mui/material/styles';
import Grid from '@mui/material/Grid';
import RefreshIcon from '@mui/icons-material/Refresh';
import IconButton from '@mui/material/IconButton';
import VisibilityOffIcon from '@mui/icons-material/VisibilityOff';
import VisibilityIcon from '@mui/icons-material/Visibility';
import {CallbacksTabsProcessBrowserTable} from './CallbacksTabsProcessBrowserTable';
import {tigerModifyStringDialog} from '../../tigerComponents/tigerDialog';
import {TaskFromUIButton} from './TaskFromUIButton';
import { tigerStyledTooltip } from '../../tigerComponents/tigerStyledTooltip';
import MenuItem from '@mui/material/MenuItem';
import FormControl from '@mui/material/FormControl';
import Select from '@mui/material/Select';
import InputLabel from '@mui/material/InputLabel';
import Input from '@mui/material/Input';
import {ViewCallbacktigerTreeGroupsDialog} from "./ViewCallbacktigerTreeGroupsDialog";
import WidgetsIcon from '@mui/icons-material/Widgets';

const treeFragment = gql`
fragment treeObjData on tigertree {
    comment
    deleted
    task_id
    filemeta {
        id
    }
    tags {
        tagtype {
            name
            color
            id
        }
        id
    }
    host
    id
    os
    can_have_children
    success
    full_path_text
    name_text
    timestamp
    parent_path_text
    tree_type
    metadata
    callback {
        id
        display_id
        tigertree_groups
    }
}
`;
const treeSubscription = gql`
    ${treeFragment}
    subscription liveData($now: timestamp!, $operation_id: Int!) {
        tigertree_stream(
            batch_size: 1000,
            cursor: {initial_value: {timestamp: $now}},
            where: { operation_id: { _eq: $operation_id }, tree_type: {_eq: "process"} }
        ) {
            ...treeObjData
        }
    }
`;
const rootQuery = gql`
    ${treeFragment}
    query myRootFolderQuery($operation_id: Int!) {
        tigertree(where: { operation_id: { _eq: $operation_id }, tree_type: {_eq: "process"} }, order_by: {id: asc}) {
            ...treeObjData
        }
    }
`;
export const uniqueSplitString = "$&%^";
export function CallbacksTabsProcessBrowserLabel(props){
    const [description, setDescription] = React.useState("Processes: " + props.tabInfo.displayID)
    const [openEditDescriptionDialog, setOpenEditDescriptionDialog] = React.useState(false);
    const contextMenuOptions = props.contextMenuOptions.concat([
        {
            name: 'Set Tab Description', 
            click: ({event}) => {
                setOpenEditDescriptionDialog(true);
            }
        },
    ]);
    useEffect( () => {
        if(props.tabInfo.customDescription !== "" && props.tabInfo.customDescription !== undefined){
            setDescription(props.tabInfo.customDescription);
        }else{
            setDescription("Processes: " + props.tabInfo.displayID);
        }
    }, [props.tabInfo.customDescription])
    useEffect( () => {
        let savedDescription = localStorage.getItem(`${props.me.user.id}-${props.tabInfo.operation_id}-${props.tabInfo.tabID}`);
        if(savedDescription && savedDescription !== ""){
            setDescription(savedDescription);
        }
    }, []);
    const editDescriptionSubmit = (description) => {
        props.onEditTabDescription(props.tabInfo, description);
        localStorage.setItem(`${props.me.user.id}-${props.tabInfo.operation_id}-${props.tabInfo.tabID}`, description);
    }
    return (
        <React.Fragment>
            <tigerTabLabel label={description} onDragTab={props.onDragTab}  {...props} contextMenuOptions={contextMenuOptions}/>
            {openEditDescriptionDialog &&
                <tigerDialog fullWidth={true} open={openEditDescriptionDialog}  onClose={() => {setOpenEditDescriptionDialog(false);}}
                    innerDialog={
                        <tigerModifyStringDialog title={"Edit Tab's Description"} onClose={() => {setOpenEditDescriptionDialog(false);}} value={description} onSubmit={editDescriptionSubmit} />
                    }
                />
            }
        </React.Fragment>  
    )
}
export const CallbacksTabsProcessBrowserPanel = ({index, value, tabInfo, me}) =>{
    const [fromNow, setFromNow] = React.useState((new Date()));
    const treeRootDataRef = React.useRef({}); // hold all the actual data
    const [treeAdjMtx, setTreeAdjMtx] = React.useState({}); // hold the simple adjacency matrix for parent/child relationships
    const [openTaskingButton, setOpenTaskingButton] = React.useState(false);
    const taskingData = React.useRef({"parameters": "", "ui_feature": "process_browser:list"});
    const mountedRef = React.useRef(true);
    const [showDeletedFiles, setShowDeletedFiles] = React.useState(false);
    const [selectedHost, setSelectedHost] = React.useState("");
    const [selectedGroup, setSelectedGroup] = React.useState("");
    useQuery(rootQuery, {
        variables: { operation_id: me?.user?.current_operation_id ||0},
        onCompleted: (data) => {
           // use an adjacency matrix but only for full_path_text -> children, not both directions
            for(let i = 0; i < data.tigertree.length; i++){
                let currentGroups = data.tigertree[i]?.["callback"]?.["tigertree_groups"] || ["Unknown Callbacks"];
                for(let j = 0; j < currentGroups.length; j++){
                    if(treeRootDataRef.current[currentGroups[j]] === undefined){
                        treeRootDataRef.current[currentGroups[j]] = {};
                    }
                    if( treeRootDataRef.current[currentGroups[j]][data.tigertree[i]["host"]] === undefined) {
                        // new host discovered
                        treeRootDataRef.current[currentGroups[j]][data.tigertree[i]["host"]] = {};
                    }
                    treeRootDataRef.current[currentGroups[j]][data.tigertree[i]["host"]][data.tigertree[i]["full_path_text"] /*+ uniqueSplitString + data.tigertree[i]["callback_id"]*/] = {...data.tigertree[i]}
                }
            }
            // create the top level data in the adjacency matrix
            const newMatrix = data.tigertree.reduce( (prev, cur) => {
                let currentGroups = cur?.["callback"]?.["tigertree_groups"] || ["Unknown Callbacks"];
                for(let j = 0; j < currentGroups.length; j++) {
                    if (prev[currentGroups[j]] === undefined) {
                        prev[currentGroups[j]] = {};
                    }
                    if (prev[currentGroups[j]][cur["host"]] === undefined) {
                        // the current host isn't tracked in the adjacency matrix, so add it
                        prev[currentGroups[j]][cur["host"]] = {}
                    }
                    if(cur["parent_path_text"] === ""){
                        if(prev[currentGroups[j]][cur['host']][""] === undefined){
                            prev[currentGroups[j]][cur['host']][""] = {}
                        }
                        prev[currentGroups[j]][cur['host']][""][cur["full_path_text"] /*+ uniqueSplitString + cur["callback_id"]*/] = 1;
                        continue
                    }
                    if (prev[currentGroups[j]][cur["host"]][cur["parent_path_text"] /*+ uniqueSplitString + cur["callback_id"]*/] === undefined) {
                        // the current parent's path isn't tracked, so add it and ourselves as children
                        prev[currentGroups[j]][cur["host"]][cur["parent_path_text"] /*+ uniqueSplitString + cur["callback_id"]*/] = {};
                    }
                    prev[currentGroups[j]][cur["host"]][cur["parent_path_text"] /*+ uniqueSplitString + cur["callback_id"]*/][cur["full_path_text"] /*+ uniqueSplitString + cur["callback_id"]*/] = 1;
                }
                return prev;
            }, {...treeAdjMtx});
           setTreeAdjMtx(newMatrix);
           // first see if we can find a group that matches our host, if not, then we can do first of each
            let groups = Object.keys(newMatrix).sort();
           if(groups.length > 0){
               for(let i = 0; i < groups.length; i++){
                   const hosts = Object.keys(newMatrix[groups[i]]).sort();
                   if(hosts.length > 0){
                       if(hosts.includes(tabInfo.host)){
                           setSelectedGroup(groups[i]);
                           setSelectedHost(tabInfo.host);
                           return;
                       }
                   }
               }
               setSelectedGroup(groups[0]);
               const hosts = Object.keys(newMatrix[groups[0]]).sort();
               if(hosts.length > 0){
                   setSelectedHost(hosts[0]);
               }
           }
        },
        fetchPolicy: 'no-cache',
    });
    useSubscription(treeSubscription, {
        variables: {now: fromNow, operation_id: me?.user?.current_operation_id ||0},
        fetchPolicy: "no-cache",
        onData: ({data}) => {
            for(let i = 0; i < data.data.tigertree_stream.length; i++){
                let currentGroups = data.data.tigertree_stream[i]?.["callback"]?.["tigertree_groups"] || ["Unknown Callbacks"];
                for(let j = 0; j < currentGroups.length; j++) {
                    if (treeRootDataRef.current[currentGroups[j]] === undefined) {
                        treeRootDataRef.current[currentGroups[j]] = {};
                    }
                    if (treeRootDataRef.current[currentGroups[j]][data.data.tigertree_stream[i]["host"]] === undefined) {
                        // new host discovered
                        treeRootDataRef.current[currentGroups[j]][data.data.tigertree_stream[i]["host"]] = {};
                    }
                    treeRootDataRef.current[currentGroups[j]][data.data.tigertree_stream[i]["host"]][data.data.tigertree_stream[i]["full_path_text"] /*+ uniqueSplitString + data.data.tigertree_stream[i]["callback_id"]*/] = {...data.data.tigertree_stream[i]};
                }
            }
            const newMatrix = data.data.tigertree_stream.reduce( (prev, cur) => {
                let currentGroups = cur?.["callback"]?.["tigertree_groups"] || ["Unknown Callbacks"];
                for(let j = 0; j < currentGroups.length; j++) {
                    if (prev[currentGroups[j]] === undefined) {
                        prev[currentGroups[j]] = {};
                    }
                    if (prev[currentGroups[j]][cur["host"]] === undefined) {
                        // the current host isn't tracked in the adjacency matrix, so add it
                        prev[currentGroups[j]][cur["host"]] = {}
                    }
                    if(cur["parent_path_text"] === ""){
                        if(prev[currentGroups[j]][cur['host']][""] === undefined){
                            prev[currentGroups[j]][cur['host']][""] = {}
                        }
                        prev[currentGroups[j]][cur['host']][""][cur["full_path_text"] /*+ uniqueSplitString + cur["callback_id"]*/] = 1;
                        continue
                    }
                    if (prev[currentGroups[j]][cur["host"]][cur["parent_path_text"] /*+ uniqueSplitString + cur["callback_id"]*/] === undefined) {
                        // the current parent's path isn't tracked, so add it and ourselves as children
                        prev[currentGroups[j]][cur["host"]][cur["parent_path_text"] /*+ uniqueSplitString + cur["callback_id"]*/] = {};
                    }
                    prev[currentGroups[j]][cur["host"]][cur["parent_path_text"] /*+ uniqueSplitString + cur["callback_id"]*/][cur["full_path_text"] /*+ uniqueSplitString + cur["callback_id"]*/] = 1;
                }
                return prev;
            }, {...treeAdjMtx});
            setTreeAdjMtx(newMatrix);
        }
    })
    const onListFilesButton = () => {
        taskingData.current = ({parameters: "",
            ui_feature: "process_browser:list",
            callback_id: tabInfo["callbackID"],
            callback_display_id: tabInfo["displayID"]});
        setOpenTaskingButton(true);
    }
    const onTaskRowAction = ({process_id, architecture, uifeature, openDialog, getConfirmation, callback_id, callback_display_id}) => {
        taskingData.current = {"parameters": {host: selectedHost, process_id, architecture},
            "ui_feature": uifeature, openDialog, getConfirmation, callback_id, callback_display_id};
        setOpenTaskingButton(true);
    }
    const toggleShowDeletedFiles = (showStatus) => {
        setShowDeletedFiles(showStatus);
    };
    const updateSelectedHost = (host) => {
        setSelectedHost(host);
    }
    const updateSelectedGroup = (group) => {
        setSelectedGroup(group);
        const hosts = Object.keys(treeAdjMtx[group]);
        if(hosts.length > 0){
            setSelectedHost(hosts[0]);
        } else {
            setSelectedHost("");
        }
    }
    React.useEffect( () => {
        return() => {
            mountedRef.current = false;
        }
         // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [])
    return (
        <tigerTabPanel index={index} value={value} >
            <div style={{display: "flex", flexGrow: 1, overflowY: "auto"}}>
                <div style={{width: "100%", display: "flex", flexDirection: "column", flexGrow: 1}}>
                    <ProcessBrowserTableTop 
                        onListFilesButton={onListFilesButton}
                        tabInfo={tabInfo}
                        host={selectedHost}
                        group={selectedGroup}
                        toggleShowDeletedFiles={toggleShowDeletedFiles}
                        updateSelectedHost={updateSelectedHost}
                        updateSelectedGroup={updateSelectedGroup}
                        groupOptions={treeRootDataRef.current}
                        hostOptions={treeRootDataRef.current[selectedGroup] || {}}
                    />
                    <CallbacksTabsProcessBrowserTable 
                        showDeletedFiles={showDeletedFiles}
                        tabInfo={tabInfo}
                        onRowDoubleClick={() => {}}
                        treeRootData={treeRootDataRef.current}
                        treeAdjMatrix={treeAdjMtx}
                        host={selectedHost}
                        group={selectedGroup}
                        onTaskRowAction={onTaskRowAction}
                        me={me}/>
                  
                </div>
                {openTaskingButton && 
                    <TaskFromUIButton ui_feature={taskingData.current?.ui_feature || " "} 
                        callback_id={taskingData.current?.callback_id || tabInfo.callbackID}
                        parameters={taskingData.current?.parameters || ""}
                        openDialog={taskingData.current?.openDialog || false}
                        getConfirmation={taskingData.current?.getConfirmation || false}
                        selectCallback={taskingData.current?.selectCallback || false}
                        onTasked={() => setOpenTaskingButton(false)}/>
                    }
            </div>            
        </tigerTabPanel>
    )
}
const ProcessBrowserTableTop = ({
    onListFilesButton,
    updateSelectedHost,
    updateSelectedGroup,
    toggleShowDeletedFiles,
    host,
    group,
    hostOptions,
    groupOptions
}) => {
    const theme = useTheme();
    const [showDeletedFiles, setLocalShowDeletedFiles] = React.useState(false);
    const [openViewGroupsDialog, setOpenViewGroupDialog] = React.useState(false);
    const inputRef = useRef(null);
    const inputGroupRef = useRef(null);
    const onLocalListFilesButton = () => {
        onListFilesButton()
    }
    const onLocalToggleShowDeletedFiles = () => {
        setLocalShowDeletedFiles(!showDeletedFiles);
        toggleShowDeletedFiles(!showDeletedFiles);
    };
    const handleChange = (event) => {
        updateSelectedHost(event.target.value);
    }
    const handleGroupChange = (event) => {
        updateSelectedGroup(event.target.value);
    }
    return (
        <Grid container spacing={0} style={{paddingTop: "10px"}}>
            <Grid item xs={12}>
                <FormControl style={{width: "30%"}}>
                    <InputLabel ref={inputGroupRef}>Available Groups</InputLabel>
                    <Select
                        labelId="demo-dialog-select-label"
                        id="demo-dialog-select"
                        value={group}
                        onChange={handleGroupChange}
                        input={<Input style={{width: "100%"}}/>}
                        endAdornment={
                            <React.Fragment>
                                <tigerStyledTooltip title="View Callbacks associated with this group">
                                    <IconButton style={{padding: "3px"}} onClick={() => {setOpenViewGroupDialog(true);}} size="large"><WidgetsIcon style={{color: theme.palette.info.main}}/></IconButton>
                                </tigerStyledTooltip>
                            </React.Fragment>
                        }
                    >
                        {Object.keys(groupOptions).sort().map( (opt) => (
                            <MenuItem value={opt} key={opt}>{opt}</MenuItem>
                        ) )}
                    </Select>
                </FormControl>
                <FormControl style={{width: "70%"}}>
                  <InputLabel ref={inputRef}>Available Hosts</InputLabel>
                  <Select
                    labelId="demo-dialog-select-label"
                    id="demo-dialog-select"
                    value={host}
                    onChange={handleChange}
                    input={<Input style={{width: "100%"}}/>}
                    endAdornment={
                <React.Fragment>
                    <tigerStyledTooltip title="Task Callback to List Processes">
                        <IconButton style={{padding: "3px"}} onClick={onLocalListFilesButton} size="large"><RefreshIcon style={{color: theme.palette.info.main}}/></IconButton>
                    </tigerStyledTooltip>
                    <tigerStyledTooltip title={showDeletedFiles ? 'Hide Deleted Processes' : 'Show Deleted Processes'}>
                            <IconButton
                                style={{ padding: '3px' }}
                                onClick={onLocalToggleShowDeletedFiles}
                                size="large">
                                {showDeletedFiles ? (
                                    <VisibilityIcon color="success" />
                                ) : (
                                    <VisibilityOffIcon color="error"  />
                                )}
                            </IconButton>
                        </tigerStyledTooltip>
                </React.Fragment>
                    }
                  >
                    {Object.keys(hostOptions).sort().map( (opt) => (
                        <MenuItem value={opt} key={opt}>{opt}</MenuItem>
                    ) )}
                  </Select>
                </FormControl>
                {openViewGroupsDialog &&
                    <tigerDialog
                        fullWidth={true}
                        maxWidth={"xl"}
                        open={openViewGroupsDialog}
                        onClose={() => {setOpenViewGroupDialog(false);}}
                        innerDialog={
                            <ViewCallbacktigerTreeGroupsDialog group_name={group}
                                                                onClose={() => {setOpenViewGroupDialog(false);}} />
                        }
                    />
                }
            </Grid>
        </Grid>
    );
}
