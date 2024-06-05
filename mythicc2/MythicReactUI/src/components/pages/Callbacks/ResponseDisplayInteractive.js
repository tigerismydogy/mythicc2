import React, { useEffect } from 'react';
import MythicTextField from '../../MythicComponents/MythicTextField';
import {createTaskingMutation, taskingDataFragment} from './CallbackMutations';
import {snackActions} from "../../utilities/Snackbar";
import { gql, useMutation, useSubscription } from '@apollo/client';
import {b64DecodeUnicode} from "./ResponseDisplay";
import {SearchBar} from './ResponseDisplay';
import Pagination from '@mui/material/Pagination';
import {Typography, CircularProgress, Select, IconButton} from '@mui/material';
import CheckCircleOutlineIcon from '@mui/icons-material/CheckCircleOutline';
import Input from '@mui/material/Input';
import MenuItem from '@mui/material/MenuItem';
import Anser from "anser";
import PaletteIcon from '@mui/icons-material/Palette';
import {MythicStyledTooltip} from "../../MythicComponents/MythicStyledTooltip";
import './ResponseDisplayInteractiveANSI.css';
import KeyboardReturnIcon from '@mui/icons-material/KeyboardReturn';
import InputLabel from "@mui/material/InputLabel";
import FormControl from "@mui/material/FormControl";
import WrapTextIcon from '@mui/icons-material/WrapText';

const getInteractiveTaskingQuery = gql`
${taskingDataFragment}
subscription getTasking($parent_task_id: Int!){
    task_stream(where: {parent_task_id: {_eq: $parent_task_id}, is_interactive_task: {_eq: true}}, batch_size: 50, cursor: {initial_value: {timestamp: "1970-01-01"}}) {
        ...taskData
    }
}
 `;
const subResponsesQuery = gql`
subscription subResponsesQuery($task_id: Int!) {
  response_stream(where: {task_id: {_eq: $task_id}}, batch_size: 50, cursor: {initial_value: {id: 0}}) {
    id
    response: response_text
    timestamp
    is_error
  }
}`;
const getTaskingStatus = (task) => {
    if(task.status === "completed"){
        return <CheckCircleOutlineIcon color={"success"} fontSize={"1rem"} style={{marginRight: "2px"}} />
    } else if(task.status === "submitted"){
        return  <CircularProgress size={"1rem"} />
    }
}
const getClassnames = (entry) => {
    let classnames = [];
    if(entry.decoration && entry.decoration === "reverse"){
        if(entry.fg){
            classnames.push(entry.fg + "-background");
        }
        if(entry.bg){
            classnames.push(entry.bg);
        }
    } else {
        if(entry.fg){
            classnames.push(entry.fg);
        }
        if(entry.bg){
            classnames.push(entry.bg + "-background");
        }
    }

    if(entry.decoration){
        classnames.push("ansi-" + entry.decoration);
    }
    for(let i = 0; i < entry.decorations.length; i++){
        classnames.push("ansi-" + entry.decorations[i])
    }
    if(entry.bg_truecolor || entry.fg_truecolor){
        console.log(entry);
    }
    //console.log(entry);
    return classnames.join(" ");
}
const GetOutputFormat = ({data, useASNIColor, messagesEndRef, showTaskStatus, wrapText}) => {
    const [dataElement, setDataElement] = React.useState(null);
    React.useEffect( () => {

        if(data.response) {
            // we're looking at response output
            if(data.is_error){
                setDataElement(<pre style={{display: "inline",backgroundColor: "#311717", margin: "0 0 0 0",
                    wordBreak: wrapText ? "break-all" : "",
                    whiteSpace: wrapText ? "pre-wrap" : ""}} key={data.timestamp + data.id}>
                    {data.response}
                </pre>)
            } else {
                if(useASNIColor){
                    let ansiJSON = Anser.ansiToJson(data.response, { use_classes: true });
                    //console.log(ansiJSON)
                    setDataElement(
                        ansiJSON.map( (a, i) => (
                            <pre style={{display: "inline", margin: "0 0 0 0",
                                wordBreak: wrapText ? "break-all" : "",
                                whiteSpace: wrapText ? "pre-wrap" : "",
                            }} className={getClassnames(a)} key={data.id + data.timestamp + i}>{a.content}</pre>
                        ))
                    )
                } else {
                    setDataElement(<pre style={{display: "inline", margin: "0 0 0 0",
                        wordBreak: wrapText ? "break-all" : "",
                        whiteSpace: wrapText ? "pre-wrap" : "",
                    }} key={data.timestamp + data.id}>{data.response}</pre>)
                }

            }
        } else {
            // we're looking at tasking
            setDataElement(
                <pre key={data.timestamp + data.id} style={{display: "inline",margin: "0 0 0 0",
                    wordBreak: wrapText ? "break-all" : "", whiteSpace: "pre-wrap"}}>
                    {showTaskStatus && getTaskingStatus(data)}
                    {data.original_params}
                </pre>
            )
        }
    }, [data.timestamp, useASNIColor, showTaskStatus, wrapText]);
    /*
    React.useEffect( () => {
        if(messagesEndRef.current){
            messagesEndRef.current.scrollIntoView({ behavior: "smooth" });
        }
    }, [dataElement]);

     */
    return (
        dataElement
    )

}

