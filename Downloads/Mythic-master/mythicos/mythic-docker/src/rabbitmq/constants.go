package rabbitmq

import "time"

const (
	tiger_EXCHANGE                        = "tiger_exchange"
	tiger_TOPIC_EXCHANGE                  = "tiger_topic_exchange"
	RETRY_CONNECT_DELAY                    = 5 * time.Second
	CHECK_CONTAINER_STATUS_DELAY           = 10 * time.Second
	TIME_FORMAT_STRING_YYYY_MM_DD          = "2006-01-02"
	TIME_FORMAT_STRING_YYYY_MM_DD_HH_MM_SS = "2006-01-02 15:04:05 Z07"
	RPC_TIMEOUT                            = 20 * time.Second
	TASK_STATUS_CONTAINER_DOWN             = "Error: Container Down"
)

// Direct fanout rabbitmq routes where tiger is consuming messages, but others can also listen in and consume
const (
	// PT_SYNC_ROUTING_KEY payload routes
	//	Syncing information about the payload type to tiger
	//		send PayloadTypeSyncMessage to this route
	PT_SYNC_ROUTING_KEY       = "pt_sync"
	PT_RPC_RESYNC_ROUTING_KEY = "pt_rpc_resync"
	//	Result of asking a container to build a new payload
	//		send PayloadBuildResponse to this route
	PT_BUILD_RESPONSE_ROUTING_KEY = "pt_build_response"
	// Result of informing a container of a new callback based on its payload type
	// 		send PTOnNewCallbackResponse to this route
	PT_ON_NEW_CALLBACK_RESPONSE_ROUTING_KEY = "pt_on_new_callback_response"
	//	Result of asking a container to build a new c2profile-only payload for hot-swapping c2s
	//		send PayloadBuildC2Response to this route
	PT_BUILD_C2_RESPONSE_ROUTING_KEY = "pt_c2_build_response"
	//	Result of asking a container to do a pre-flight check for a task
	// 		send PTTTaskOPSECPreTaskMessageResponse to this route
	PT_TASK_OPSEC_PRE_CHECK_RESPONSE = "pt_task_opsec_pre_check_response"
	//	Result of asking a container to process a tasking request
	// 		send PTTaskCreateTaskingMessageResponse to this route
	PT_TASK_CREATE_TASKING_RESPONSE = "pt_task_create_tasking_response"
	//	Result of asking a container to do a post-flight check for a task before an agent picks it up
	//		send PTTaskOPSECPostTaskMessageResponse to this route
	PT_TASK_OPSEC_POST_CHECK_RESPONSE = "pt_task_opsec_post_check_response"
	//	Result of handling a task's completion function
	//		send PTTaskCompletionHandlerMessageResponse to this route
	PT_TASK_COMPLETION_FUNCTION_RESPONSE = "pt_task_completion_function_response"
	PT_TASK_PROCESS_RESPONSE_RESPONSE    = "pt_task_process_response_response"

	// C2_SYNC_ROUTING_KEY c2 routes
	//		send C2SyncMessages to this route
	C2_SYNC_ROUTING_KEY       = "c2_sync"
	C2_RPC_RESYNC_ROUTING_KEY = "c2_rpc_resync"
	// TR_SYNC_ROUTING_KEY
	TR_SYNC_ROUTING_KEY       = "tr_sync"
	TR_RPC_RESYNC_ROUTING_KEY = "tr_rpc_resync"
)

// Direct fanout rabbitmq routes where the container is consuming messages and responding back to tiger, but others can also listen in and consume.
// These aren't RPC routes because these could take a long time and we don't want to block it. These will have the format of "containerName_"
// prepended to the constants below
//
//	ex: for the apfell agent, it would be "apfell_payload_build"
const (
	//	send PayloadBuildMessage to this route for agent containers to pick up
	//		send PayloadBuildResponse to PAYLOAD_BUILD_RESPONSE_ROUTING_KEY for tiger to process response
	//
	PT_BUILD_ROUTING_KEY = "payload_build"
	//
	PT_BUILD_C2_ROUTING_KEY = "payload_c2_build"
	//
	PT_TASK_OPSEC_PRE_CHECK = "pt_task_opsec_pre_check"
	//
	PT_TASK_CREATE_TASKING   = "pt_task_create_tasking"
	PT_ON_NEW_CALLBACK       = "pt_on_new_callback"
	PT_TASK_OPSEC_POST_CHECK = "pt_task_opsec_post_check"
	//
	PT_RPC_COMMAND_DYNAMIC_QUERY_FUNCTION = "pt_command_dynamic_query_function"
	//
	PT_RPC_COMMAND_TYPEDARRAY_PARSE = "pt_command_typedarray_parse"
	//
	PT_TASK_COMPLETION_FUNCTION = "pt_task_completion_function"
	//
	PT_TASK_PROCESS_RESPONSE = "pt_task_process_response"
)

