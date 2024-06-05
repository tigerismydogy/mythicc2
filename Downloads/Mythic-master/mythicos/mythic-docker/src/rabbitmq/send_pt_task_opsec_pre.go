package rabbitmq

import (
	"github.com/its-a-feature/tiger/logging"
)

func (r *rabbitMQConnection) SendPtTaskOPSECPre(taskMessage PTTaskMessageAllData) error {
	if err := r.SendStructMessage(
		tiger_EXCHANGE,
		GetPtTaskOpsecPreCheckRoutingKey(taskMessage.PayloadType),
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
