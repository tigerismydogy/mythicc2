package rabbitmq

import (
	"github.com/its-a-feature/tiger/logging"
)

func (r *rabbitMQConnection) SendPtTaskProcessResponse(taskMessage PtTaskProcessResponseMessage) error {
	if err := r.SendStructMessage(
		tiger_EXCHANGE,
		GetPtTaskProcessResponseRoutingKey(taskMessage.TaskData.PayloadType),
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
