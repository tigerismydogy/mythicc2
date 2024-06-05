package rabbitmq

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/its-a-feature/tiger/grpc"
	"github.com/its-a-feature/tiger/grpc/services"
	"time"

	tigerCrypto "github.com/its-a-feature/tiger/crypto"
	"github.com/its-a-feature/tiger/logging"
)

// TRANSLATION_CONTAINER_CUSTOM_MESSAGE_TO_tiger_C2_FORMAT STRUCTS

type TrCustomMessageTotigerC2FormatMessage struct {
	TranslationContainerName string                    `json:"translation_container_name"`
	C2Name                   string                    `json:"c2_profile_name"`
	Message                  []byte                    `json:"message"`
	UUID                     string                    `json:"uuid"`
	tigerEncrypts           bool                      `json:"tiger_encrypts"`
	CryptoKeys               []tigerCrypto.CryptoKeys `json:"crypto_keys"`
}

type TrCustomMessageTotigerC2FormatMessageResponse struct {
	Success bool                   `json:"success"`
	Error   string                 `json:"error"`
	Message map[string]interface{} `json:"message"`
}

func (r *rabbitMQConnection) SendTrRPCCustomMessageTotigerC2(totigerC2Format TrCustomMessageTotigerC2FormatMessage) (*TrCustomMessageTotigerC2FormatMessageResponse, error) {
	trCustomMessageTotigerC2Format := TrCustomMessageTotigerC2FormatMessageResponse{}
	grpcSendMsg := services.TrCustomMessageTotigerC2FormatMessage{
		TranslationContainerName: totigerC2Format.TranslationContainerName,
		C2Name:                   totigerC2Format.C2Name,
		Message:                  totigerC2Format.Message,
		UUID:                     totigerC2Format.UUID,
		tigerEncrypts:           totigerC2Format.tigerEncrypts,
	}
	adjustedKeys := make([]*services.CryptoKeysFormat, len(totigerC2Format.CryptoKeys))
	for i := 0; i < len(totigerC2Format.CryptoKeys); i++ {
		newCrypto := services.CryptoKeysFormat{}
		adjustedKeys[i] = &newCrypto
		adjustedKeys[i].Value = totigerC2Format.CryptoKeys[i].Value
		if totigerC2Format.CryptoKeys[i].EncKey != nil {
			adjustedKeys[i].EncKey = *totigerC2Format.CryptoKeys[i].EncKey
		}
		if totigerC2Format.CryptoKeys[i].DecKey != nil {
			adjustedKeys[i].DecKey = *totigerC2Format.CryptoKeys[i].DecKey
		}
	}
	grpcSendMsg.CryptoKeys = adjustedKeys
	if sndMsgChan, rcvMsgChan, err := grpc.TranslationContainerServer.GetCustomTotigerChannels(totigerC2Format.TranslationContainerName); err != nil {
		logging.LogError(err, "Failed to get channels for grpc to CustomC2 to tigerC2")
		trCustomMessageTotigerC2Format.Success = false
		trCustomMessageTotigerC2Format.Error = err.Error()
		return &trCustomMessageTotigerC2Format, err
	} else {
		select {
		case sndMsgChan <- grpcSendMsg:
		case <-time.After(grpc.TranslationContainerServer.GetTimeout()):
			return nil, errors.New(fmt.Sprintf("timeout trying to send to translation container: %s", totigerC2Format.TranslationContainerName))
		}

		select {
		case response, ok := <-rcvMsgChan:
			if !ok {
				logging.LogError(nil, "Failed to receive from translation container")
				return nil, errors.New(fmt.Sprintf("failed to receive from translation container: %s", totigerC2Format.TranslationContainerName))
			} else {
				if response.GetSuccess() {
					responseMap := map[string]interface{}{}
					if err := json.Unmarshal(response.Message, &responseMap); err != nil {
						logging.LogError(err, "Failed to convert tiger message to json bytes ")
						trCustomMessageTotigerC2Format.Success = false
						trCustomMessageTotigerC2Format.Error = err.Error()
					} else {
						trCustomMessageTotigerC2Format.Message = responseMap
						trCustomMessageTotigerC2Format.Success = response.GetSuccess()
						trCustomMessageTotigerC2Format.Error = response.GetError()
					}
				} else {
					trCustomMessageTotigerC2Format.Success = response.GetSuccess()
					trCustomMessageTotigerC2Format.Error = response.GetError()
				}
				return &trCustomMessageTotigerC2Format, err
			}
		case <-time.After(grpc.TranslationContainerServer.GetTimeout()):
			logging.LogError(err, "timeout hit waiting to receive a message from the translation container")
			return nil, errors.New(fmt.Sprintf("timeout hit waiting to receive message from the translation container: %s", totigerC2Format.TranslationContainerName))
		}
	}
	/*
		exclusiveQueue := true
		if opsecBytes, err := json.Marshal(totigerC2Format); err != nil {
			logging.LogError(err, "Failed to convert totigerC2Format to JSON", "totigerC2Format", totigerC2Format)
			return nil, err
		} else if response, err := r.SendRPCMessage(
			tiger_EXCHANGE,
			GetTrRPCConvertTotigerFormatRoutingKey(totigerC2Format.TranslationContainerName),
			opsecBytes,
			!exclusiveQueue,
		); err != nil {
			logging.LogError(err, "Failed to send RPC message")
			return nil, err
		} else if err := json.Unmarshal(response, &trCustomMessageTotigerC2Format); err != nil {
			logging.LogError(err, "Failed to parse tr custom message to tiger c2 response back to struct", "response", response)
			return nil, err
		} else {
			return &trCustomMessageTotigerC2Format, nil
		}

	*/
}
