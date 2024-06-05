package rabbitmq

import (
	"encoding/json"
	"github.com/mitchellh/mapstructure"

	"github.com/its-a-feature/tiger/database"
	databaseStructs "github.com/its-a-feature/tiger/database/structs"
	"github.com/its-a-feature/tiger/logging"
	amqp "github.com/rabbitmq/amqp091-go"
)

type tigerRPCKeylogSearchMessage struct {
	TaskID        int                             `json:"task_id"` //required
	SearchKeylogs tigerRPCKeylogSearchKeylogData `json:"keylogs"`
}
type tigerRPCKeylogSearchMessageResponse struct {
	Success bool                              `json:"success"`
	Error   string                            `json:"error"`
	Keylogs []tigerRPCKeylogSearchKeylogData `json:"keylogs"`
}
type tigerRPCKeylogSearchKeylogData struct {
	User        *string `json:"user" mapstructure:"user"`             // optional
	WindowTitle *string `json:"window_title" mapstructure:"window"`   // optional
	Keystrokes  *[]byte `json:"keystrokes" mapstructure:"keystrokes"` // optional
}

func init() {
	RabbitMQConnection.AddRPCQueue(RPCQueueStruct{
		Exchange:   tiger_EXCHANGE,
		Queue:      tiger_RPC_KEYLOG_SEARCH,
		RoutingKey: tiger_RPC_KEYLOG_SEARCH,
		Handler:    processtigerRPCKeylogSearch,
	})
}

// Endpoint: tiger_RPC_KEYLOG_SEARCH
func tigerRPCKeylogSearch(input tigerRPCKeylogSearchMessage) tigerRPCKeylogSearchMessageResponse {
	response := tigerRPCKeylogSearchMessageResponse{
		Success: false,
	}
	paramDict := make(map[string]interface{})
	task := databaseStructs.Task{}
	if err := database.DB.Get(&task, `SELECT 
	task.id,
	callback.operation_id "callback.operation_id",
	callback.host "callback.host"
	FROM task
	JOIN callback ON task.callback_id = callback.id
	WHERE task.id=$1`, input.TaskID); err != nil {
		response.Error = err.Error()
		return response
	} else {
		keylogs := []databaseStructs.Keylog{}
		paramDict["operation_id"] = task.Callback.OperationID
		searchString := `SELECT * FROM keylog WHERE operation_id=:operation_id  `
		if input.SearchKeylogs.User != nil {
			paramDict["user"] = *input.SearchKeylogs.User
			searchString += "AND user ILIKE %:user% "
		}
		if input.SearchKeylogs.WindowTitle != nil {
			paramDict["window_title"] = *input.SearchKeylogs.WindowTitle
			searchString += "AND window ILIKE %:window_title% "
		}
		if input.SearchKeylogs.Keystrokes != nil {
			paramDict["keystrokes"] = *input.SearchKeylogs.Keystrokes
			searchString += "AND keystrokes LIKE %:keystrokes% "
		}

		if err := database.DB.Select(&keylogs, searchString, paramDict); err != nil {
			response.Error = err.Error()
			return response
		} else {
			returnedProcesses := []tigerRPCKeylogSearchKeylogData{}
			if err := mapstructure.Decode(keylogs, &returnedProcesses); err != nil {
				response.Error = err.Error()
				return response
			} else {
				response.Success = true
				response.Keylogs = returnedProcesses
				return response
			}

		}
	}
}
func processtigerRPCKeylogSearch(msg amqp.Delivery) interface{} {
	incomingMessage := tigerRPCKeylogSearchMessage{}
	responseMsg := tigerRPCKeylogSearchMessageResponse{
		Success: false,
	}
	if err := json.Unmarshal(msg.Body, &incomingMessage); err != nil {
		logging.LogError(err, "Failed to unmarshal JSON into struct")
		responseMsg.Error = err.Error()
	} else {
		return tigerRPCKeylogSearch(incomingMessage)
	}
	return responseMsg
}
