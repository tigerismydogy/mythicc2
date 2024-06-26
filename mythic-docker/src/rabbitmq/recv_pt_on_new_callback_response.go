package rabbitmq

import (
	"encoding/json"
	"fmt"
	"github.com/its-a-feature/tiger/database"
	databaseStructs "github.com/its-a-feature/tiger/database/structs"
	"github.com/its-a-feature/tiger/logging"
	amqp "github.com/rabbitmq/amqp091-go"
)

// PAYLOAD_BUILD STRUCTS

type PTOnNewCallbackResponse struct {
	AgentCallbackID string `json:"agent_callback_id"`
	Success         bool   `json:"success"`
	Error           string `json:"error"`
}

func init() {
	RabbitMQConnection.AddDirectQueue(DirectQueueStruct{
		Exchange:   tiger_EXCHANGE,
		Queue:      "tiger_consume_new_callback_response",
		RoutingKey: PT_ON_NEW_CALLBACK_RESPONSE_ROUTING_KEY,
		Handler:    processOnNewCallbackResponse,
	})

}

// handle payload build response messages coming back on the queue
func processOnNewCallbackResponse(msg amqp.Delivery) {
	logging.LogInfo("got message", "routingKey", msg.RoutingKey)
	newCallbackResponse := PTOnNewCallbackResponse{}
	if err := json.Unmarshal(msg.Body, &newCallbackResponse); err != nil {
		logging.LogError(err, "Failed to process new callback response message")
	} else {
		//logging.LogInfo("got build response", "buildMsg", payloadBuildResponse)
		if newCallbackResponse.Success {
			return // don't need to send any alert messages on success
		}
		databaseCallback := databaseStructs.Callback{}
		err := database.DB.Get(&databaseCallback, `SELECT 
			callback.display_id, callback.operation_id
			FROM callback 
			WHERE agent_callback_id=$1 
			LIMIT 1`, newCallbackResponse.AgentCallbackID)
		if err != nil {
			logging.LogError(err, "Failed to get payload from the database")
			return
		}
		go SendAllOperationsMessage(fmt.Sprintf("Failed to handle new callback processing for callback %d\n%s", databaseCallback.DisplayID, newCallbackResponse.Error),
			databaseCallback.OperationID, "", database.MESSAGE_LEVEL_WARNING)
	}
}
