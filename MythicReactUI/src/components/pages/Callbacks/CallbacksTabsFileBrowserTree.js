import React from 'react';
import FileBrowserVirtualTree from '../../tigerComponents/tigerFileBrowserVirtualTree';
import {tigerDialog} from "../../tigerComponents/tigerDialog";
import {ViewCallbacktigerTreeGroupsDialog} from "./ViewCallbacktigerTreeGroupsDialog";


export const CallbacksTabsFileBrowserTree = ({ treeRootData, treeAdjMatrix, fetchFolderData, setTableData, taskListing, tableOpenedPathId, showDeletedFiles, tabInfo}) => {
    const [openNodes, setOpenNodes] = React.useState({});
    const groupName = React.useRef("");
    const [openViewGroupsDialog, setOpenViewGroupDialog] = React.useState(false);
    const toggleNodeExpanded = (nodeId, nodeData) => {
        //console.log("toggleNodeExpanded", nodeId, nodeData);
        setTableData(nodeData);
        fetchFolderData(nodeData);
        setOpenNodes({
          ...openNodes,
          [nodeId]: true
        });
      };
    const toggleNodeCollapsed = (nodeId, nodeData) => {
        setOpenNodes({
          ...openNodes,
          [nodeId]: false
        });
      };
    const onSelectNode = (nodeId, nodeData) => {
        if(nodeData.root){

        }else if(nodeData.is_group){
            groupName.current = nodeData.group;
            setOpenViewGroupDialog(true);
        }else {
            //console.log(nodeData);
            setTableData(nodeData);
        }

    };
    React.useEffect( () => {
      setOpenNodes({
        ...openNodes,
        [tableOpenedPathId]: true
      });
    }, [tableOpenedPathId]);
    const contextMenuOptions = [
      {
          name: 'Task Listing', 
          click: ({event, node, callback_id, callback_display_id}) => {
              taskListing(node, callback_id, callback_display_id);
          }
      },
  ];
  return(
      <>
          <FileBrowserVirtualTree
              showDeletedFiles={showDeletedFiles}
              treeRootData={treeRootData}
              treeAdjMatrix={treeAdjMatrix}
              openNodes={openNodes}
              onSelectNode={onSelectNode}
              onExpandNode={toggleNodeExpanded}
              onCollapseNode={toggleNodeCollapsed}
              contextMenuOptions={contextMenuOptions}
              tabInfo={tabInfo}
          />
          {openViewGroupsDialog &&
              <tigerDialog
                  fullWidth={true}
                  maxWidth={"xl"}
                  open={openViewGroupsDialog}
                  onClose={() => {setOpenViewGroupDialog(false);}}
                  innerDialog={
                      <ViewCallbacktigerTreeGroupsDialog group_name={groupName.current}
                                                          onClose={() => {setOpenViewGroupDialog(false);}} />
                  }
              />
          }
      </>

  )
};