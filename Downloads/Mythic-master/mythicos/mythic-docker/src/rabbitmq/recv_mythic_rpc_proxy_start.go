package rabbitmq

import (
	"encoding/json"
	"github.com/its-a-feature/tiger/database"
	databaseStructs "github.com/its-a-feature/tiger/database/structs"
	"github.com/its-a-feature/tiger/logging"
	amqp "github.com/rabbitmq/amqp091-go"
)

type tigerRPCProxyStartMessage struct {
	TaskID     int    `json:"task_id"`
	LocalPort  int    `json:"local_port"`
	RemotePort int    `json:"remote_port"`
	RemoteIP   string `json:"remote_ip"`
	PortType   string `json:"port_type"`
}
type tigerRPCProxyStartMessageResponse struct {
	Success   bool   `json:"success"`
	LocalPort int    `json:"local_port"`
	Error     string `json:"error"`
}

func init() {
	RabbitMQConnection.AddRPCQueue(RPCQueueStruct{
		Exchange:   tiger_EXCHANGE,
		Queue:      tiger_RPC_PROXY_START,
		RoutingKey: tiger_RPC_PROXY_START,
		Handler:    processtigerRPCProxyStart,
	})
}

// Endpoint: tiger_RPC_PROXY
//
// Creates a FileMeta object for a specific task in tiger's database and writes contents to disk with a random UUID filename.
func tigerRPCProxyStart(input tigerRPCProxyStartMessage) tigerRPCProxyStartMessageResponse {
	response := tigerRPCProxyStartMessageResponse{
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
		case CALLBACK_PORT_TYPE_SOCKS:
			fallthrough
		case CALLBACK_PORT_TYPE_INTERACTIVE:
			if input.LocalPort == 0 {
				input.LocalPort = int(proxyPorts.GetNextAvailableLocalPort())
				if input.LocalPort == 0 {
					response.Error = "No more ports available through docker, please modify your .env's configuration and restart tiger"
					return response
				}
			}
		}
		response.LocalPort = input.LocalPort
		if err := proxyPorts.Add(task.CallbackID, input.PortType, input.LocalPort, input.RemotePort, input.RemoteIP, task.ID, task.OperationID,
			0, 0, 0); err != nil {
			logging.LogError(err, "Failed to add new callback port")
			response.Error = err.Error()
			return response
		} else {
			response.Success = true
			return response
		}
	}

}
func processtigerRPCProxyStart(msg amqp.Delivery) interface{} {
	incomingMessage := tigerRPCProxyStartMessage{}
	responseMsg := tigerRPCProxyStartMessageResponse{
		Success: false,
	}
	if err := json.Unmarshal(msg.Body, &incomingMessage); err != nil {
		logging.LogError(err, "Failed to unmarshal JSON into struct")
		responseMsg.Error = err.Error()
	} else {
		return tigerRPCProxyStart(incomingMessage)
	}
	return responseMsg
}
