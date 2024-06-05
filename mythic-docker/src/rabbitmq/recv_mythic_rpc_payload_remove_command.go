package rabbitmq

import (
	"database/sql"
	"encoding/json"

	"github.com/its-a-feature/tiger/database"
	databaseStructs "github.com/its-a-feature/tiger/database/structs"
	"github.com/its-a-feature/tiger/logging"
	amqp "github.com/rabbitmq/amqp091-go"
)

type tigerRPCPayloadRemoveCommandMessage struct {
	PayloadUUID string   `json:"payload_uuid"` //required
	Commands    []string `json:"commands"`     // required
}
type tigerRPCPayloadRemoveCommandMessageResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
}

func init() {
	RabbitMQConnection.AddRPCQueue(RPCQueueStruct{
		Exchange:   tiger_EXCHANGE,
		Queue:      tiger_RPC_PAYLOAD_REMOVE_COMMAND,
		RoutingKey: tiger_RPC_PAYLOAD_REMOVE_COMMAND,
		Handler:    processtigerRPCPayloadRemoveCommand,
	})
}

// Endpoint: tiger_RPC_PAYLOAD_REMOVE_COMMAND
func tigerRPCPayloadRemoveCommand(input tigerRPCPayloadRemoveCommandMessage) tigerRPCPayloadRemoveCommandMessageResponse {
	response := tigerRPCPayloadRemoveCommandMessageResponse{
		Success: false,
	}
	payload := databaseStructs.Payload{}
	if err := database.DB.Get(&payload, `SELECT payload.id, payload.payload_type_id 
	FROM payload
	WHERE uuid=$1`, input.PayloadUUID); err != nil {
		logging.LogError(err, "Failed to fetch callback in tigerRPCPayloadRemoveCommand")
		response.Error = err.Error()
		return response
	} else {
		if err := PayloadRemoveCommand(payload.ID, payload.PayloadTypeID, input.Commands); err != nil {
			logging.LogError(err, "Failed to remove commands in payload")
			response.Error = err.Error()
			return response
		} else {
			response.Success = true
			return response
		}
	}
}
func PayloadRemoveCommand(PayloadID int, payloadtypeID int, commands []string) error {
	for _, command := range commands {
		// first check if the command is already loaded
		databaseCommand := databaseStructs.Command{}
		loadedCommand := databaseStructs.Payloadcommand{}
		if err := database.DB.Get(&databaseCommand, `SELECT
		id, "version"
		FROM command
		WHERE command.cmd=$1 AND command.payload_type_id=$2`,
			command, payloadtypeID); err != nil {
			logging.LogError(err, "Failed to find command to load")
			return err
		} else if err := database.DB.Get(&loadedCommand, `SELECT id
		FROM payloadcommand
		WHERE command_id=$1 AND payload_id=$2`,
			databaseCommand.ID, PayloadID); err == nil {
			// the command is loaded, so remove it
			if _, err := database.DB.NamedExec(`DELETE FROM payloadcommand WHERE id=:id`, loadedCommand); err != nil {
				logging.LogError(err, "Failed to remove command from payload")
				return err
			}
		} else if err == sql.ErrNoRows {
			// this never existed, so move on
			continue
		} else {
			// we got some other sort of error
			logging.LogError(err, "Failed to query database for loaded command")
			return err
		}
	}
	return nil
}
func processtigerRPCPayloadRemoveCommand(msg amqp.Delivery) interface{} {
	incomingMessage := tigerRPCPayloadRemoveCommandMessage{}
	responseMsg := tigerRPCPayloadRemoveCommandMessageResponse{
		Success: false,
	}
	if err := json.Unmarshal(msg.Body, &incomingMessage); err != nil {
		logging.LogError(err, "Failed to unmarshal JSON into struct")
		responseMsg.Error = err.Error()
	} else {
		return tigerRPCPayloadRemoveCommand(incomingMessage)
	}
	return responseMsg
}
