package rabbitmq

import (
	"database/sql"
	"encoding/json"

	"github.com/its-a-feature/tiger/database"
	databaseStructs "github.com/its-a-feature/tiger/database/structs"
	"github.com/its-a-feature/tiger/logging"
	amqp "github.com/rabbitmq/amqp091-go"
)

type tigerRPCCallbackRemoveCommandMessage struct {
	TaskID   int      `json:"task_id"`  // required
	Commands []string `json:"commands"` // required
}
type tigerRPCCallbackRemoveCommandMessageResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
}

func init() {
	RabbitMQConnection.AddRPCQueue(RPCQueueStruct{
		Exchange:   tiger_EXCHANGE,
		Queue:      tiger_RPC_CALLBACK_REMOVE_COMMAND,
		RoutingKey: tiger_RPC_CALLBACK_REMOVE_COMMAND,
		Handler:    processtigerRPCCallbackAddCommand,
	})
}

// Endpoint: tiger_RPC_CALLBACK_REMOVE_COMMAND
func tigerRPCCallbackRemoveCommand(input tigerRPCCallbackRemoveCommandMessage) tigerRPCCallbackRemoveCommandMessageResponse {
	response := tigerRPCCallbackRemoveCommandMessageResponse{
		Success: false,
	}
	task := databaseStructs.Task{}
	if err := database.DB.Get(&task, `SELECT
		task.operator_id,
		callback.id "callback.id",
		payload.payload_type_id "payload.payload_type_id"
		FROM task
		JOIN callback on task.callback_id = callback.id
		JOIN payload on callback.registered_payload_id = payload.id
		WHERE task.id=$1`, input.TaskID); err != nil {
		logging.LogError(err, "Failed to find task in tigerRPCCallbackRemoveCommand")
		response.Error = err.Error()
		return response
	} else {
		if err := CallbackRemoveCommand(task.Callback.ID, task.Callback.Payload.PayloadTypeID, task.OperatorID, input.Commands); err != nil {
			logging.LogError(err, "Failed to remove commands to callback")
			response.Error = err.Error()
			return response
		} else {
			response.Success = true
			return response
		}
	}
}
func CallbackRemoveCommand(callbackID int, payloadtypeID int, operatorID int, commands []string) error {
	for _, command := range commands {
		// first check if the command is already loaded
		databaseCommand := databaseStructs.Command{}
		loadedCommand := databaseStructs.Loadedcommands{}
		if err := database.DB.Get(&databaseCommand, `SELECT
		id, "version"
		FROM command
		WHERE command.cmd=$1 AND command.payload_type_id=$2`,
			command, payloadtypeID); err != nil {
			logging.LogError(err, "Failed to find command to load")
			return err
		} else if err := database.DB.Get(&loadedCommand, `SELECT id
		FROM loadedcommands
		WHERE command_id=$1 AND callback_id=$2`,
			databaseCommand.ID, callbackID); err == nil {
			// the command is loaded, so remove it
			if _, err := database.DB.NamedExec(`DELETE FROM loadedcommands WHERE id=:id`, loadedCommand); err != nil {
				logging.LogError(err, "Failed to remove command from callback")
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
func processtigerRPCCallbackRemoveCommand(msg amqp.Delivery) interface{} {
	incomingMessage := tigerRPCCallbackRemoveCommandMessage{}
	responseMsg := tigerRPCCallbackRemoveCommandMessageResponse{
		Success: false,
	}
	if err := json.Unmarshal(msg.Body, &incomingMessage); err != nil {
		logging.LogError(err, "Failed to unmarshal JSON into struct")
		responseMsg.Error = err.Error()
	} else {
		return tigerRPCCallbackRemoveCommand(incomingMessage)
	}
	return responseMsg
}
