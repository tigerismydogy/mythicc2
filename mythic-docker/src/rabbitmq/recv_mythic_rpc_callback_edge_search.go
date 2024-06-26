package rabbitmq

import (
	"encoding/json"
	"github.com/mitchellh/mapstructure"
	"time"

	"github.com/its-a-feature/tiger/database"
	databaseStructs "github.com/its-a-feature/tiger/database/structs"
	"github.com/its-a-feature/tiger/logging"
	amqp "github.com/rabbitmq/amqp091-go"
)

type tigerRPCCallbackEdgeSearchMessage struct {
	AgentCallbackUUID     string  `json:"agent_callback_id"`
	AgentCallbackID       int     `json:"callback_id"`
	SearchC2ProfileName   *string `json:"search_c2_profile_name"`
	SearchActiveEdgesOnly *bool   `json:"search_active_edges_only"`
}
type tigerRPCCallbackEdgeSearchMessageResult struct {
	ID             int                      `mapstructure:"id" json:"id"`
	StartTimestamp time.Time                `mapstructure:"start_timestamp" json:"start_timestamp"`
	EndTimestamp   time.Time                `mapstructure:"end_timestamp" json:"end_timestamp"`
	Source         databaseStructs.Callback `mapstructure:"source" json:"source"`
	Destination    databaseStructs.Callback `mapstructure:"destination" json:"destination"`
	C2Profile      string                   `mapstructure:"c2profile_name" json:"c2profile"`
}
type tigerRPCCallbackEdgeSearchMessageResponse struct {
	Success bool                                       `json:"success"`
	Error   string                                     `json:"error"`
	Results []tigerRPCCallbackEdgeSearchMessageResult `json:"results"`
}

func init() {
	RabbitMQConnection.AddRPCQueue(RPCQueueStruct{
		Exchange:   tiger_EXCHANGE,
		Queue:      tiger_RPC_CALLBACK_EDGE_SEARCH,
		RoutingKey: tiger_RPC_CALLBACK_EDGE_SEARCH,
		Handler:    processtigerRPCCallbackEdgeSearch,
	})
}

// Endpoint: tiger_RPC_CALLBACK_EDGE_SEARCH
func tigerRPCCallbackEdgeSearch(input tigerRPCCallbackEdgeSearchMessage) tigerRPCCallbackEdgeSearchMessageResponse {
	response := tigerRPCCallbackEdgeSearchMessageResponse{
		Success: false,
	}

	searchString := `SELECT 
    		callbackgraphedge.*
			FROM callbackgraphedge
			JOIN callback s on callbackgraphedge.source_id = s.id
			JOIN callback d on callbackgraphedge.destination_id = d.id
			WHERE s.id != d.id AND (s.agent_callback_id=:agent_callback_id OR s.id=:id OR 
			      d.agent_callback_id=:agent_callback_id OR d.id=:id)`

	rows, err := database.DB.NamedQuery(searchString, databaseStructs.Callback{
		ID:              input.AgentCallbackID,
		AgentCallbackID: input.AgentCallbackUUID,
	})
	if err != nil {
		logging.LogError(err, "Failed to search callback edge information")
		response.Error = err.Error()
		return response
	}
	searchResults := databaseStructs.Callbackgraphedge{}
	for rows.Next() {
		result := tigerRPCCallbackEdgeSearchMessageResult{}
		err = rows.StructScan(&searchResults)
		if err != nil {
			logging.LogError(err, "Failed to get row from callbacks for search")
			continue
		}
		if input.SearchActiveEdgesOnly != nil && *input.SearchActiveEdgesOnly && searchResults.EndTimestamp.Valid {
			// want to only search for active edges and this one has a valid end timestamp
			continue
		}
		if input.SearchC2ProfileName != nil && *input.SearchC2ProfileName != searchResults.C2Profile.Name {
			// looking for a specific c2 profile edge and this isn't it
			continue
		}
		err = mapstructure.Decode(searchResults, &result)
		if err != nil {
			logging.LogError(err, "Failed to map callback search results into array")
			response.Error = err.Error()
			return response
		}
		result.C2Profile = searchResults.C2Profile.Name
		response.Results = append(response.Results, result)
	}
	response.Success = true
	return response

}
func processtigerRPCCallbackEdgeSearch(msg amqp.Delivery) interface{} {
	incomingMessage := tigerRPCCallbackEdgeSearchMessage{}
	responseMsg := tigerRPCCallbackEdgeSearchMessageResponse{
		Success: false,
	}
	if err := json.Unmarshal(msg.Body, &incomingMessage); err != nil {
		logging.LogError(err, "Failed to unmarshal JSON into struct")
		responseMsg.Error = err.Error()
	} else {
		return tigerRPCCallbackEdgeSearch(incomingMessage)
	}
	return responseMsg
}
