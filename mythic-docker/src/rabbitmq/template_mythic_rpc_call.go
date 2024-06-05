package rabbitmq

import (
	"encoding/json"

	"github.com/its-a-feature/tiger/logging"
	amqp "github.com/rabbitmq/amqp091-go"
)

type tigerRPCObjectActionFakeNotRealMessage struct {
}

// Every tigerRPC function call must return a response that includes the following two values
type tigerRPCObjectActionFakeNotRealMessageResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
}

func init() {
	RabbitMQConnection.AddRPCQueue(RPCQueueStruct{
		Exchange:   tiger_EXCHANGE,
		Queue:      tiger_RPC_BLANK,                        // swap out with queue in rabbitmq.constants.go file
		RoutingKey: tiger_RPC_BLANK,                        // swap out with routing key in rabbitmq.constants.go file
		Handler:    processFaketigerRPCObjectActionNotReal, // points to function that takes in amqp.Delivery and returns interface{}
	})
}

//tiger_RPC_OBJECT_ACTION - Say what the function does
func tigerRPCObjectActionFakeNotReal(input tigerRPCObjectActionFakeNotRealMessage) tigerRPCObjectActionFakeNotRealMessageResponse {
	response := tigerRPCObjectActionFakeNotRealMessageResponse{}
	return response
}
func processFaketigerRPCObjectActionNotReal(msg amqp.Delivery) interface{} {
	logging.LogDebug("got message", "routingKey", msg.RoutingKey, "message", msg)
	incomingMessage := tigerRPCObjectActionFakeNotRealMessage{}
	responseMsg := tigerRPCObjectActionFakeNotRealMessageResponse{
		Success: false,
	}
	if err := json.Unmarshal(msg.Body, &incomingMessage); err != nil {
		logging.LogError(err, "Failed to unmarshal JSON into struct")
		responseMsg.Error = err.Error()
	} else {
		return tigerRPCObjectActionFakeNotReal(incomingMessage)
	}
	return responseMsg
}
