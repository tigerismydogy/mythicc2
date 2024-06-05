package rabbitmq

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/its-a-feature/tiger/database"
	databaseStructs "github.com/its-a-feature/tiger/database/structs"
	"github.com/its-a-feature/tiger/logging"
	amqp "github.com/rabbitmq/amqp091-go"
)

type tigerRPCCallbackAddCommandMessage struct {
	TaskID   int      `json:"task_id"`  // required
	Commands []string `json:"commands"` // required
}
type tigerRPCCallbackAddCommandMessageResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
}

func init() {
	RabbitMQConnection.AddRPCQueue(RPCQueueStruct{
		Exchange:   tiger_EXCHANGE,
		Queue:      tiger_RPC_CALLBACK_ADD_COMMAND,
		RoutingKey: tiger_RPC_CALLBACK_ADD_COMMAND,
		Handler:    processtigerRPCCallbackAddCommand,
	})
}

// Endpoint: tiger_RPC_CALLBACK_ADD_COMMAND
func tigerRPCCallbackAddCommand(input tigerRPCCallbackAddCommandMessage) tigerRPCCallbackAddCommandMessageResponse {
	response := tigerRPCCallbackAddCommandMessageResponse{
		Success: false,
	}
	task := databaseStructs.Task{}
	if err := database.DB.Get(&task, `SELECT
		task.operator_id,
		callback.id "callback.id",
		payload.payload_type_id "callback.payload.payload_type_id"
		FROM task
		JOIN callback on task.callback_id = callback.id
		JOIN payload on callback.registered_payload_id = payload.id
		WHERE task.id=$1`, input.TaskID); err != nil {
		logging.LogError(err, "Failed to find task in tigerRPCCallbackAddCommand")
		response.Error = err.Error()
		return response
	} else {
		if err := CallbackAddCommand(task.Callback.ID, task.Callback.Payload.PayloadTypeID, task.OperatorID, input.Commands); err != nil {
			logging.LogError(err, "Failed to add commands to callback")
			response.Error = err.Error()
			return response
		} else {
			response.Success = true
			return response
		}
	}

}
func CallbackAddCommand(callbackID int, payloadtypeID int, operatorID int, commands []string) error {
	for _, command := range commands {
		logging.LogDebug("trying to add command", "cmd", command)
		// first check if the command is already loaded
		// if not, try to add it as a loaded command
		databaseCommand := databaseStructs.Command{}
		loadedCommand := databaseStructs.Loadedcommands{}
		if err := database.DB.Get(&databaseCommand, `SELECT
		id, "version"
		FROM command
		WHERE command.cmd=$1 AND command.payload_type_id=$2`,
			command, payloadtypeID); err != nil {
			logging.LogError(err, "Failed to find command to load")
			return errors.New("Failed to find command: " + command)
		} else if err := database.DB.Get(&loadedCommand, `SELECT id
		FROM loadedcommands
		WHERE command_id=$1 AND callback_id=$2`,
			databaseCommand.ID, callbackID); err == nil {
			continue
		} else if err == sql.ErrNoRows {
			// this never existed, so let's add it as a loaded command
			loadedCommand.Version = databaseCommand.Version
			loadedCommand.CallbackID = callbackID
			loadedCommand.CommandID = databaseCommand.ID
			loadedCommand.OperatorID = operatorID
			if _, err := database.DB.NamedExec(`INSERT INTO loadedcommands 
			("version", callback_id, command_id, operator_id)
			VALUES (:version, :callback_id, :command_id, :operator_id)`,
				loadedCommand); err != nil {
				logging.LogError(err, "Failed to mark command as loaded in callback")
				return err
			} else {
				continue
			}
		} else {
			// we got some other sort of error
			logging.LogError(err, "Failed to query database for loaded command")
			return err
		}
	}
	return nil
}
func processtigerRPCCallbackAddCommand(msg amqp.Delivery) interface{} {
	incomingMessage := tigerRPCCallbackAddCommandMessage{}
	responseMsg := tigerRPCCallbackAddCommandMessageResponse{
		Success: false,
	}
	if err := json.Unmarshal(msg.Body, &incomingMessage); err != nil {
		logging.LogError(err, "Failed to unmarshal JSON into struct")
		responseMsg.Error = err.Error()
	} else {
		return tigerRPCCallbackAddCommand(incomingMessage)
	}
	return responseMsg
}
