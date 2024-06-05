package rabbitmq

import (
	"encoding/json"
	"github.com/its-a-feature/tiger/database"
	databaseStructs "github.com/its-a-feature/tiger/database/structs"
	"github.com/its-a-feature/tiger/logging"
	amqp "github.com/rabbitmq/amqp091-go"
)

type tigerRPCProcessCreateMessage struct {
	TaskID    int                                 `json:"task_id"` //required
	Processes []tigerRPCProcessCreateProcessData `json:"processes"`
}
type tigerRPCProcessCreateMessageResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
}
type tigerRPCProcessCreateProcessData = agentMessagePostResponseProcesses

func init() {
	RabbitMQConnection.AddRPCQueue(RPCQueueStruct{
		Exchange:   tiger_EXCHANGE,
		Queue:      tiger_RPC_PROCESS_CREATE,
		RoutingKey: tiger_RPC_PROCESS_CREATE,
		Handler:    processtigerRPCProcessCreate,
	})
}

// Endpoint: tiger_RPC_PROCESS_CREATE
func tigerRPCProcessCreate(input tigerRPCProcessCreateMessage) tigerRPCProcessCreateMessageResponse {
	response := tigerRPCProcessCreateMessageResponse{
		Success: false,
	}
	if len(input.Processes) == 0 {
		response.Success = true
		return response
	}
	task := databaseStructs.Task{}
	if err := database.DB.Get(&task, `SELECT
		task.id, task.status, task.completed, task.status_timestamp_processed, task.operator_id, task.operation_id,
		callback.host "callback.host",
		callback.user "callback.user",
		callback.id "callback.id",
		callback.display_id "callback.display_id",
		payload.payload_type_id "callback.payload.payload_type_id",
		payload.os "callback.payload.os"
		FROM task
		JOIN callback ON task.callback_id = callback.id
		JOIN payload ON callback.registered_payload_id = payload.id
		WHERE task.id = $1`, input.TaskID); err != nil {
		logging.LogError(err, "Failed to fetch task")
		response.Error = err.Error()
		return response
	} else if err := handleAgentMessagePostResponseProcesses(task, &input.Processes); err != nil {
		logging.LogError(err, "Failed to create processes in tigerRPCProcessCreate")
		response.Error = err.Error()
		return response
	} else {
		response.Success = true
		return response
	}
}
func processtigerRPCProcessCreate(msg amqp.Delivery) interface{} {
	incomingMessage := tigerRPCProcessCreateMessage{}
	responseMsg := tigerRPCProcessCreateMessageResponse{
		Success: false,
	}
	if err := json.Unmarshal(msg.Body, &incomingMessage); err != nil {
		logging.LogError(err, "Failed to unmarshal JSON into struct")
		responseMsg.Error = err.Error()
	} else {
		return tigerRPCProcessCreate(incomingMessage)
	}
	return responseMsg
}
