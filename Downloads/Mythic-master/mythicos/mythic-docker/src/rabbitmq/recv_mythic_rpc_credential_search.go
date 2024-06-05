package rabbitmq

import (
	"encoding/json"
	"fmt"

	"github.com/its-a-feature/tiger/database"
	databaseStructs "github.com/its-a-feature/tiger/database/structs"
	"github.com/its-a-feature/tiger/logging"
	amqp "github.com/rabbitmq/amqp091-go"
)

type tigerRPCCredentialSearchMessage struct {
	TaskID            int                                     `json:"task_id"` //required
	SearchCredentials tigerRPCCredentialSearchCredentialData `json:"credentials"`
}
type tigerRPCCredentialSearchMessageResponse struct {
	Success     bool                                      `json:"success"`
	Error       string                                    `json:"error"`
	Credentials []tigerRPCCredentialSearchCredentialData `json:"credentials"`
}
type tigerRPCCredentialSearchCredentialData struct {
	Type       *string `json:"type" `      // optional
	Account    *string `json:"account" `   // optional
	Realm      *string `json:"realm" `     // optional
	Credential *string `json:"credential"` // optional
	Comment    *string `json:"comment"`    // optional
	Metadata   *string `json:"metadata"`   // optional
}

func init() {
	RabbitMQConnection.AddRPCQueue(RPCQueueStruct{
		Exchange:   tiger_EXCHANGE,
		Queue:      tiger_RPC_CREDENTIAL_SEARCH,
		RoutingKey: tiger_RPC_CREDENTIAL_SEARCH,
		Handler:    processtigerRPCCredentialSearch,
	})
}

// Endpoint: tiger_RPC_CREDENTIAL_SEARCH
func tigerRPCCredentialSearch(input tigerRPCCredentialSearchMessage) tigerRPCCredentialSearchMessageResponse {
	response := tigerRPCCredentialSearchMessageResponse{
		Success: false,
	}
	params := []interface{}{}
	task := databaseStructs.Task{}
	if err := database.DB.Get(&task, `SELECT 
	task.id,
	task.operation_id
	FROM task
	WHERE task.id=$1`, input.TaskID); err != nil {
		response.Error = err.Error()
		return response
	} else {
		credentials := []databaseStructs.Credential{}
		params = append(params, task.OperationID)
		searchString := fmt.Sprintf(`SELECT * FROM credential WHERE operation_id=$%d `, len(params))

		if input.SearchCredentials.Type != nil {
			params = append(params, "%"+*input.SearchCredentials.Type+"%")
			searchString += fmt.Sprintf("AND \"type\" ILIKE $%d ", len(params))
		}
		if input.SearchCredentials.Credential != nil {
			params = append(params, "%"+*input.SearchCredentials.Credential+"%")
			searchString += fmt.Sprintf("AND credential ILIKE $%d ", len(params))
		}
		if input.SearchCredentials.Account != nil {
			params = append(params, "%"+*input.SearchCredentials.Account+"%")
			searchString += fmt.Sprintf("AND account ILIKE $%d ", len(params))
		}
		if input.SearchCredentials.Realm != nil {
			params = append(params, "%"+*input.SearchCredentials.Realm+"%")
			searchString += fmt.Sprintf("AND realm ILIKE $%d ", len(params))
		}
		if input.SearchCredentials.Comment != nil {
			params = append(params, "%"+*input.SearchCredentials.Comment+"%")
			searchString += fmt.Sprintf("AND comment ILIKE $%d ", len(params))
		}
		if input.SearchCredentials.Metadata != nil {
			params = append(params, "%"+*input.SearchCredentials.Metadata+"%")
			searchString += fmt.Sprintf("AND metadata ILIKE $%d ", len(params))
		}
		if err := database.DB.Select(&credentials, searchString, params...); err != nil {
			logging.LogError(err, "Failed to search for credentials")
			response.Error = err.Error()
			return response
		} else {
			returnedCredentials := []tigerRPCCredentialSearchCredentialData{}
			if marshalledBytes, err := json.Marshal(credentials); err != nil {
				logging.LogError(err, "Failed to marshal credential result")
				response.Error = err.Error()
				return response
			} else if err := json.Unmarshal(marshalledBytes, &returnedCredentials); err != nil {
				logging.LogError(err, "Failed to unmarshal credential results")
				response.Error = err.Error()
				return response
			} else {
				response.Success = true
				response.Credentials = returnedCredentials
				return response
			}

		}
	}
}
func processtigerRPCCredentialSearch(msg amqp.Delivery) interface{} {
	incomingMessage := tigerRPCCredentialSearchMessage{}
	responseMsg := tigerRPCCredentialSearchMessageResponse{
		Success: false,
	}
	if err := json.Unmarshal(msg.Body, &incomingMessage); err != nil {
		logging.LogError(err, "Failed to unmarshal JSON into struct")
		responseMsg.Error = err.Error()
	} else {
		return tigerRPCCredentialSearch(incomingMessage)
	}
	return responseMsg
}
