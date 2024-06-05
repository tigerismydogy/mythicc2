package rabbitmq

import (
	"encoding/json"
	"os"

	"github.com/its-a-feature/tiger/database"
	databaseStructs "github.com/its-a-feature/tiger/database/structs"
	"github.com/its-a-feature/tiger/logging"
	amqp "github.com/rabbitmq/amqp091-go"
)

type tigerRPCFileGetContentMessage struct {
	AgentFileID string `json:"file_id"`
}

// Every tigerRPC function call must return a response that includes the following two values
type tigerRPCFileGetContentMessageResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
	Content []byte `json:"content"`
}

func init() {
	RabbitMQConnection.AddRPCQueue(RPCQueueStruct{
		Exchange:   tiger_EXCHANGE,
		Queue:      tiger_RPC_FILE_GET_CONTENT,    // swap out with queue in rabbitmq.constants.go file
		RoutingKey: tiger_RPC_FILE_GET_CONTENT,    // swap out with routing key in rabbitmq.constants.go file
		Handler:    processtigerRPCFileGetContent, // points to function that takes in amqp.Delivery and returns interface{}
	})
}

//tiger_RPC_OBJECT_ACTION - Say what the function does
func tigerRPCFileGetContent(input tigerRPCFileGetContentMessage) tigerRPCFileGetContentMessageResponse {
	response := tigerRPCFileGetContentMessageResponse{
		Success: false,
	}
	file := databaseStructs.Filemeta{
		AgentFileID: input.AgentFileID,
	}
	if err := database.DB.Get(&file, `SELECT id, path FROM filemeta WHERE agent_file_id=$1`, input.AgentFileID); err != nil {
		response.Error = err.Error()
		return response
	}

	if diskFile, err := os.OpenFile(file.Path, os.O_RDONLY, 0644); err != nil {
		response.Error = err.Error()
		return response
	} else if fileInfo, err := diskFile.Stat(); err != nil {
		response.Error = err.Error()
		return response
	} else {
		fileContents := make([]byte, fileInfo.Size())
		if _, err := diskFile.Read(fileContents); err != nil {
			response.Error = err.Error()
			return response
		} else {
			response.Content = fileContents
			response.Success = true
			return response
		}
	}
}
func processtigerRPCFileGetContent(msg amqp.Delivery) interface{} {
	incomingMessage := tigerRPCFileGetContentMessage{}
	responseMsg := tigerRPCFileGetContentMessageResponse{
		Success: false,
	}
	if err := json.Unmarshal(msg.Body, &incomingMessage); err != nil {
		logging.LogError(err, "Failed to unmarshal JSON into struct")
		responseMsg.Error = err.Error()
	} else {
		return tigerRPCFileGetContent(incomingMessage)
	}
	return responseMsg
}
