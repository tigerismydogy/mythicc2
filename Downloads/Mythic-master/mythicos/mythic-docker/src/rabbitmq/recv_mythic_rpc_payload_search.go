package rabbitmq

import (
	"encoding/json"

	"github.com/its-a-feature/tiger/database"
	databaseStructs "github.com/its-a-feature/tiger/database/structs"
	"github.com/its-a-feature/tiger/logging"
	"github.com/its-a-feature/tiger/utils"
	amqp "github.com/rabbitmq/amqp091-go"
)

type tigerRPCPayloadSearchMessage struct {
	CallbackID                   int                                    `json:"callback_id"`
	PayloadUUID                  string                                 `json:"uuid"`
	Description                  string                                 `json:"description"`
	Filename                     string                                 `json:"filename"`
	PayloadTypes                 []string                               `json:"payload_types"`
	IncludeAutoGeneratedPayloads bool                                   `json:"include_auto_generated"`
	BuildParameters              []tigerRPCPayloadSearchBuildParameter `json:"build_parameters"`
}

type tigerRPCPayloadSearchBuildParameter struct {
	PayloadType          string            `json:"payload_type"`
	BuildParameterValues map[string]string `json:"build_parameter_values"`
}

// Every tigerRPC function call must return a response that includes the following two values
type tigerRPCPayloadSearchMessageResponse struct {
	Success               bool                   `json:"success"`
	Error                 string                 `json:"error"`
	PayloadConfigurations []PayloadConfiguration `json:"payloads"`
}

func init() {
	RabbitMQConnection.AddRPCQueue(RPCQueueStruct{
		Exchange:   tiger_EXCHANGE,
		Queue:      tiger_RPC_PAYLOAD_SEARCH,     // swap out with queue in rabbitmq.constants.go file
		RoutingKey: tiger_RPC_PAYLOAD_SEARCH,     // swap out with routing key in rabbitmq.constants.go file
		Handler:    processtigerRPCPayloadSearch, // points to function that takes in amqp.Delivery and returns interface{}
	})
}

// tiger_RPC_OBJECT_ACTION - Say what the function does
func tigerRPCPayloadSearch(input tigerRPCPayloadSearchMessage) tigerRPCPayloadSearchMessageResponse {
	response := tigerRPCPayloadSearchMessageResponse{
		Success: false,
	}
	if input.PayloadUUID != "" {
		if config, err := getPayloadConfigFromUUID(input.PayloadUUID); err != nil {
			response.Error = err.Error()
			return response
		} else {
			response.PayloadConfigurations = append(response.PayloadConfigurations, config)
			response.Success = true
			return response
		}
	} else {
		// search payloads based on the supplied information
		callback := databaseStructs.Callback{}
		payloads := []databaseStructs.Payload{}
		if err := database.DB.Get(&callback, `SELECT operation_id FROM callback WHERE id=$1`, input.CallbackID); err != nil {
			response.Error = err.Error()
			return response
		} else if err := database.DB.Select(&payloads, `SELECT
		payload.uuid, payload.auto_generated, payload.id,
		payloadtype.name "payloadtype.name",
		payloadtype.id "payloadtype.id"
		FROM payload
		JOIN payloadtype ON payload.payload_type_id = payloadtype.id
		WHERE payload.operation_id=$1
		`, callback.OperationID); err != nil {
			response.Error = err.Error()
			return response
		} else {
			finalPayloads := []PayloadConfiguration{}
			for _, payload := range payloads {
				if payload.AutoGenerated && !input.IncludeAutoGeneratedPayloads {
					continue
				} else if len(input.PayloadTypes) > 0 && !utils.SliceContains(input.PayloadTypes, payload.Payloadtype.Name) {
					continue
				} else if len(input.BuildParameters) > 0 {
					allBuildParametersAreGood := true
					for _, buildRequirement := range input.BuildParameters {
						if buildRequirement.PayloadType == payload.Payloadtype.Name {
							// only care about checking if it's the right type
							// now we need to try to find the matching build parameter to see if the value matches
							for key, val := range buildRequirement.BuildParameterValues {
								logging.LogInfo("searching build param values", "search key", key, "search val", val)
								buildParamInstance := databaseStructs.Buildparameterinstance{}
								if err := database.DB.Get(&buildParamInstance, `
								SELECT value,
								buildparameter.name "buildparameter.name"
								FROM buildparameterinstance
								JOIN buildparameter ON buildparameterinstance.build_parameter_id = buildparameter.id
								WHERE buildparameterinstance.payload_id=$1 and buildparameter.name=$2`, payload.ID, key); err != nil {
									logging.LogError(err, "Failed to get build parameters for payload type")
									response.Error = err.Error()
									return response
								} else if buildParamInstance.Value != val {
									allBuildParametersAreGood = false
								}
							}
						}
					}
					if allBuildParametersAreGood {
						if finalPayload, err := getPayloadConfigFromUUID(payload.UuID); err != nil {
							logging.LogError(err, "Failed to get configuration for payload")
							response.Error = err.Error()
							return response
						} else {
							finalPayloads = append(finalPayloads, finalPayload)
						}
					}
				} else if finalPayload, err := getPayloadConfigFromUUID(payload.UuID); err != nil {
					logging.LogError(err, "Failed to get configuration for payload")
					response.Error = err.Error()
					return response
				} else {
					finalPayloads = append(finalPayloads, finalPayload)
				}
			}
			response.PayloadConfigurations = finalPayloads
		}
	}
	response.Success = true
	return response
}
func processtigerRPCPayloadSearch(msg amqp.Delivery) interface{} {
	incomingMessage := tigerRPCPayloadSearchMessage{}
	responseMsg := tigerRPCPayloadSearchMessageResponse{
		Success: false,
	}
	if err := json.Unmarshal(msg.Body, &incomingMessage); err != nil {
		logging.LogError(err, "Failed to unmarshal JSON into struct")
		responseMsg.Error = err.Error()
	} else {
		return tigerRPCPayloadSearch(incomingMessage)
	}
	return responseMsg
}

