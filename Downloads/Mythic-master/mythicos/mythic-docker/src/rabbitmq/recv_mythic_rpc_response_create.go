package rabbitmq

import (
	"encoding/json"

	"github.com/its-a-feature/tiger/database"
	databaseStructs "github.com/its-a-feature/tiger/database/structs"
	"github.com/its-a-feature/tiger/logging"
	amqp "github.com/rabbitmq/amqp091-go"
)

type tigerRPCResponseCreateMessage struct {
	TaskID   int    `json:"task_id"`
	Response []byte `json:"response"`
}
type tigerRPCResponseCreateMessageResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
}

func init() {
	RabbitMQConnection.AddRPCQueue(RPCQueueStruct{
		Exchange:   tiger_EXCHANGE,
		Queue:      tiger_RPC_RESPONSE_CREATE,
		RoutingKey: tiger_RPC_RESPONSE_CREATE,
		Handler:    processtigerRPCPayloadResponseCreate,
	})
}

// Endpoint: tiger_RPC_RESPONSE_CREATE
//
// Creates a FileMeta object for a specific task in tiger's database and writes contents to disk with a random UUID filename.
func tigerRPCResponseCreate(input tigerRPCResponseCreateMessage) tigerRPCResponseCreateMessageResponse {
	response := tigerRPCResponseCreateMessageResponse{
		Success: false,
	}
	databaseResponse := databaseStructs.Response{
		TaskID:   input.TaskID,
		Response: input.Response,
	}
	if len(input.Response) == 0 {
		response.Error = "Response must have actual bytes"
		return response
	}
	if err := database.DB.Get(&databaseResponse.OperationID, `SELECT operation_id FROM task 
		WHERE id=$1`, input.TaskID); err != nil {
		logging.LogError(err, "failed to fetch task from database")
		response.Error = err.Error()
		return response
	} else if _, err := database.DB.NamedExec(`INSERT INTO response 
	(task_id, response, operation_id)
	VALUES (:task_id, :response, :operation_id)`, databaseResponse); err != nil {
		logging.LogError(err, "Failed to create response for task", "response", input.Response)
		response.Error = err.Error()
		return response
	} else {
		response.Success = true
		return response
	}
}
func processtigerRPCPayloadResponseCreate(msg amqp.Delivery) interface{} {
	incomingMessage := tigerRPCResponseCreateMessage{}
	responseMsg := tigerRPCResponseCreateMessageResponse{
		Success: false,
	}
	if err := json.Unmarshal(msg.Body, &incomingMessage); err != nil {
		logging.LogError(err, "Failed to unmarshal JSON into struct")
		responseMsg.Error = err.Error()
	} else {
		return tigerRPCResponseCreate(incomingMessage)
	}
	return responseMsg
}
