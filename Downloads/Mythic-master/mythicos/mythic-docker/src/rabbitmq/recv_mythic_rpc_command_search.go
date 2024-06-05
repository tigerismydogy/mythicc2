package rabbitmq

import (
	"encoding/json"
	"fmt"
	"github.com/its-a-feature/tiger/utils"
	"reflect"

	"github.com/its-a-feature/tiger/database"
	databaseStructs "github.com/its-a-feature/tiger/database/structs"
	"github.com/its-a-feature/tiger/logging"
	amqp "github.com/rabbitmq/amqp091-go"
)

type tigerRPCCommandSearchMessage struct {
	SearchCommandNames        *[]string `json:"command_names,omitempty"`
	SearchPayloadTypeName     *string   `json:"payload_type_name,omitempty"`
	SearchSupportedUIFeatures *string   `json:"supported_ui_features,omitempty"`
	SearchScriptOnly          *bool     `json:"script_only,omitempty"`
	SearchOs                  *string   `json:"os,omitempty"`
	// this is an exact match search
	SearchAttributes map[string]interface{} `json:"params,omitempty"`
}

// Every tigerRPC function call must return a response that includes the following two values
type tigerRPCCommandSearchMessageResponse struct {
	Success  bool                                `json:"success"`
	Error    string                              `json:"error"`
	Commands []tigerRPCCommandSearchCommandData `json:"commands"`
}

type tigerRPCCommandSearchCommandData struct {
	Name                string                 `json:"cmd"`
	Version             int                    `json:"version"`
	Attributes          map[string]interface{} `json:"attributes"`
	NeedsAdmin          bool                   `json:"needs_admin"`
	HelpCmd             string                 `json:"help_cmd"`
	Description         string                 `json:"description"`
	SupportedUiFeatures []string               `json:"supported_ui_features"`
	Author              string                 `json:"author"`
	ScriptOnly          bool                   `json:"script_only"`
}

func init() {
	RabbitMQConnection.AddRPCQueue(RPCQueueStruct{
		Exchange:   tiger_EXCHANGE,
		Queue:      tiger_RPC_COMMAND_SEARCH,     // swap out with queue in rabbitmq.constants.go file
		RoutingKey: tiger_RPC_COMMAND_SEARCH,     // swap out with routing key in rabbitmq.constants.go file
		Handler:    processtigerRPCCommandSearch, // points to function that takes in amqp.Delivery and returns interface{}
	})
}

func interfaceArrayToStringArray(input []interface{}) []string {
	stringArray := make([]string, len(input))
	for i := 0; i < len(input); i++ {
		switch v := input[i].(type) {
		case string:
			stringArray[i] = v
		default:
			stringArray[i] = fmt.Sprintf("%v", v)
		}
	}
	return stringArray
}

