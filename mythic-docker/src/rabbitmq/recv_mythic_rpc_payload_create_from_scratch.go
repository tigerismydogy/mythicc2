package rabbitmq

import (
	"encoding/json"
	"strings"

	"github.com/its-a-feature/tiger/database"
	databaseStructs "github.com/its-a-feature/tiger/database/structs"
	"github.com/its-a-feature/tiger/logging"
	amqp "github.com/rabbitmq/amqp091-go"
)

type tigerRPCPayloadCreateFromScratchMessage struct {
	TaskID               int                  `json:"task_id"`
	PayloadConfiguration PayloadConfiguration `json:"payload_configuration"`
	RemoteHost           *string              `json:"remote_host"`
}

// Every tigerRPC function call must return a response that includes the following two values
type tigerRPCPayloadCreateFromScratchMessageResponse struct {
	Success        bool   `json:"success"`
	Error          string `json:"error"`
	NewPayloadUUID string `json:"new_payload_uuid"`
}

func init() {
	RabbitMQConnection.AddRPCQueue(RPCQueueStruct{
		Exchange:   tiger_EXCHANGE,
		Queue:      tiger_RPC_PAYLOAD_CREATE_FROM_SCRATCH,   // swap out with queue in rabbitmq.constants.go file
		RoutingKey: tiger_RPC_PAYLOAD_CREATE_FROM_SCRATCH,   // swap out with routing key in rabbitmq.constants.go file
		Handler:    processtigerRPCPayloadCreateFromScratch, // points to function that takes in amqp.Delivery and returns interface{}
	})
}

// tiger_RPC_OBJECT_ACTION - Say what the function does
func tigerRPCPayloadCreateFromScratch(input tigerRPCPayloadCreateFromScratchMessage) tigerRPCPayloadCreateFromScratchMessageResponse {
	response := tigerRPCPayloadCreateFromScratchMessageResponse{
		Success: false,
	}
	task := databaseStructs.Task{}
	if err := database.DB.Get(&task, `SELECT 
    operator_id, operation_id
	FROM task 
	WHERE id=$1`, input.TaskID); err != nil {
		logging.LogError(err, "Failed to get operator_id from task when generating payload")
		response.Error = err.Error()
		return response
	}
	newUUID, newID, err := RegisterNewPayload(input.PayloadConfiguration, &databaseStructs.Operatoroperation{
		CurrentOperation: databaseStructs.Operation{ID: task.OperationID},
		CurrentOperator:  databaseStructs.Operator{ID: task.OperatorID},
	})
	if err != nil {
		response.Error = err.Error()
		return response
	}
	_, err = database.DB.Exec(`UPDATE payload SET auto_generated=true WHERE id=$1`, newID)
	if err != nil {
		logging.LogError(err, "failed to update payload auto_generated status")
		response.Error = err.Error()
		return response
	}
	if input.RemoteHost != nil {
		if _, err := database.DB.Exec(`INSERT INTO payloadonhost 
					(host, payload_id, operation_id, task_id) 
					VALUES 
					($1, $2, $3, $4)`, strings.ToUpper(*input.RemoteHost), newID, task.OperationID, input.TaskID); err != nil {
			logging.LogError(err, "Failed to register payload on host")
		}
	}
	response.NewPayloadUUID = newUUID
	response.Success = true
	return response

}
func processtigerRPCPayloadCreateFromScratch(msg amqp.Delivery) interface{} {
	incomingMessage := tigerRPCPayloadCreateFromScratchMessage{}
	responseMsg := tigerRPCPayloadCreateFromScratchMessageResponse{
		Success: false,
	}
	if err := json.Unmarshal(msg.Body, &incomingMessage); err != nil {
		logging.LogError(err, "Failed to unmarshal JSON into struct")
		responseMsg.Error = err.Error()
	} else {
		return tigerRPCPayloadCreateFromScratch(incomingMessage)
	}
	return responseMsg
}