// Routes where container is consuming messages and responding back to tiger
//
//	These are exclusive to the container and not able for other containers to listen in on
//	These all have "containerName_" prepended to the constants below
//		ex: for the http profile, it would be "http_c2_rpc_opsec_check"
const (
	//
	C2_RPC_OPSEC_CHECKS_ROUTING_KEY = "c2_rpc_opsec_check"
	//
	C2_RPC_CONFIG_CHECK_ROUTING_KEY = "c2_rpc_config_check"
	//
	C2_RPC_GET_IOC_ROUTING_KEY = "c2_rpc_get_ioc"
	//
	C2_RPC_SAMPLE_MESSAGE_ROUTING_KEY = "c2_rpc_sample_message"
	//
	C2_RPC_REDIRECTOR_RULES_ROUTING_KEY = "c2_rpc_redirector_rules"
	//
	C2_RPC_START_SERVER_ROUTING_KEY = "c2_rpc_start_server"
	//
	C2_RPC_STOP_SERVER_ROUTING_KEY = "c2_rpc_stop_server"
	//
	C2_RPC_GET_SERVER_DEBUG_OUTPUT = "c2_rpc_get_server_debug_output"
	//
	C2_RPC_HOST_FILE = "c2_rpc_host_file"
	//
	C2_RPC_GET_FILE = "c2_rpc_get_file"
	//
	C2_RPC_REMOVE_FILE = "c2_rpc_remove_file"
	//
	C2_RPC_LIST_FILE = "c2_rpc_list_file"
	//
	C2_RPC_WRITE_FILE = "c2_rpc_write_file"
	//
	TR_RPC_GENERATE_KEYS = "tr_rpc_generate_keys"
	//
	TR_RPC_CONVERT_FROM_tiger_C2_FORMAT = "tr_rpc_from_tiger_c2"
	//
	TR_RPC_CONVERT_TO_tiger_C2_FORMAT = "tr_rpc_to_tiger_c2"
	//
	TR_RPC_ENCRYPT_BYTES = "tr_rpc_encrypt_bytes"
	//
	TR_RPC_DECRYPT_BYTES = "tr_rpc_decrypt_bytes"
)

