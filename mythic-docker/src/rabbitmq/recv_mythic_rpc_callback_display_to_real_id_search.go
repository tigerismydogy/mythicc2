package rabbitmq

import (
	"encoding/json"
	"github.com/its-a-feature/tiger/database"
	databaseStructs "github.com/its-a-feature/tiger/database/structs"

	"github.com/its-a-feature/tiger/logging"
	amqp "github.com/rabbitmq/amqp091-go"
)

type tigerRPCCallbackDisplayToRealIdSearchMessage struct {
	CallbackDisplayID int     `json:"callback_display_id"`
	OperationName     *string `json:"operation_name"`
	OperationID       *int    `json:"operation_id"`
}

// Every tigerRPC function call must return a response that includes the following two values
type tigerRPCCallbackDisplayToRealIdSearchMessageResponse struct {
	Success    bool   `json:"success"`
	Error      string `json:"error"`
	CallbackID int    `json:"callback_id"`
}

func init() {
	RabbitMQConnection.AddRPCQueue(RPCQueueStruct{
		Exchange:   tiger_EXCHANGE,
		Queue:      tiger_RPC_CALLBACK_DISPLAY_TO_REAL_ID_SEARCH, // swap out with queue in rabbitmq.constants.go file
		RoutingKey: tiger_RPC_CALLBACK_DISPLAY_TO_REAL_ID_SEARCH, // swap out with routing key in rabbitmq.constants.go file
		Handler:    processtigerRPCCallbackDisplayToRealIdSearch, // points to function that takes in amqp.Delivery and returns interface{}
	})
}

// tiger_RPC_OBJECT_ACTION - Say what the function does
func tigerRPCCallbackDisplayToRealIdSearch(input tigerRPCCallbackDisplayToRealIdSearchMessage) tigerRPCCallbackDisplayToRealIdSearchMessageResponse {
	response := tigerRPCCallbackDisplayToRealIdSearchMessageResponse{
		Success: false,
	}
	searchString := ""
	if input.OperationName != nil {
		searchString = `SELECT 
    		callback.id 
			FROM 
			callback
			JOIN operation on callback.operation_id = operation.id
			WHERE callback.display_id=$1 AND operation.name=$2`
		callback := databaseStructs.Callback{}
		if err := database.DB.Get(&callback, searchString, input.CallbackDisplayID, *input.OperationName); err != nil {
			logging.LogError(err, "Failed to find task based on task id and operation name")
			response.Error = err.Error()
			return response
		} else {
			response.CallbackID = callback.ID
			response.Success = true
			return response
		}
	} else if input.OperationID != nil {
		searchString = `SELECT 
    		callback.id 
			FROM 
			callback
			WHERE callback.display_id=$1 AND callback.operation_id=$2`
		callback := databaseStructs.Callback{}
		if err := database.DB.Get(&callback, searchString, input.CallbackDisplayID, *input.OperationID); err != nil {
			logging.LogError(err, "Failed to find task based on task id and operation id")
			response.Error = err.Error()
			return response
		} else {
			response.CallbackID = callback.ID
			response.Success = true
			return response
		}
	} else {
		response.Error = "Must specify operation name or operation id"
		return response
	}
}
func processtigerRPCCallbackDisplayToRealIdSearch(msg amqp.Delivery) interface{} {
	incomingMessage := tigerRPCCallbackDisplayToRealIdSearchMessage{}
	responseMsg := tigerRPCCallbackDisplayToRealIdSearchMessageResponse{
		Success: false,
	}
	if err := json.Unmarshal(msg.Body, &incomingMessage); err != nil {
		logging.LogError(err, "Failed to unmarshal JSON into struct")
		responseMsg.Error = err.Error()
	} else {
		return tigerRPCCallbackDisplayToRealIdSearch(incomingMessage)
	}
	return responseMsg
}
