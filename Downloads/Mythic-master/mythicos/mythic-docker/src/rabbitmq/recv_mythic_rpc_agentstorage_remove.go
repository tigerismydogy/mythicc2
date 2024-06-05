package rabbitmq

import (
	"encoding/json"

	"github.com/its-a-feature/tiger/database"
	databaseStructs "github.com/its-a-feature/tiger/database/structs"
	"github.com/its-a-feature/tiger/logging"
	amqp "github.com/rabbitmq/amqp091-go"
)

type tigerRPCAgentstorageRemoveMessage struct {
	UniqueID string `json:"unique_id"`
}
type tigerRPCAgentstorageRemoveMessageResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
}

func init() {
	RabbitMQConnection.AddRPCQueue(RPCQueueStruct{
		Exchange:   tiger_EXCHANGE,
		Queue:      tiger_RPC_AGENTSTORAGE_REMOVE,
		RoutingKey: tiger_RPC_AGENTSTORAGE_REMOVE,
		Handler:    processtigerRPCAgentstorageRemove,
	})
}

// Endpoint: tiger_RPC_AGENTSTORAGE_REMOVE
//
func tigerRPCAgentstorageRemove(input tigerRPCAgentstorageRemoveMessage) tigerRPCAgentstorageRemoveMessageResponse {
	response := tigerRPCAgentstorageRemoveMessageResponse{
		Success: false,
	}
	agentStorage := databaseStructs.Agentstorage{
		UniqueID: input.UniqueID,
	}
	if _, err := database.DB.NamedExec(`DELETE FROM agentstorage 
			WHERE unique_id=:unique_id`,
		agentStorage); err != nil {
		logging.LogError(err, "Failed to save agentstorage data to database")
		response.Error = err.Error()
		return response
	} else {
		logging.LogDebug("Removed agentstorage entries", "unique_id", input.UniqueID)
		response.Success = true
		return response
	}
}
func processtigerRPCAgentstorageRemove(msg amqp.Delivery) interface{} {
	incomingMessage := tigerRPCAgentstorageRemoveMessage{}
	responseMsg := tigerRPCAgentstorageRemoveMessageResponse{
		Success: false,
	}
	if err := json.Unmarshal(msg.Body, &incomingMessage); err != nil {
		logging.LogError(err, "Failed to unmarshal JSON into struct")
		responseMsg.Error = err.Error()
	} else {
		return tigerRPCAgentstorageRemove(incomingMessage)
	}
	return responseMsg
}
