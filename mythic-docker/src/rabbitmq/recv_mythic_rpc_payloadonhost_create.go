package rabbitmq

import (
	"encoding/json"
	"fmt"
	"github.com/its-a-feature/tiger/database"
	databaseStructs "github.com/its-a-feature/tiger/database/structs"
	"github.com/its-a-feature/tiger/logging"
	amqp "github.com/rabbitmq/amqp091-go"
	"strings"
)

type tigerRPCPayloadOnHostCreateMessage struct {
	TaskID        int                              `json:"task_id"` //required
	PayloadOnHost tigerRPCPayloadOnHostCreateData `json:"payload_on_host"`
}
type tigerRPCPayloadOnHostCreateMessageResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
}
type tigerRPCPayloadOnHostCreateData struct {
	Host        string  `json:"host"`
	PayloadId   *int    `json:"payload_id"`
	PayloadUUID *string `json:"payload_uuid"`
}

func init() {
	RabbitMQConnection.AddRPCQueue(RPCQueueStruct{
		Exchange:   tiger_EXCHANGE,
		Queue:      tiger_RPC_PAYLOADONHOST_CREATE,
		RoutingKey: tiger_RPC_PAYLOADONHOST_CREATE,
		Handler:    processtigerRPCPayloadOnHostCreate,
	})
}

// Endpoint: tiger_RPC_PAYLOADONHOST_CREATE
func tigerRPCPayloadOnHostCreate(input tigerRPCPayloadOnHostCreateMessage) tigerRPCPayloadOnHostCreateMessageResponse {
	response := tigerRPCPayloadOnHostCreateMessageResponse{
		Success: false,
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
	} else {
		if input.PayloadOnHost.Host == "" {
			input.PayloadOnHost.Host = task.Callback.Host
		} else {
			input.PayloadOnHost.Host = strings.ToUpper(input.PayloadOnHost.Host)
		}
		payloadId := 0
		if input.PayloadOnHost.PayloadId != nil {
			payloadId = *input.PayloadOnHost.PayloadId
		} else if input.PayloadOnHost.PayloadUUID != nil {
			payload := databaseStructs.Payload{}
			if err := database.DB.Get(&payload, `SELECT id FROM payload WHERE uuid=$1 AND operation_id=$2`,
				*input.PayloadOnHost.PayloadUUID, task.OperationID); err != nil {
				logging.LogError(err, "Failed to find specified payload UUID")
				response.Error = err.Error()
				return response
			} else {
				payloadId = payload.ID
			}
		}
		if payloadId == 0 {
			response.Error = fmt.Sprintf("Failed to find the specified payload")
			return response
		} else if _, err := database.DB.Exec(`INSERT INTO payloadonhost 
			(host, payload_id, operation_id, task_id) VALUES 
			($1, $2, $3, $4)`,
			input.PayloadOnHost.Host, payloadId, task.OperationID, task.ID); err != nil {
			logging.LogError(err, "Failed to create new payload on host value")
			response.Error = err.Error()
			return response
		} else {
			response.Success = true
			return response
		}
	}
}
func processtigerRPCPayloadOnHostCreate(msg amqp.Delivery) interface{} {
	incomingMessage := tigerRPCPayloadOnHostCreateMessage{}
	responseMsg := tigerRPCPayloadOnHostCreateMessageResponse{
		Success: false,
	}
	if err := json.Unmarshal(msg.Body, &incomingMessage); err != nil {
		logging.LogError(err, "Failed to unmarshal JSON into struct")
		responseMsg.Error = err.Error()
	} else {
		return tigerRPCPayloadOnHostCreate(incomingMessage)
	}
	return responseMsg
}
