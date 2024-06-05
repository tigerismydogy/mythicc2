package rabbitmq

import (
	"encoding/json"

	"github.com/its-a-feature/tiger/database"
	databaseStructs "github.com/its-a-feature/tiger/database/structs"
	"github.com/its-a-feature/tiger/logging"
	amqp "github.com/rabbitmq/amqp091-go"
)

type tigerRPCTaskCreateSubtaskMessage struct {
	TaskID                  int     `json:"task_id"`
	SubtaskCallbackFunction *string `json:"subtask_callback_function,omitempty"`
	CommandName             string  `json:"command_name"`
	Params                  string  `json:"params"`
	ParameterGroupName      *string `json:"parameter_group_name,omitempty"`
	Token                   *int    `json:"token,omitempty"`
}

// Every tigerRPC function call must return a response that includes the following two values
type tigerRPCTaskCreateSubtaskMessageResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
	TaskID  int    `json:"task_id"`
}

func init() {
	RabbitMQConnection.AddRPCQueue(RPCQueueStruct{
		Exchange:   tiger_EXCHANGE,
		Queue:      tiger_RPC_TASK_CREATE_SUBTASK,    // swap out with queue in rabbitmq.constants.go file
		RoutingKey: tiger_RPC_TASK_CREATE_SUBTASK,    // swap out with routing key in rabbitmq.constants.go file
		Handler:    processtigerRPCTaskCreateSubtask, // points to function that takes in amqp.Delivery and returns interface{}
	})
}

// tiger_RPC_OBJECT_ACTION - Say what the function does
func tigerRPCTaskCreateSubtask(input tigerRPCTaskCreateSubtaskMessage) tigerRPCTaskCreateSubtaskMessageResponse {
	response := tigerRPCTaskCreateSubtaskMessageResponse{
		Success: false,
	}
	taskingLocation := "tiger_rpc"
	createTaskInput := CreateTaskInput{
		ParentTaskID:            &input.TaskID,
		CommandName:             input.CommandName,
		Params:                  input.Params,
		Token:                   input.Token,
		ParameterGroupName:      input.ParameterGroupName,
		SubtaskCallbackFunction: input.SubtaskCallbackFunction,
		TaskingLocation:         &taskingLocation,
	}
	task := databaseStructs.Task{}
	operatorOperation := databaseStructs.Operatoroperation{}
	if err := database.DB.Get(&task, `SELECT 
	callback.id "callback.id",
	callback.display_id "callback.display_id",
	callback.operation_id "callback.operation_id",
	operator.id "operator.id",
	operator.admin "operator.admin" 
	FROM task
	JOIN callback ON task.callback_id = callback.id 
	JOIN operator ON task.operator_id = operator.id
	WHERE task.id=$1`, input.TaskID); err != nil {
		response.Error = err.Error()
		logging.LogError(err, "Failed to fetch task/callback information when creating subtask")
		return response
	} else if err := database.DB.Get(&operatorOperation, `SELECT
	base_disabled_commands_id
	FROM operatoroperation
	WHERE operator_id = $1 AND operation_id = $2
	`, task.Operator.ID, task.Callback.OperationID); err != nil {
		logging.LogError(err, "Failed to get operation information when creating subtask")
		response.Error = err.Error()
		return response
	} else {
		createTaskInput.IsOperatorAdmin = task.Operator.Admin
		createTaskInput.CallbackDisplayID = task.Callback.DisplayID
		createTaskInput.CurrentOperationID = task.Callback.OperationID
		if operatorOperation.BaseDisabledCommandsID.Valid {
			baseDisabledCommandsID := int(operatorOperation.BaseDisabledCommandsID.Int64)
			createTaskInput.DisabledCommandID = &baseDisabledCommandsID
		}
		createTaskInput.OperatorID = task.Operator.ID
		// create a subtask of this task
		creationResponse := CreateTask(createTaskInput)
		if creationResponse.Status == "success" {
			response.Success = true
			response.TaskID = creationResponse.TaskID
		} else {
			response.Error = creationResponse.Error
		}
		return response
	}

}
func processtigerRPCTaskCreateSubtask(msg amqp.Delivery) interface{} {
	incomingMessage := tigerRPCTaskCreateSubtaskMessage{}
	responseMsg := tigerRPCTaskCreateSubtaskMessageResponse{
		Success: false,
	}
	if err := json.Unmarshal(msg.Body, &incomingMessage); err != nil {
		logging.LogError(err, "Failed to unmarshal JSON into struct")
		responseMsg.Error = err.Error()
	} else {
		return tigerRPCTaskCreateSubtask(incomingMessage)
	}
	return responseMsg
}
