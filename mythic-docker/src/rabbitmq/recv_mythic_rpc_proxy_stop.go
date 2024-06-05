package rabbitmq

import (
	"encoding/json"
	"fmt"

	"github.com/its-a-feature/tiger/database"
	databaseStructs "github.com/its-a-feature/tiger/database/structs"
	"github.com/its-a-feature/tiger/logging"
	amqp "github.com/rabbitmq/amqp091-go"
)

type tigerRPCProxyStopMessage struct {
	TaskID   int    `json:"task_id"`
	Port     int    `json:"port"`
	PortType string `json:"port_type"`
}
type tigerRPCProxyStopMessageResponse struct {
	Success   bool   `json:"success"`
	Error     string `json:"error"`
	LocalPort int    `json:"local_port"`
}

func init() {
	RabbitMQConnection.AddRPCQueue(RPCQueueStruct{
		Exchange:   tiger_EXCHANGE,
		Queue:      tiger_RPC_PROXY_STOP,
		RoutingKey: tiger_RPC_PROXY_STOP,
		Handler:    processtigerRPCProxyStop,
	})
}

// Endpoint: tiger_RPC_PROXY
//
// Creates a FileMeta object for a specific task in tiger's database and writes contents to disk with a random UUID filename.
func tigerRPCProxyStop(input tigerRPCProxyStopMessage) tigerRPCProxyStopMessageResponse {
	response := tigerRPCProxyStopMessageResponse{
		Success: false,
	}
	task := databaseStructs.Task{ID: input.TaskID}
	if err := database.DB.Get(&task, `SELECT id, operation_id, callback_id FROM task WHERE id=$1`, task.ID); err != nil {
		logging.LogError(err, "Failed to get task from database to start socks")
		response.Error = err.Error()
		return response
	} else {
		switch input.PortType {
		case CALLBACK_PORT_TYPE_RPORTFWD:
			fallthrough
		case CALLBACK_PORT_TYPE_SOCKS:
			fallthrough
		case CALLBACK_PORT_TYPE_INTERACTIVE:
			if input.Port == 0 {
				// lookup the port that might need to be closed for this PortType and CallbackID
				input.Port = proxyPorts.GetPortForTypeAndCallback(task.ID, task.CallbackID, input.PortType)
				if input.Port == 0 {
					response.Error = fmt.Sprintf("Failed to find port for type, %s, and task, %d", input.PortType, task.ID)
					return response
				}
			}
		}
		response.LocalPort = input.Port
		if err := proxyPorts.Remove(task.CallbackID, input.PortType, input.Port, task.OperationID); err != nil {
			logging.LogError(err, "Failed to stop callback port")
			response.Error = err.Error()
			return response
		} else {
			response.Success = true
			return response
		}
	}

}
func processtigerRPCProxyStop(msg amqp.Delivery) interface{} {
	incomingMessage := tigerRPCProxyStopMessage{}
	responseMsg := tigerRPCProxyStopMessageResponse{
		Success: false,
	}
	if err := json.Unmarshal(msg.Body, &incomingMessage); err != nil {
		logging.LogError(err, "Failed to unmarshal JSON into struct")
		responseMsg.Error = err.Error()
	} else {
		return tigerRPCProxyStop(incomingMessage)
	}
	return responseMsg
}
