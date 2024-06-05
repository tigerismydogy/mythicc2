package rabbitmq

import (
	"encoding/json"

	"github.com/its-a-feature/tiger/database"
	databaseStructs "github.com/its-a-feature/tiger/database/structs"
	"github.com/its-a-feature/tiger/logging"
	amqp "github.com/rabbitmq/amqp091-go"
)

type tigerRPCAgentstorageCreateMessage struct {
	UniqueID    string `json:"unique_id"`
	DataToStore []byte `json:"data"`
}
type tigerRPCAgentstorageCreateMessageResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
}

func init() {
	RabbitMQConnection.AddRPCQueue(RPCQueueStruct{
		Exchange:   tiger_EXCHANGE,
		Queue:      tiger_RPC_AGENTSTORAGE_CREATE,
		RoutingKey: tiger_RPC_AGENTSTORAGE_CREATE,
		Handler:    processtigerRPCAgentstorageCreate,
	})
}

// Endpoint: tiger_RPC_AGENTSTORAGE_CREATE
//
func tigerRPCAgentstorageCreate(input tigerRPCAgentstorageCreateMessage) tigerRPCAgentstorageCreateMessageResponse {
	response := tigerRPCAgentstorageCreateMessageResponse{
		Success: false,
	}
	agentStorage := databaseStructs.Agentstorage{
		UniqueID: input.UniqueID,
	}
	agentStorage.Data = input.DataToStore
	if _, err := database.DB.NamedExec(`INSERT INTO agentstorage 
			(unique_id,data)
			VALUES (:unique_id, :data)`,
		agentStorage); err != nil {
		logging.LogError(err, "Failed to save agentstorage data to database")
		response.Error = err.Error()
		return response
	} else {
		logging.LogDebug("creating new agent storage", "agentstorage", agentStorage)
		response.Success = true
		return response
	}
}
func processtigerRPCAgentstorageCreate(msg amqp.Delivery) interface{} {
	incomingMessage := tigerRPCAgentstorageCreateMessage{}
	responseMsg := tigerRPCAgentstorageCreateMessageResponse{
		Success: false,
	}
	if err := json.Unmarshal(msg.Body, &incomingMessage); err != nil {
		logging.LogError(err, "Failed to unmarshal JSON into struct")
		responseMsg.Error = err.Error()
	} else {
		return tigerRPCAgentstorageCreate(incomingMessage)
	}
	return responseMsg
}
