package rabbitmq

import (
	"encoding/json"
	"fmt"
	"github.com/mitchellh/mapstructure"

	"github.com/its-a-feature/tiger/database"
	databaseStructs "github.com/its-a-feature/tiger/database/structs"
	"github.com/its-a-feature/tiger/logging"
	amqp "github.com/rabbitmq/amqp091-go"
)

type tigerRPCArtifactSearchMessage struct {
	TaskID          int                                 `json:"task_id"` //required
	SearchArtifacts tigerRPCArtifactSearchArtifactData `json:"artifact"`
}
type tigerRPCArtifactSearchMessageResponse struct {
	Success   bool                                  `json:"success"`
	Error     string                                `json:"error"`
	Artifacts []tigerRPCArtifactSearchArtifactData `json:"artifacts"`
}
type tigerRPCArtifactSearchArtifactData struct {
	Host            *string `json:"host" `            // optional
	ArtifactType    *string `json:"artifact_type"`    //optional
	ArtifactMessage *string `json:"artifact_message"` //optional
	TaskID          *int    `json:"task_id"`          //optional
}

func init() {
	RabbitMQConnection.AddRPCQueue(RPCQueueStruct{
		Exchange:   tiger_EXCHANGE,
		Queue:      tiger_RPC_ARTIFACT_SEARCH,
		RoutingKey: tiger_RPC_ARTIFACT_SEARCH,
		Handler:    processtigerRPCArtifactSearch,
	})
}

// Endpoint: tiger_RPC_PROCESS_SEARCH
func tigerRPCArtifactSearch(input tigerRPCArtifactSearchMessage) tigerRPCArtifactSearchMessageResponse {
	response := tigerRPCArtifactSearchMessageResponse{
		Success:   false,
		Artifacts: []tigerRPCArtifactSearchArtifactData{},
	}
	paramDict := make(map[string]interface{})
	task := databaseStructs.Task{}
	if err := database.DB.Get(&task, `SELECT 
	task.id,
	callback.operation_id "callback.operation_id"
	FROM task
	JOIN callback ON task.callback_id = callback.id
	WHERE task.id=$1`, input.TaskID); err != nil {
		logging.LogError(err, "Failed to search for task information when searching for artifacts")
		response.Error = err.Error()
		return response
	} else {
		paramDict["operation_id"] = task.Callback.OperationID
		searchString := `SELECT * FROM taskartifact WHERE operation_id=:operation_id `
		if input.SearchArtifacts.Host != nil {
			paramDict["host"] = fmt.Sprintf("%%%s%%", *input.SearchArtifacts.Host)
			searchString += "AND host ILIKE :host "
		}
		if input.SearchArtifacts.ArtifactMessage != nil {
			paramDict["artifact"] = fmt.Sprintf("%%%s%%", *input.SearchArtifacts.ArtifactMessage)
			searchString += "AND artifact LIKE :artifact "
		}
		if input.SearchArtifacts.TaskID != nil {
			paramDict["task_id"] = *input.SearchArtifacts.TaskID
			searchString += "AND task_id=:task_id "
		}
		if input.SearchArtifacts.ArtifactType != nil {
			paramDict["base_artifact"] = fmt.Sprintf("%%%s%%", *input.SearchArtifacts.ArtifactType)
			searchString += "AND base_artifact LIKE :base_artifact "
		}
		if rows, err := database.DB.NamedQuery(searchString, paramDict); err != nil {
			logging.LogError(err, "Failed to search artifact information")
			response.Error = err.Error()
			return response
		} else {
			for rows.Next() {
				result := tigerRPCArtifactSearchArtifactData{}
				searchResult := databaseStructs.Taskartifact{}
				if err = rows.StructScan(&searchResult); err != nil {
					logging.LogError(err, "Failed to get row from artifacts for search")
				} else if err = mapstructure.Decode(searchResult, &result); err != nil {
					logging.LogError(err, "Failed to map artifact search results into array")
				} else {
					response.Artifacts = append(response.Artifacts, result)
				}
			}
			response.Success = true
			return response
		}
	}
}
func processtigerRPCArtifactSearch(msg amqp.Delivery) interface{} {
	incomingMessage := tigerRPCArtifactSearchMessage{}
	responseMsg := tigerRPCArtifactSearchMessageResponse{
		Success: false,
	}
	if err := json.Unmarshal(msg.Body, &incomingMessage); err != nil {
		logging.LogError(err, "Failed to unmarshal JSON into struct")
		responseMsg.Error = err.Error()
	} else {
		return tigerRPCArtifactSearch(incomingMessage)
	}
	return responseMsg
}
