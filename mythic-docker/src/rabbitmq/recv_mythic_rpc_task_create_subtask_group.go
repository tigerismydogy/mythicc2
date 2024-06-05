package rabbitmq

import (
	"encoding/json"

	"github.com/its-a-feature/tiger/database"
	databaseStructs "github.com/its-a-feature/tiger/database/structs"
	"github.com/its-a-feature/tiger/logging"
	amqp "github.com/rabbitmq/amqp091-go"
)

type tigerRPCTaskCreateSubtaskGroupMessage struct {
	TaskID                int                                    `json:"task_id"`    // required
	GroupName             string                                 `json:"group_name"` // required
	GroupCallbackFunction *string                                `json:"group_callback_function,omitempty"`
	Tasks                 []tigerRPCTaskCreateSubtaskGroupTasks `json:"tasks"` // required

}

type tigerRPCTaskCreateSubtaskGroupTasks struct {
	SubtaskCallbackFunction *string `json:"subtask_callback_function,omitempty"`
	CommandName             string  `json:"command_name"` // required
	Params                  string  `json:"params"`       // required
	ParameterGroupName      *string `json:"parameter_group_name,omitempty"`
	Token                   *int    `json:"token,omitempty"`
}

// Every tigerRPC function call must return a response that includes the following two values
type tigerRPCTaskCreateSubtaskGroupMessageResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
	TaskIDs []int  `json:"task_ids"`
}

func init() {
	RabbitMQConnection.AddRPCQueue(RPCQueueStruct{
		Exchange:   tiger_EXCHANGE,
		Queue:      tiger_RPC_TASK_CREATE_SUBTASK_GROUP,   // swap out with queue in rabbitmq.constants.go file
		RoutingKey: tiger_RPC_TASK_CREATE_SUBTASK_GROUP,   // swap out with routing key in rabbitmq.constants.go file
		Handler:    processtigerRPCTaskCreateSubtaskGroup, // points to function that takes in amqp.Delivery and returns interface{}
	})
}

// tiger_RPC_OBJECT_ACTION - Say what the function does
func tigerRPCTaskCreateSubtaskGroup(input tigerRPCTaskCreateSubtaskGroupMessage) tigerRPCTaskCreateSubtaskGroupMessageResponse {
	response := tigerRPCTaskCreateSubtaskGroupMessageResponse{
		Success: false,
	}
	createdSubTasks := []int{}
	parentTask := databaseStructs.Task{}
	taskingLocation := "tiger_rpc"
	operatorOperation := databaseStructs.Operatoroperation{}
	if err := database.DB.Get(&parentTask, `SELECT 
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
	`, parentTask.Operator.ID, parentTask.Callback.OperationID); err != nil {
		logging.LogError(err, "Failed to get operation information when creating subtask")
		response.Error = err.Error()
		return response
	} else {
		for _, task := range input.Tasks {
			createTaskInput := CreateTaskInput{
				ParentTaskID:            &input.TaskID,
				CommandName:             task.CommandName,
				Params:                  task.Params,
				Token:                   task.Token,
				ParameterGroupName:      task.ParameterGroupName,
				SubtaskCallbackFunction: task.SubtaskCallbackFunction,
				GroupName:               &input.GroupName,
				GroupCallbackFunction:   input.GroupCallbackFunction,
				IsOperatorAdmin:         parentTask.Operator.Admin,
				CallbackDisplayID:       parentTask.Callback.DisplayID,
				CurrentOperationID:      parentTask.Callback.OperationID,
				OperatorID:              parentTask.Operator.ID,
				TaskingLocation:         &taskingLocation,
			}
			if operatorOperation.BaseDisabledCommandsID.Valid {
				baseDisabledCommandsID := int(operatorOperation.BaseDisabledCommandsID.Int64)
				createTaskInput.DisabledCommandID = &baseDisabledCommandsID
			}
			// create a subtask of this task
			creationResponse := CreateTask(createTaskInput)
			if creationResponse.Status == "success" {
				createdSubTasks = append(createdSubTasks, creationResponse.TaskID)
			} else {
				response.Error = creationResponse.Error
				return response
			}
		}
		response.Success = true
		response.TaskIDs = createdSubTasks
		return response
	}

}
func processtigerRPCTaskCreateSubtaskGroup(msg amqp.Delivery) interface{} {
	incomingMessage := tigerRPCTaskCreateSubtaskGroupMessage{}
	responseMsg := tigerRPCTaskCreateSubtaskGroupMessageResponse{
		Success: false,
	}
	if err := json.Unmarshal(msg.Body, &incomingMessage); err != nil {
		logging.LogError(err, "Failed to unmarshal JSON into struct")
		responseMsg.Error = err.Error()
	} else {
		return tigerRPCTaskCreateSubtaskGroup(incomingMessage)
	}
	return responseMsg
}