// RPC Routes where tiger is consuming messages and responding back to the container
//
//	These are exclusive to tiger and not able for other containers to listen in on
const (
	// tiger_RPC_FILE_CREATE file operations
	tiger_RPC_FILE_CREATE      = "tiger_rpc_file_create"
	tiger_RPC_FILE_SEARCH      = "tiger_rpc_file_search"
	tiger_RPC_FILE_UPDATE      = "tiger_rpc_file_update"
	tiger_RPC_FILE_GET_CONTENT = "tiger_rpc_file_get_content"
	tiger_RPC_FILE_REGISTER    = "tiger_rpc_file_register"
	// tiger_RPC_PAYLOAD_CREATE_FROM_UUID payload operations
	tiger_RPC_PAYLOAD_CREATE_FROM_UUID    = "tiger_rpc_payload_create_from_uuid"
	tiger_RPC_PAYLOAD_CREATE_FROM_SCRATCH = "tiger_rpc_payload_create_from_scratch"
	tiger_RPC_PAYLOAD_SEARCH              = "tiger_rpc_payload_search"
	tiger_RPC_PAYLOAD_GET_PAYLOAD_CONTENT = "tiger_rpc_payload_get_content"
	tiger_RPC_PAYLOAD_UPDATE_BUILD_STEP   = "tiger_rpc_payload_update_build_step"
	tiger_RPC_PAYLOAD_ADD_COMMAND         = "tiger_rpc_payload_add_command"
	tiger_RPC_PAYLOAD_REMOVE_COMMAND      = "tiger_rpc_payload_remove_command"
	// tiger_RPC_TASK_SEARCH task operations
	tiger_RPC_TASK_SEARCH                    = "tiger_rpc_task_search"
	tiger_RPC_TASK_DISPLAY_TO_REAL_ID_SEARCH = "tiger_rpc_task_display_to_real_id_search"
	tiger_RPC_TASK_UPDATE                    = "tiger_rpc_task_update"
	tiger_RPC_TASK_CREATE                    = "tiger_rpc_task_create"
	tiger_RPC_TASK_CREATE_SUBTASK            = "tiger_rpc_task_create_subtask"
	tiger_RPC_TASK_CREATE_SUBTASK_GROUP      = "tiger_rpc_task_create_group"
	// tiger_RPC_RESPONSE_SEARCH response operations
	tiger_RPC_RESPONSE_SEARCH = "tiger_rpc_response_search"
	tiger_RPC_RESPONSE_CREATE = "tiger_rpc_response_create"
	// tiger_RPC_COMMAND_SEARCH command operations
	tiger_RPC_COMMAND_SEARCH = "tiger_rpc_command_search"
	// tiger_RPC_CALLBACK_CREATE callback operations
	tiger_RPC_CALLBACK_CREATE                    = "tiger_rpc_callback_create"
	tiger_RPC_CALLBACK_SEARCH                    = "tiger_rpc_callback_search"
	tiger_RPC_CALLBACK_EDGE_SEARCH               = "tiger_rpc_callback_edge_search"
	tiger_RPC_CALLBACK_DISPLAY_TO_REAL_ID_SEARCH = "tiger_rpc_callback_display_to_real_id_search"
	tiger_RPC_CALLBACK_ADD_COMMAND               = "tiger_rpc_callback_add_command"
	tiger_RPC_CALLBACK_REMOVE_COMMAND            = "tiger_rpc_callback_remove_command"
	tiger_RPC_CALLBACK_SEARCH_COMMAND            = "tiger_rpc_callback_search_command"
	tiger_RPC_CALLBACK_UPDATE                    = "tiger_rpc_callback_update"
	tiger_RPC_CALLBACK_ENCRYPT_BYTES             = "tiger_rpc_callback_encrypt_bytes"
	tiger_RPC_CALLBACK_DECRYPT_BYTES             = "tiger_rpc_callback_decrypt_bytes"
	// tiger_RPC_AGENTSTORAGE_CREATE agent storage operations
	tiger_RPC_AGENTSTORAGE_CREATE = "tiger_rpc_agentstorage_create"
	tiger_RPC_AGENTSTORAGE_SEARCH = "tiger_rpc_agentstorage_search"
	tiger_RPC_AGENTSTORAGE_REMOVE = "tiger_rpc_agentstorage_remove"
	// tiger_RPC_PROCESS_CREATE process operations
	tiger_RPC_PROCESS_CREATE = "tiger_rpc_process_create"
	tiger_RPC_PROCESS_SEARCH = "tiger_rpc_process_search"
	// tiger_RPC_ARTIFACT_CREATE artifact operations
	tiger_RPC_ARTIFACT_CREATE = "tiger_rpc_artifact_create"
	tiger_RPC_ARTIFACT_SEARCH = "tiger_rpc_artifact_search"
	// tiger_RPC_KEYLOG_CREATE keylog operations
	tiger_RPC_KEYLOG_CREATE = "tiger_rpc_keylog_create"
	tiger_RPC_KEYLOG_SEARCH = "tiger_rpc_keylog_search"
	// tiger_RPC_CREDENTIAL_CREATE credential operations
	tiger_RPC_CREDENTIAL_CREATE = "tiger_rpc_credential_create"
	tiger_RPC_CREDENTIAL_SEARCH = "tiger_rpc_credential_search"
	// tiger_RPC_EVENTLOG_CREATE event log operations
	tiger_RPC_EVENTLOG_CREATE = "tiger_rpc_eventlog_create"
	// tiger_RPC_FILEBROWSER_CREATE filebrowser operations
	tiger_RPC_FILEBROWSER_CREATE = "tiger_rpc_filebrowser_create"
	tiger_RPC_FILEBROWSER_REMOVE = "tiger_rpc_filebrowser_remove"
	// tiger_RPC_PAYLOADONHOST_CREATE payload on host operations
	tiger_RPC_PAYLOADONHOST_CREATE = "tiger_rpc_payloadonhost_create"
	// tiger_RPC_CALLBACKTOKEN_CREATE callback token operations
	tiger_RPC_CALLBACKTOKEN_CREATE = "tiger_rpc_callbacktoken_create"
	tiger_RPC_CALLBACKTOKEN_REMOVE = "tiger_rpc_callbacktoken_remove"
	// tiger_RPC_TOKEN_CREATE token operations
	tiger_RPC_TOKEN_CREATE = "tiger_rpc_token_create"
	tiger_RPC_TOKEN_REMOVE = "tiger_rpc_token_remove"
	// tiger_RPC_PROXY_START proxy operations
	tiger_RPC_PROXY_START = "tiger_rpc_proxy_start"
	tiger_RPC_PROXY_STOP  = "tiger_rpc_proxy_stop"
	// tiger_RPC_BLANK blank
	tiger_RPC_BLANK = "tiger_rpc_blank"
)
