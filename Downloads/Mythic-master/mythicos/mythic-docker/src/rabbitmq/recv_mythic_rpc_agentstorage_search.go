package rabbitmq

import (
	"encoding/json"
	"fmt"

	"github.com/its-a-feature/tiger/database"
	databaseStructs "github.com/its-a-feature/tiger/database/structs"
	"github.com/its-a-feature/tiger/logging"
	amqp "github.com/rabbitmq/amqp091-go"
)

type tigerRPCAgentstorageSearchMessage struct {
	SearchUniqueID string `json:"unique_id" db:"unique_id"` // required
}
type tigerRPCAgentstorageSearchMessageResponse struct {
	Success              bool                                `json:"success"`
	Error                string                              `json:"error"`
	AgentStorageMessages []tigerRPCAgentstorageSearchResult `json:"agentstorage_messages"`
}

type tigerRPCAgentstorageSearchResult struct {
	UniqueID string `json:"unique_id"`
	Data     []byte `json:"data"`
}

func init() {
	RabbitMQConnection.AddRPCQueue(RPCQueueStruct{
		Exchange:   tiger_EXCHANGE,
		Queue:      tiger_RPC_AGENTSTORAGE_SEARCH,
		RoutingKey: tiger_RPC_AGENTSTORAGE_SEARCH,
		Handler:    processtigerRPCAgentstorageSearch,
	})
}

// Endpoint: tiger_RPC_AGENTSTORAGE_SEARCH
func tigerRPCAgentstorageSearch(input tigerRPCAgentstorageSearchMessage) tigerRPCAgentstorageSearchMessageResponse {
	response := tigerRPCAgentstorageSearchMessageResponse{
		Success: false,
	}
	agentStorageMessages := []databaseStructs.Agentstorage{}
	searchUniqueID := fmt.Sprintf("%%%s%%", input.SearchUniqueID)
	if err := database.DB.Select(&agentStorageMessages, `SELECT
	*
	FROM agentstorage
	WHERE unique_id ILIKE $1`, searchUniqueID); err != nil {
		logging.LogError(err, "Failed to search agentstorage data")
		response.Error = err.Error()
		return response
	} else {
		response.Success = true
		agentStorageResponses := make([]tigerRPCAgentstorageSearchResult, len(agentStorageMessages))
		for i, msg := range agentStorageMessages {
			agentStorageResponses[i] = tigerRPCAgentstorageSearchResult{
				UniqueID: msg.UniqueID,
				Data:     msg.Data,
			}
		}
		response.AgentStorageMessages = agentStorageResponses
		return response
	}
}
func processtigerRPCAgentstorageSearch(msg amqp.Delivery) interface{} {
	incomingMessage := tigerRPCAgentstorageSearchMessage{}
	responseMsg := tigerRPCAgentstorageSearchMessageResponse{
		Success: false,
	}
	if err := json.Unmarshal(msg.Body, &incomingMessage); err != nil {
		logging.LogError(err, "Failed to unmarshal JSON into struct")
		responseMsg.Error = err.Error()
	} else {
		return tigerRPCAgentstorageSearch(incomingMessage)
	}
	return responseMsg
}
