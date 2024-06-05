package rabbitmq

import (
	"github.com/its-a-feature/tiger/logging"
)

func (r *rabbitMQConnection) SendPtTaskCompletionFunction(taskMessage PTTaskCompletionFunctionMessage) error {
	if err := r.SendStructMessage(
		tiger_EXCHANGE,
		GetPtTaskCompletionHandlerRoutingKey(taskMessage.TaskData.PayloadType),
		"",
		taskMessage,
		false,
	); err != nil {
		logging.LogError(err, "Failed to send message")
		return err
	} else {
		return nil
	}
}