const InteractiveMessageTypes = [
    {"name": "None", "value": -1, "text": "None"},
    {"name": "Tab", "value": 13, "text": "^I"},
    {"name": "Backspace", "value": 12, "text": "^H"},
    {"name": "Escape", "value": 4, "text": "^["},
    {"name": "Ctrl+A", "value": 5, "text": "^A"},
    {"name": "Ctrl+B", "value": 6, "text": "^B"},
    {"name": "Ctrl+C", "value": 7, "text": "^C"},
    {"name": "Ctrl+D", "value": 8, "text": "^D"},
    {"name": "Ctrl+E", "value": 9, "text": "^E"},
    {"name": "Ctrl+F", "value": 10, "text": "^F"},
    {"name": "Ctrl+G", "value": 11, "text": "^G"},
    {"name": "Ctrl+K", "value": 14, "text": "^K"},
    {"name": "Ctrl+L", "value": 15, "text": "^L"},
    {"name": "Ctrl+N", "value": 16, "text": "^N"},
    {"name": "Ctrl+P", "value": 17, "text": "^P"},
    {"name": "Ctrl+Q", "value": 18, "text": "^Q"},
    {"name": "Ctrl+R", "value": 19, "text": "^R"},
    {"name": "Ctrl+S", "value": 20, "text": "^S"},
    {"name": "Ctrl+U", "value": 21, "text": "^U"},
    {"name": "Ctrl+W", "value": 22, "text": "^W"},
    {"name": "Ctrl+Y", "value": 23, "text": "^Y"},
    {"name": "Ctrl+Z", "value": 24, "text": "^Z"},
]
const EnterOptions = [
    {"value": "", "name": "None"},
    {"value": "\n", "name": "LF"},
    {"value": "\r", "name": "CR"},
    {"value": "\r\n", "name": "CRLF"}
];
export const ResponseDisplayInteractive = (props) =>{
    const pageSize = React.useRef(100);
    const highestFetched = React.useRef(0);
    const [taskData, setTaskData] = React.useState([]);
    const [alloutput, setAllOutput] = React.useState([]);
    const [rawResponses, setRawResponses] = React.useState([]);
    const [search, setSearch] = React.useState("");
    const [totalCount, setTotalCount] = React.useState(0);
    const messagesEndRef = React.useRef();
    const page = React.useRef(1);
    const taskIDRef = React.useRef(props.task.id);
    const taskIDTaskingRef = React.useRef(props.task.id);
    const [useASNIColor, setUseANSIColor] = React.useState(true);
    const [showTaskStatus, setShowTaskStatus] = React.useState(true);
    const [wrapText, setWrapText] = React.useState(true);
    useSubscription(getInteractiveTaskingQuery, {
      variables: {parent_task_id: props.task.id},
      shouldResubscribe: true,
      onError: data => {
          console.error(data)
      },
      fetchPolicy: "no-cache",
      onData: ({data}) => {
          if(props.task.id !== taskIDTaskingRef.current){
              const newTaskData = data.data.task_stream;
              setTaskData(newTaskData);
              taskIDTaskingRef.current = props.task.id;
          } else {
              const newTaskData = data.data.task_stream.reduce( (prev, cur) => {
                  let taskIndex = prev.findIndex(t => t.id === cur.id);
                  if(taskIndex >= 0){
                      prev[taskIndex] = {...cur}
                      return prev
                  }
                  return [...prev, cur]
              }, [...taskData]);
              setTaskData(newTaskData);
          }
      }
  })
    const subscriptionDataCallback = React.useCallback( ({data}) => {
        // we still have some room to view more, but only room for fetchLimit - totalFetched.current
        if(props.task.id !== taskIDRef.current){
            const newResponses = data.data.response_stream;
            const newerResponses = newResponses.map( (r) => { return {...r, response: b64DecodeUnicode(r.response)}});
            newerResponses.sort( (a,b) => a.id > b.id ? 1 : -1);
            let rawResponseArray = [...rawResponses];
            let highestFetchedId = 0;
            for(let i = 0; i < newerResponses.length; i++){
                rawResponseArray.push(newerResponses[i]);
                highestFetchedId = newerResponses[i]["id"];
            }
            setRawResponses(rawResponseArray);
            highestFetched.current = highestFetchedId;
            taskIDRef.current = props.task.id;
        } else {
            const newResponses = data.data.response_stream.filter( r => r.id > highestFetched.current);
            const newerResponses = newResponses.map( (r) => { return {...r, response: b64DecodeUnicode(r.response)}});
            newerResponses.sort( (a,b) => a.id > b.id ? 1 : -1);
            let rawResponseArray = [...rawResponses];
            let highestFetchedId = highestFetched.current;
            for(let i = 0; i < newerResponses.length; i++){
                rawResponseArray.push(newerResponses[i]);
                highestFetchedId = newerResponses[i]["id"];
            }
            setRawResponses(rawResponseArray);
            highestFetched.current = highestFetchedId;
        }

    }, [highestFetched.current, rawResponses, props.task.id]);
    useSubscription(subResponsesQuery, {
        variables: {task_id: props.task.id},
        fetchPolicy: "network_only",
        onData: subscriptionDataCallback
    });
    const onSubmitPageChange = (currentPage) => {

        if(search === undefined || search === ""){
            let allData = [...rawResponses, ...taskData];
            allData.sort( (a,b) => {
                let aDate = new Date(a.timestamp);
                let bDate = new Date(b.timestamp);
                return aDate < bDate ? -1 : bDate < aDate ? 1 : 0;
            })
            if(page.current === currentPage){
                const pageCount = Math.max(1, Math.ceil(allData.length / pageSize.current));
                if(page.current === pageCount -1){
                    console.log("updating pageSize");
                    // we just streamed more data and we're on the latest page, increase pageSize
                    //pageSize.current += Math.abs(allData.length - totalCount);
                    currentPage += 1;
                }
            }
            setAllOutput(allData.slice((currentPage-1)*pageSize.current, currentPage*pageSize.current));
            setTotalCount(allData.length);
        }else{
            let allData = [...rawResponses, ...taskData];
            allData = allData.filter( (r) => {
                if(r.response){
                    return r.response.includes(search);
                } else {
                    return r.display_params.includes(search);
                }
            })
            allData.sort( (a,b) => {
                let aDate = new Date(a.timestamp);
                let bDate = new Date(b.timestamp);
                return aDate < bDate ? -1 : bDate < aDate ? 1 : 0;
            });

            if(page.current === currentPage){
                const pageCount = Math.max(1, Math.ceil(allData.length / pageSize.current));
                if(page.current === pageCount -1){
                    // we just streamed more data and we're on the latest page, increase pageSize
                    //pageSize.current += Math.abs(allData.length - totalCount);
                    currentPage += 1;
                }
            }
            setTotalCount(allData.length);
            setAllOutput(allData.slice((currentPage-1)*pageSize.current, currentPage*pageSize.current));
        }
        page.current = currentPage;
    }
    React.useEffect( () => {
        onSubmitPageChange(1);
    }, [search]);
    const onSubmitSearch = React.useCallback( (newSearch) => {
        setSearch(newSearch);
    }, []);
    useEffect( () =>{
        onSubmitPageChange(page.current);
        /*
        let allData = [...rawResponses, ...taskData];
        allData.sort( (a,b) => {
            let aDate = new Date(a.timestamp);
            let bDate = new Date(b.timestamp);
            return aDate < bDate ? -1 : bDate < aDate ? 1 : 0;
        });
        setAllOutput(allData);
        setTotalCount(allData.length);

         */
    }, [rawResponses, taskData]);
    const toggleANSIColor = () => {
        setUseANSIColor(!useASNIColor);
    }
    const toggleShowTaskStatus = () => {
        setShowTaskStatus(!showTaskStatus);
    }
    const toggleWrapText = () => {
        setWrapText(!wrapText);
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

    <div style={{ display: "flex", overflowY: "auto",
        position: "relative", height: props.expand ? "100%": undefined, maxHeight: props.expand ? "100%": "500px",
        flexDirection: "column"}}>
        {props.searchOutput &&
            <SearchBar onSubmitSearch={onSubmitSearch} />
        }
        <div style={{overflowY: "auto", width: "100%", flexGrow: 1}} ref={props.responseRef}>
            {alloutput.map( (e, index) => (
                <GetOutputFormat key={"getoutput" + index} data={e} useASNIColor={useASNIColor} messagesEndRef={messagesEndRef}
                                 showTaskStatus={showTaskStatus} wrapText={wrapText} />
            ))}
            <div ref={messagesEndRef}/>
        </div>

        <div style={{width: "100%", display: "inline-flex", alignItems: "flex-end"}} >
            <InteractiveTaskingBar task={props.task} taskData={taskData}
                                   useASNIColor={useASNIColor} toggleANSIColor={toggleANSIColor}
                                   showTaskStatus={showTaskStatus} toggleShowTaskStatus={toggleShowTaskStatus}
                                   wrapText={wrapText} toggleWrapText={toggleWrapText}
            />
        </div>
            <InteractivePaginationBar totalCount={totalCount} currentPage={page.current}
                                      onSubmitPageChange={onSubmitPageChange} expand={props.expand}
                                      pageSize={pageSize.current} />
    </div>
  )
      
}
const InteractiveTaskingBar = ({task, taskData, useASNIColor, toggleANSIColor,
                                   showTaskStatus, toggleShowTaskStatus, wrapText, toggleWrapText}) => {
    const [inputText, setInputText] = React.useState("");
    const [createTask] = useMutation(createTaskingMutation, {
        update: (cache, {data}) => {
            if(data.createTask.status === "error"){
                snackActions.error(data.createTask.error);
            }else{
            }
        },
        onError: data => {
            console.error(data);
        }
    });
    const [taskOptionsIndex, setTaskOptionsIndex] = React.useState(-1);
    const [taskOptions, setTaskOptions] = React.useState([]);
    const [selectedEnterOption, setSelectedEnterOption] = React.useState(EnterOptions[1]);
    const [selectedControlOption, setSelectedControlOption] = React.useState(InteractiveMessageTypes[0]);
    React.useEffect( () => {
        const newTaskOptions = taskData.filter( t => t.display_params.length > 1 && (t.interactive_task_type === 0 || t.interactive_task_type === 8));
        newTaskOptions.sort( (a,b) => a.id > b.id ? -1 : 1);
        setTaskOptions(newTaskOptions);
    }, [taskData]);
    const onInputChange = (name, value, error, event) => {
        if(event.key === "ArrowUp"){
            event.preventDefault();
            event.stopPropagation();
            if(taskOptions.length === 0){
                snackActions.warning("No previous tasks")
                return;
            }
            let newIndex = (taskOptionsIndex + 1);
            if(newIndex > taskOptions.length -1){
                newIndex = taskOptions.length -1;
            }
            setTaskOptionsIndex(newIndex);
            setInputText(taskOptions[newIndex].display_params.trim());
        }else if(event.key === "ArrowDown"){
            event.preventDefault();
            event.stopPropagation();
            if(taskData.length === 0){
                snackActions.warning("No next tasks")
                return;
            }
            let newIndex = (taskOptionsIndex -1);
            if(newIndex < 0){
                setTaskOptionsIndex(-1);
                setInputText("");
            } else {
                setTaskOptionsIndex(newIndex);
                setInputText(taskOptions[newIndex].display_params.trim());
            }

        }else{
            setInputText(value);
        }

    }
    const submitTask = (event) => {
        event.stopPropagation();
        event.preventDefault();
        if(event.shiftKey){
            setInputText(inputText + selectedEnterOption.value);
            return;
        }
        if(selectedControlOption.value > 0){
            let ctrlSequence = selectedControlOption.text;
            let enterOption = selectedEnterOption.value;
            let originalParams = inputText + ctrlSequence + enterOption;
            // if we're looking at a tab, never send enter along with it
            if (selectedControlOption.value === 8){
                originalParams = inputText + ctrlSequence;
                enterOption = "";
            // if we're looking at escape, never send enter along with it
            } else if(selectedControlOption.value === 4){
                originalParams = ctrlSequence + inputText;
                enterOption = "";
            }
            createTask({variables: {
                    callback_id: task.callback.display_id,
                    command: task.command.cmd,
                    params: inputText + enterOption,
                    tasking_location: "command_line",
                    original_params: originalParams,
                    parameter_group_name: "default",
                    parent_task_id: task.id,
                    is_interactive_task: true,
                    interactive_task_type: selectedControlOption.value,
                }})
        }else {
            // no control option selected, just send data along as input
            createTask({variables: {
                    callback_id: task.callback.display_id,
                    command: task.command.cmd,
                    params: inputText + selectedEnterOption.value,
                    tasking_location: "command_line",
                    original_params: inputText + selectedEnterOption.value,
                    parameter_group_name: "default",
                    parent_task_id: task.id,
                    is_interactive_task: true,
                    interactive_task_type: 0,
                }})
        }
        setInputText("");
        setSelectedControlOption(InteractiveMessageTypes[0]);
    }
    const onChangeSelect = (event) => {
        event.stopPropagation();
        event.preventDefault();
        setSelectedControlOption(event.target.value);
    }
    const onChangeSelectEnterOption = (event) => {
        event.stopPropagation();
        event.preventDefault();
        setSelectedEnterOption(event.target.value);
    }
    return (
        <>
            <FormControl style={{width: "10rem", marginTop: "2px"}} >
                <InputLabel id="control-label" style={{}}>Controls</InputLabel>
                <Select
                    labelId="control-label"
                    id="control-select"
                    value={selectedControlOption}
                    onChange={onChangeSelect}
                    input={<Input style={{margin: 0}} />}
                >
                    {InteractiveMessageTypes.map( (opt) => (
                        <MenuItem value={opt} key={opt.name}>{opt.name}</MenuItem>
                    ) )}
                </Select>
            </FormControl>

            <MythicTextField autoFocus={true} maxRows={5} multiline={true} onChange={onInputChange} onEnter={submitTask}
                             value={inputText} variant={"standard"} placeholder={">_ type here..."} inline={true}
                             marginBottom={"0px"} InputProps={{style: { width: "100%"}}}/>
            <FormControl style={{width: "6rem"}} >
                <Select
                    labelId="linefeed-label"
                    id="linefeed-control"
                    value={selectedEnterOption}
                    autoWidth
                    onChange={onChangeSelectEnterOption}
                    input={<Input />}
                    IconComponent={KeyboardReturnIcon}
                >
                    {EnterOptions.map( (opt) => (
                        <MenuItem value={opt} key={opt.name}>{opt.name}</MenuItem>
                    ) )}
                </Select>
            </FormControl>
            <MythicStyledTooltip title={useASNIColor ?  "Disable ANSI Color" : "Enable ANSI Color"} >
                <IconButton onClick={toggleANSIColor} style={{paddingLeft: 0, paddingRight: 0}}>
                    <PaletteIcon color={useASNIColor ? "success" : "secondary"}
                                 style={{cursor: "pointer"}}
                    />
                </IconButton>

            </MythicStyledTooltip>
            <MythicStyledTooltip title={showTaskStatus ?  "Hide Task Status" : "Show Task Status"} >
                <IconButton onClick={toggleShowTaskStatus} style={{paddingLeft: 0, paddingRight: 0}}>
                    <CheckCircleOutlineIcon color={showTaskStatus ? "success" : "secondary"}
                                            style={{cursor: "pointer",}}
                    />
                </IconButton>

            </MythicStyledTooltip>
            <MythicStyledTooltip title={wrapText ?  "Unwrap Text" : "Wrap Text"} >
                <IconButton onClick={toggleWrapText} style={{paddingLeft: 0, paddingRight: 0}}>
                    <WrapTextIcon color={wrapText ? "success" : "secondary"}
                                  style={{cursor: "pointer"}}
                    />
                </IconButton>

            </MythicStyledTooltip>
        </>
    )
}
const InteractivePaginationBar = ({totalCount, currentPage, onSubmitPageChange, pageSize, expand}) => {
    const onChangePage =  (event, value) => {
        onSubmitPageChange(value);
    };
    const pageCount = Math.max(1, Math.ceil(totalCount / pageSize));
    if(pageCount < 2){
        if(expand){
            return (
                <div style={{height: "50px"}}>

                </div>
            )
        } else {
            return null
        }

    }
    return (
        <div style={{background: "transparent", display: "flex", justifyContent: "center", alignItems: "center", paddingBottom: "10px",}} >
            <Pagination count={pageCount} page={currentPage} variant="contained" color="primary"
                        boundaryCount={2} onChange={onChangePage} style={{margin: "10px"}}
                        showFirstButton showLastButton siblingCount={2}
            />
            <Typography style={{paddingLeft: "10px"}}>Total: {totalCount}</Typography>
        </div>
    )
}