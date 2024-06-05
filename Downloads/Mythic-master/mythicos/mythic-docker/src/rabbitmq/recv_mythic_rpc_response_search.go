package rabbitmq

import (
	"encoding/json"

	"github.com/its-a-feature/tiger/database"
	databaseStructs "github.com/its-a-feature/tiger/database/structs"
	"github.com/its-a-feature/tiger/logging"
	amqp "github.com/rabbitmq/amqp091-go"
)

type tigerRPCResponseSearchMessage struct {
	TaskID   int    `json:"task_id"`
	Response string `json:"response"`
}
type tigerRPCResponseSearchMessageResponse struct {
	Success   bool                `json:"success"`
	Error     string              `json:"error"`
	Responses []tigerRPCResponse `json:"responses"`
}
type tigerRPCResponse struct {
	ResponseID int    `json:"response_id"`
	Response   []byte `json:"response"`
	TaskID     int    `json:"task_id"`
}

func init() {
	RabbitMQConnection.AddRPCQueue(RPCQueueStruct{
		Exchange:   tiger_EXCHANGE,
		Queue:      tiger_RPC_RESPONSE_SEARCH,
		RoutingKey: tiger_RPC_RESPONSE_SEARCH,
		Handler:    processtigerRPCPayloadResponseSearch,
	})
}

// Endpoint: tiger_RPC_RESPONSE_SEARCH
//
// Creates a FileMeta object for a specific task in tiger's database and writes contents to disk with a random UUID filename.
func tigerRPCResponseSearch(input tigerRPCResponseSearchMessage) tigerRPCResponseSearchMessageResponse {
	response := tigerRPCResponseSearchMessageResponse{
		Success: false,
	}
	results := []databaseStructs.Response{}
	responseSearch := input.Response
	if responseSearch == "" {
		responseSearch = "%_%"
	} else {
		responseSearch = "%" + responseSearch + "%"
	}
	task := databaseStructs.Task{}
	if err := database.DB.Get(&task, `SELECT 
	operation_id, display_id, id
	FROM task
	WHERE id=$1`, input.TaskID); err != nil {
		logging.LogError(err, "Failed to fetch task for tigerRPCResponseSearch")
		response.Error = err.Error()
		return response
	} else if err := database.DB.Select(&results, `SELECT 
	id, response, task_id 
	FROM response
	WHERE operation_id=$1 AND response LIKE $2 AND task_id=$3`,
		task.OperationID, responseSearch, input.TaskID); err != nil {
		logging.LogError(err, "Failed to fetch responses for tigerRPCResponseSearch")
		response.Error = err.Error()
		return response
	} else {
		for _, resp := range results {
			response.Responses = append(response.Responses, tigerRPCResponse{
				ResponseID: resp.ID,
				Response:   resp.Response,
				TaskID:     resp.TaskID,
			})
		}
		response.Success = true
		return response
	}
}
func processtigerRPCPayloadResponseSearch(msg amqp.Delivery) interface{} {
	incomingMessage := tigerRPCResponseSearchMessage{}
	responseMsg := tigerRPCResponseSearchMessageResponse{
		Success: false,
	}
	if err := json.Unmarshal(msg.Body, &incomingMessage); err != nil {
		logging.LogError(err, "Failed to unmarshal JSON into struct")
		responseMsg.Error = err.Error()
	} else {
		return tigerRPCResponseSearch(incomingMessage)
	}
	return responseMsg
}
