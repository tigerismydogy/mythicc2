package rabbitmq

import (
	"encoding/json"
	"github.com/its-a-feature/tiger/database"
	databaseStructs "github.com/its-a-feature/tiger/database/structs"
	"github.com/its-a-feature/tiger/logging"
	amqp "github.com/rabbitmq/amqp091-go"
)

type tigerRPCTokenRemoveMessage struct {
	TaskID int                             `json:"task_id"` //required
	Tokens []tigerRPCTokenRemoveTokenData `json:"tokens"`
}
type tigerRPCTokenRemoveMessageResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
}
type tigerRPCTokenRemoveTokenData = agentMessagePostResponseToken

func init() {
	RabbitMQConnection.AddRPCQueue(RPCQueueStruct{
		Exchange:   tiger_EXCHANGE,
		Queue:      tiger_RPC_TOKEN_REMOVE,
		RoutingKey: tiger_RPC_TOKEN_REMOVE,
		Handler:    processtigerRPCTokenRemove,
	})
}

// Endpoint: tiger_RPC_TOKEN_REMOVE
func tigerRPCTokenRemove(input tigerRPCTokenRemoveMessage) tigerRPCTokenRemoveMessageResponse {
	response := tigerRPCTokenRemoveMessageResponse{
		Success: false,
	}
	if len(input.Tokens) == 0 {
		response.Success = true
		return response
	}
	for i := 0; i < len(input.Tokens); i++ {
		input.Tokens[i].Action = "remove"
	}
	task := databaseStructs.Task{}
	if err := database.DB.Get(&task, `SELECT
	callback.operation_id "callback.operation_id",
	callback.host "callback.host"
	FROM task
	JOIN callback ON task.callback_id = callback.id
	WHERE task.id = $1`, input.TaskID); err != nil {
		logging.LogError(err, "Failed to fetch task")
		response.Error = err.Error()
		return response
	} else if err := handleAgentMessagePostResponseTokens(task, &input.Tokens); err != nil {
		logging.LogError(err, "Failed to create processes in tigerRPCProcessCreate")
		response.Error = err.Error()
		return response
	} else {
		response.Success = true
		return response
	}
}
func processtigerRPCTokenRemove(msg amqp.Delivery) interface{} {
	incomingMessage := tigerRPCTokenRemoveMessage{}
	responseMsg := tigerRPCTokenRemoveMessageResponse{
		Success: false,
	}
	if err := json.Unmarshal(msg.Body, &incomingMessage); err != nil {
		logging.LogError(err, "Failed to unmarshal JSON into struct")
		responseMsg.Error = err.Error()
	} else {
		return tigerRPCTokenRemove(incomingMessage)
	}
	return responseMsg
}