// tiger_RPC_OBJECT_ACTION - Say what the function does
func tigerRPCCommandSearch(input tigerRPCCommandSearchMessage) tigerRPCCommandSearchMessageResponse {
	response := tigerRPCCommandSearchMessageResponse{
		Success: false,
	}
	foundCommands := []tigerRPCCommandSearchCommandData{}
	commands := []databaseStructs.Command{}
	if err := database.DB.Select(&commands, `SELECT
		command.*,
		payloadtype.name "payloadtype.name"
		FROM
        command
        JOIN payloadtype on command.payload_type_id = payloadtype.id
        WHERE payloadtype.name=$1 AND command.deleted=false`, input.SearchPayloadTypeName); err != nil {
		logging.LogError(err, "Failed to search for commands")
		response.Error = err.Error()
		return response
	}
	for _, command := range commands {
		if input.SearchCommandNames != nil {
			if !utils.SliceContains(*input.SearchCommandNames, command.Cmd) {
				continue
			}
		}
		if input.SearchScriptOnly != nil {
			if command.ScriptOnly != *input.SearchScriptOnly {
				continue
			}
		}
		uiFeatures := command.SupportedUiFeatures.StructValue()
		stringUIFeatures := interfaceArrayToStringArray(uiFeatures)
		attributes := map[string]interface{}{}
		if err := command.Attributes.Unmarshal(&attributes); err != nil {
			logging.LogError(err, "Failed to get attributes from command")
			response.Error = "Failed to get attributes from command"
			return response
		}
		if input.SearchSupportedUIFeatures != nil {
			if !utils.SliceContains(stringUIFeatures, *input.SearchSupportedUIFeatures) {
				continue
			}
		}
		if input.SearchOs != nil {
			logging.LogInfo("Searching supported OS", "os", *input.SearchOs)
			attributeSupportedOS := attributes["supported_os"].([]interface{})
			supportedOS := interfaceArrayToStringArray(attributeSupportedOS)
			if len(supportedOS) != 0 {
				if !utils.SliceContains(supportedOS, *input.SearchOs) {
					continue
				}
			}
			logging.LogInfo("matched OS requirement")
		}
		if input.SearchAttributes != nil {
			matchedValues := true
			logging.LogInfo("Search attributes", "raw", input.SearchAttributes, "command attributes", attributes)
			for searchKey, searchValue := range input.SearchAttributes {
				if actualValue, ok := attributes[searchKey]; ok {
					logging.LogInfo("Searching attributes", "search value", searchValue, "search key", searchKey, "actual value", actualValue)
					switch v := actualValue.(type) {
					case []interface{}:
						actualArray := interfaceArrayToStringArray(v)
						var searchArray []string
						switch s := searchValue.(type) {
						case []interface{}:
							searchArray = interfaceArrayToStringArray(s)
						default:
							searchArray = []string{fmt.Sprintf("%v", s)}
						}
						logging.LogInfo("checking arrays", "search", searchArray, "actual", actualArray)
						// need to make sure every value in searchArray is in actualArray
						for i, _ := range searchArray {
							if !utils.SliceContains(actualArray, searchArray[i]) {
								matchedValues = false
							}
						}
					case map[string]interface{}:
						searchMap := searchValue.(map[string]interface{})
						logging.LogInfo("checking maps", "search", searchMap, "actual", v)
						for skey, sval := range searchMap {
							if aval, ok := v[skey]; !ok {
								matchedValues = false
							} else if !reflect.DeepEqual(sval, aval) {
								matchedValues = false
							}
						}
					default:
						if !reflect.DeepEqual(searchValue, v) {
							matchedValues = false
						}
					}

				} else {
					matchedValues = false
				}
			}
			if matchedValues {
				newSearchCommandData := tigerRPCCommandSearchCommandData{
					Name:                command.Cmd,
					NeedsAdmin:          command.NeedsAdmin,
					Version:             command.Version,
					HelpCmd:             command.HelpCmd,
					Description:         command.Description,
					Author:              command.Author,
					ScriptOnly:          command.ScriptOnly,
					SupportedUiFeatures: stringUIFeatures,
					Attributes:          attributes,
				}
				foundCommands = append(foundCommands, newSearchCommandData)
			}
		} else {
			newFoundCommand := tigerRPCCommandSearchCommandData{
				Name:                command.Cmd,
				NeedsAdmin:          command.NeedsAdmin,
				Version:             command.Version,
				HelpCmd:             command.HelpCmd,
				Description:         command.Description,
				Author:              command.Author,
				ScriptOnly:          command.ScriptOnly,
				SupportedUiFeatures: stringUIFeatures,
				Attributes:          attributes,
			}
			foundCommands = append(foundCommands, newFoundCommand)
		}
	}
	response.Success = true
	response.Commands = foundCommands
	return response

}
func processtigerRPCCommandSearch(msg amqp.Delivery) interface{} {
	incomingMessage := tigerRPCCommandSearchMessage{}
	responseMsg := tigerRPCCommandSearchMessageResponse{
		Success: false,
	}
	if err := json.Unmarshal(msg.Body, &incomingMessage); err != nil {
		logging.LogError(err, "Failed to unmarshal JSON into struct")
		responseMsg.Error = err.Error()
	} else {
		return tigerRPCCommandSearch(incomingMessage)
	}
	return responseMsg
}