func getPayloadConfigFromUUID(payloadUUID string) (PayloadConfiguration, error) {
	payloadConfiguration := PayloadConfiguration{}
	payload := databaseStructs.Payload{}
	if err := database.DB.Get(&payload, `SELECT
	payload.id, payload.description, payload.uuid, payload.os, payload.wrapped_payload_id, payload.build_phase,
	payloadtype.name "payloadtype.name",
	filemeta.filename "filemeta.filename",
	filemeta.agent_file_id "filemeta.agent_file_id"
	FROM
	payload
	JOIN payloadtype ON payload.payload_type_id = payloadtype.id
	JOIN filemeta ON payload.file_id = filemeta.id
	WHERE 
	payload.uuid=$1`, payloadUUID); err != nil {
		logging.LogError(err, "Failed to get payload when searching for payloads")
		return payloadConfiguration, err
	} else {
		payloadConfiguration.Description = payload.Description
		payloadConfiguration.SelectedOS = payload.Os
		payloadConfiguration.PayloadType = payload.Payloadtype.Name
		payloadConfiguration.C2Profiles = GetPayloadC2ProfileInformation(payload)
		payloadConfiguration.BuildParameters = GetBuildParameterInformation(payload.ID)
		payloadConfiguration.Commands = GetPayloadCommandInformation(payload)
		payloadConfiguration.Filename = string(payload.Filemeta.Filename)
		payloadConfiguration.AgentFileID = payload.Filemeta.AgentFileID
		payloadConfiguration.UUID = payload.UuID
		payloadConfiguration.BuildPhase = payload.BuildPhase
		if payload.WrappedPayloadID.Valid {
			// get the associated UUID for the wrapped payload
			wrappedPayload := databaseStructs.Payload{}
			if err := database.DB.Get(&wrappedPayload, `SELECT uuid FROM payload WHERE id=$1`, payload.WrappedPayloadID.Int64); err != nil {
				logging.LogError(err, "Failed to fetch wrapped payload information")
			} else {
				payloadConfiguration.WrappedPayloadUUID = wrappedPayload.UuID
			}
		}
		return payloadConfiguration, nil
	}
}