package grpc

import (
	"errors"
	"fmt"
	"github.com/its-a-feature/tiger/grpc/services"
	"github.com/its-a-feature/tiger/logging"
	"github.com/its-a-feature/tiger/utils"
	"google.golang.org/grpc"
	"math"
	"net"
	"sync"
	"time"
)

const (
	connectionTimeoutSeconds  = 3
	channelSendTimeoutSeconds = 1
)

type translationContainerServer struct {
	services.UnimplementedTranslationContainerServer
	sync.RWMutex
	clients            map[string]*grpcTranslationContainerClientConnections
	connectionTimeout  time.Duration
	channelSendTimeout time.Duration
	listening          bool
	latestError        string
}
type PushC2ServerConnected struct {
	PushC2MessagesTotiger   chan RabbitMQProcessAgentMessageFromPushC2
	DisconnectProcessingChan chan bool
}
type pushC2Server struct {
	services.UnimplementedPushC2Server
	sync.RWMutex
	clients                              map[int]*grpcPushC2ClientConnections
	clientsOneToMany                     map[string]*grpcPushC2ClientConnections
	connectionTimeout                    time.Duration
	channelSendTimeout                   time.Duration
	rabbitmqProcessPushC2AgentConnection chan PushC2ServerConnected
	listening                            bool
	latestError                          string
}

type grpcTranslationContainerClientConnections struct {
	generateKeysMessage                          chan services.TrGenerateEncryptionKeysMessage
	generateKeysMessageResponse                  chan services.TrGenerateEncryptionKeysMessageResponse
	translateCustomTotigerFormatMessage         chan services.TrCustomMessageTotigerC2FormatMessage
	translateCustomTotigerFormatMessageResponse chan services.TrCustomMessageTotigerC2FormatMessageResponse
	translatetigerToCustomFormatMessage         chan services.TrtigerC2ToCustomMessageFormatMessage
	translatetigerToCustomFormatMessageResponse chan services.TrtigerC2ToCustomMessageFormatMessageResponse
	connectedGenerateKeys                        bool
	connectedCustomTotiger                      bool
	connectedtigerToCustom                      bool
	sync.RWMutex
}
type grpcPushC2ClientConnections struct {
	pushC2MessageFromtiger chan services.PushC2MessageFromtiger
	connected               bool
	callbackUUID            string
	base64Encoded           bool
	c2ProfileName           string
	AgentUUIDSize           int
	sync.RWMutex
}

var TranslationContainerServer translationContainerServer
var PushC2Server pushC2Server

// translationContainerServer functions
func (t *translationContainerServer) addNewGenerateKeysClient(translationContainerName string) (chan services.TrGenerateEncryptionKeysMessage, chan services.TrGenerateEncryptionKeysMessageResponse, error) {
	t.Lock()
	if _, ok := t.clients[translationContainerName]; !ok {
		t.clients[translationContainerName] = &grpcTranslationContainerClientConnections{}
		t.clients[translationContainerName].generateKeysMessage = make(chan services.TrGenerateEncryptionKeysMessage)
		t.clients[translationContainerName].generateKeysMessageResponse = make(chan services.TrGenerateEncryptionKeysMessageResponse)
		t.clients[translationContainerName].translateCustomTotigerFormatMessage = make(chan services.TrCustomMessageTotigerC2FormatMessage)
		t.clients[translationContainerName].translateCustomTotigerFormatMessageResponse = make(chan services.TrCustomMessageTotigerC2FormatMessageResponse)
		t.clients[translationContainerName].translatetigerToCustomFormatMessage = make(chan services.TrtigerC2ToCustomMessageFormatMessage)
		t.clients[translationContainerName].translatetigerToCustomFormatMessageResponse = make(chan services.TrtigerC2ToCustomMessageFormatMessageResponse)
	}
	msg := t.clients[translationContainerName].generateKeysMessage
	rsp := t.clients[translationContainerName].generateKeysMessageResponse
	t.clients[translationContainerName].connectedGenerateKeys = true
	t.Unlock()
	return msg, rsp, nil
}
func (t *translationContainerServer) addNewCustomTotigerClient(translationContainerName string) (chan services.TrCustomMessageTotigerC2FormatMessage, chan services.TrCustomMessageTotigerC2FormatMessageResponse, error) {
	t.Lock()
	if _, ok := t.clients[translationContainerName]; !ok {
		t.clients[translationContainerName] = &grpcTranslationContainerClientConnections{}
		t.clients[translationContainerName].generateKeysMessage = make(chan services.TrGenerateEncryptionKeysMessage)
		t.clients[translationContainerName].generateKeysMessageResponse = make(chan services.TrGenerateEncryptionKeysMessageResponse)
		t.clients[translationContainerName].translateCustomTotigerFormatMessage = make(chan services.TrCustomMessageTotigerC2FormatMessage)
		t.clients[translationContainerName].translateCustomTotigerFormatMessageResponse = make(chan services.TrCustomMessageTotigerC2FormatMessageResponse)
		t.clients[translationContainerName].translatetigerToCustomFormatMessage = make(chan services.TrtigerC2ToCustomMessageFormatMessage)
		t.clients[translationContainerName].translatetigerToCustomFormatMessageResponse = make(chan services.TrtigerC2ToCustomMessageFormatMessageResponse)
	}
	msg := t.clients[translationContainerName].translateCustomTotigerFormatMessage
	rsp := t.clients[translationContainerName].translateCustomTotigerFormatMessageResponse
	t.clients[translationContainerName].connectedCustomTotiger = true
	t.Unlock()
	return msg, rsp, nil
}
func (t *translationContainerServer) addNewtigerToCustomClient(translationContainerName string) (chan services.TrtigerC2ToCustomMessageFormatMessage, chan services.TrtigerC2ToCustomMessageFormatMessageResponse, error) {
	t.Lock()
	if _, ok := t.clients[translationContainerName]; !ok {
		t.clients[translationContainerName] = &grpcTranslationContainerClientConnections{}
		t.clients[translationContainerName].generateKeysMessage = make(chan services.TrGenerateEncryptionKeysMessage)
		t.clients[translationContainerName].generateKeysMessageResponse = make(chan services.TrGenerateEncryptionKeysMessageResponse)
		t.clients[translationContainerName].translateCustomTotigerFormatMessage = make(chan services.TrCustomMessageTotigerC2FormatMessage)
		t.clients[translationContainerName].translateCustomTotigerFormatMessageResponse = make(chan services.TrCustomMessageTotigerC2FormatMessageResponse)
		t.clients[translationContainerName].translatetigerToCustomFormatMessage = make(chan services.TrtigerC2ToCustomMessageFormatMessage)
		t.clients[translationContainerName].translatetigerToCustomFormatMessageResponse = make(chan services.TrtigerC2ToCustomMessageFormatMessageResponse)
	}
	msg := t.clients[translationContainerName].translatetigerToCustomFormatMessage
	rsp := t.clients[translationContainerName].translatetigerToCustomFormatMessageResponse
	t.clients[translationContainerName].connectedtigerToCustom = true
	t.Unlock()
	return msg, rsp, nil
}
func (t *translationContainerServer) GetGenerateKeysChannels(translationContainerName string) (chan services.TrGenerateEncryptionKeysMessage, chan services.TrGenerateEncryptionKeysMessageResponse, error) {
	t.RLock()
	defer t.RUnlock()
	if _, ok := t.clients[translationContainerName]; ok {
		return t.clients[translationContainerName].generateKeysMessage,
			t.clients[translationContainerName].generateKeysMessageResponse,
			nil
	}
	return nil, nil, errors.New("no translation container by that name currently connected")
}
func (t *translationContainerServer) SetGenerateKeysChannelExited(translationContainerName string) {
	t.RLock()
	defer t.RUnlock()
	if _, ok := t.clients[translationContainerName]; ok {
		t.clients[translationContainerName].connectedGenerateKeys = false
	}
}
func (t *translationContainerServer) GetCustomTotigerChannels(translationContainerName string) (chan services.TrCustomMessageTotigerC2FormatMessage, chan services.TrCustomMessageTotigerC2FormatMessageResponse, error) {
	t.RLock()
	defer t.RUnlock()
	if _, ok := t.clients[translationContainerName]; ok {
		return t.clients[translationContainerName].translateCustomTotigerFormatMessage,
			t.clients[translationContainerName].translateCustomTotigerFormatMessageResponse,
			nil
	}
	return nil, nil, errors.New("no translation container by that name currently connected")
}
func (t *translationContainerServer) SetCustomTotigerChannelExited(translationContainerName string) {
	t.RLock()
	defer t.RUnlock()
	if _, ok := t.clients[translationContainerName]; ok {
		t.clients[translationContainerName].connectedCustomTotiger = false
	}
}
func (t *translationContainerServer) GettigerToCustomChannels(translationContainerName string) (chan services.TrtigerC2ToCustomMessageFormatMessage, chan services.TrtigerC2ToCustomMessageFormatMessageResponse, error) {
	t.RLock()
	defer t.RUnlock()
	if _, ok := t.clients[translationContainerName]; ok {
		return t.clients[translationContainerName].translatetigerToCustomFormatMessage,
			t.clients[translationContainerName].translatetigerToCustomFormatMessageResponse,
			nil
	}
	return nil, nil, errors.New("no translation container by that name currently connected")
}
func (t *translationContainerServer) SettigerToCustomChannelExited(translationContainerName string) {
	t.RLock()
	if _, ok := t.clients[translationContainerName]; ok {
		t.clients[translationContainerName].connectedtigerToCustom = false
	}
	t.RUnlock()
}
func (t *translationContainerServer) CheckClientConnected(translationContainerName string) bool {
	t.RLock()
	if _, ok := t.clients[translationContainerName]; ok {
		t.RUnlock()
		return t.clients[translationContainerName].connectedtigerToCustom &&
			t.clients[translationContainerName].connectedGenerateKeys &&
			t.clients[translationContainerName].connectedCustomTotiger
	} else {
		t.RUnlock()
		return false
	}
}
func (t *translationContainerServer) CheckListening() (listening bool, latestError string) {
	return t.listening, t.latestError
}
func (t *translationContainerServer) GetTimeout() time.Duration {
	return t.connectionTimeout
}
func (t *translationContainerServer) GetChannelTimeout() time.Duration {
	return t.channelSendTimeout
}

// pushC2Server functions
func (t *pushC2Server) GetRabbitMqProcessAgentMessageChannel() chan PushC2ServerConnected {
	return t.rabbitmqProcessPushC2AgentConnection
}
func (t *pushC2Server) addNewPushC2Client(CallbackAgentID int, callbackUUID string, base64Encoded bool, c2ProfileName string, agentUUIDSize int) (chan services.PushC2MessageFromtiger, error) {
	t.Lock()
	if _, ok := t.clients[CallbackAgentID]; !ok {
		t.clients[CallbackAgentID] = &grpcPushC2ClientConnections{}
		t.clients[CallbackAgentID].pushC2MessageFromtiger = make(chan services.PushC2MessageFromtiger, 100)
	}
	fromtiger := t.clients[CallbackAgentID].pushC2MessageFromtiger
	t.clients[CallbackAgentID].connected = true
	t.clients[CallbackAgentID].callbackUUID = callbackUUID
	t.clients[CallbackAgentID].base64Encoded = base64Encoded
	t.clients[CallbackAgentID].c2ProfileName = c2ProfileName
	t.clients[CallbackAgentID].AgentUUIDSize = agentUUIDSize
	t.Unlock()
	return fromtiger, nil
}
func (t *pushC2Server) addNewPushC2OneToManyClient(c2ProfileName string) (chan services.PushC2MessageFromtiger, error) {
	t.Lock()
	if _, ok := t.clientsOneToMany[c2ProfileName]; !ok {
		t.clientsOneToMany[c2ProfileName] = &grpcPushC2ClientConnections{}
		t.clientsOneToMany[c2ProfileName].pushC2MessageFromtiger = make(chan services.PushC2MessageFromtiger, 100)
	}
	fromtiger := t.clientsOneToMany[c2ProfileName].pushC2MessageFromtiger
	t.clientsOneToMany[c2ProfileName].connected = true
	t.clientsOneToMany[c2ProfileName].c2ProfileName = c2ProfileName
	t.Unlock()
	return fromtiger, nil
}
func (t *pushC2Server) GetPushC2ClientInfo(CallbackAgentID int) (chan services.PushC2MessageFromtiger, string, bool, string, string, int, error) {
	t.RLock()
	if _, ok := t.clients[CallbackAgentID]; ok {
		if t.clients[CallbackAgentID].connected {
			t.RUnlock()
			return t.clients[CallbackAgentID].pushC2MessageFromtiger,
				t.clients[CallbackAgentID].callbackUUID,
				t.clients[CallbackAgentID].base64Encoded,
				t.clients[CallbackAgentID].c2ProfileName,
				t.clients[CallbackAgentID].callbackUUID,
				t.clients[CallbackAgentID].AgentUUIDSize,
				nil
		} else {
			t.RUnlock()
			return nil, "", false, "", "", 0, errors.New("push c2 channel for that callback is no longer available")
		}

	}
	for c2, _ := range t.clientsOneToMany {
		c2ProfileToCallbackIDsMapLock.RLock()
		if _, ok := c2ProfileToCallbackIDsMap[c2]; ok {
			if _, ok = c2ProfileToCallbackIDsMap[c2][CallbackAgentID]; ok {
				t.RUnlock()
				c2ProfileToCallbackIDsMapLock.RUnlock()
				return t.clientsOneToMany[c2].pushC2MessageFromtiger,
					c2ProfileToCallbackIDsMap[c2][CallbackAgentID].CallbackUUID,
					c2ProfileToCallbackIDsMap[c2][CallbackAgentID].Base64Encoded,
					c2,
					c2ProfileToCallbackIDsMap[c2][CallbackAgentID].TrackingID,
					c2ProfileToCallbackIDsMap[c2][CallbackAgentID].AgentUUIDSize,
					nil
			}
		}
		c2ProfileToCallbackIDsMapLock.RUnlock()
	}
	t.RUnlock()
	return nil, "", false, "", "", 0, errors.New("no push c2 channel for that callback available")
}
func (t *pushC2Server) SetPushC2ChannelExited(CallbackAgentID int) {
	t.RLock()
	if _, ok := t.clients[CallbackAgentID]; ok {
		t.clients[CallbackAgentID].connected = false
	}
	t.RUnlock()
}
func (t *pushC2Server) SetPushC2OneToManyChannelExited(c2ProfileName string) {
	t.RLock()
	if _, ok := t.clientsOneToMany[c2ProfileName]; ok {
		t.clientsOneToMany[c2ProfileName].connected = false
	}
	t.RUnlock()
}
func (t *pushC2Server) CheckListening() (listening bool, latestError string) {
	return t.listening, t.latestError
}
func (t *pushC2Server) CheckClientConnected(CallbackAgentID int) bool {
	t.RLock()
	defer t.RUnlock()
	if _, ok := t.clients[CallbackAgentID]; ok {
		return t.clients[CallbackAgentID].connected
	}
	for c2, _ := range t.clientsOneToMany {
		c2ProfileToCallbackIDsMapLock.RLock()
		if _, ok := c2ProfileToCallbackIDsMap[c2]; ok {
			if _, ok = c2ProfileToCallbackIDsMap[c2][CallbackAgentID]; ok {
				c2ProfileToCallbackIDsMapLock.RUnlock()
				return t.clientsOneToMany[c2].connected
			}
		}
		c2ProfileToCallbackIDsMapLock.RUnlock()
	}
	return false
}
func (t *pushC2Server) GetConnectedClients() []int {
	clientIDs := []int{}
	t.RLock()
	defer t.RUnlock()
	for clientID, _ := range t.clients {
		if t.clients[clientID].connected {
			clientIDs = append(clientIDs, clientID)
		}
	}
	for c2, _ := range t.clientsOneToMany {
		if t.clientsOneToMany[c2].connected {
			c2ProfileToCallbackIDsMapLock.RLock()
			for newId, _ := range c2ProfileToCallbackIDsMap[c2] {
				clientIDs = append(clientIDs, newId)
			}
			c2ProfileToCallbackIDsMapLock.RUnlock()
		}
	}
	return clientIDs
}
func (t *pushC2Server) GetTimeout() time.Duration {
	return t.connectionTimeout
}
func (t *pushC2Server) GetChannelTimeout() time.Duration {
	return t.channelSendTimeout
}

func Initialize() {
	// need to open a port to accept gRPC connections
	var (
		connectString string
	)
	// initialize the clients
	TranslationContainerServer.clients = make(map[string]*grpcTranslationContainerClientConnections)
	TranslationContainerServer.connectionTimeout = connectionTimeoutSeconds * time.Second
	TranslationContainerServer.channelSendTimeout = channelSendTimeoutSeconds * time.Second
	// initial for push c2 servers
	PushC2Server.clients = make(map[int]*grpcPushC2ClientConnections)
	PushC2Server.clientsOneToMany = make(map[string]*grpcPushC2ClientConnections)
	PushC2Server.rabbitmqProcessPushC2AgentConnection = make(chan PushC2ServerConnected, 20)
	PushC2Server.connectionTimeout = connectionTimeoutSeconds * time.Second
	PushC2Server.channelSendTimeout = channelSendTimeoutSeconds * time.Second
	connectString = fmt.Sprintf("0.0.0.0:%d", utils.tigerConfig.ServerGRPCPort)
	go serveGRPCInBackground(connectString)

}
func serveGRPCInBackground(connectString string) {
	grpcServer := grpc.NewServer(grpc.MaxSendMsgSize(math.MaxInt), grpc.MaxRecvMsgSize(math.MaxInt))
	services.RegisterTranslationContainerServer(grpcServer, &TranslationContainerServer)
	services.RegisterPushC2Server(grpcServer, &PushC2Server)
	logging.LogInfo("Initializing grpc connections...")
	for {
		TranslationContainerServer.listening = false
		PushC2Server.listening = false
		if listen, err := net.Listen("tcp", connectString); err != nil {
			logging.LogError(err, "Failed to open port for gRPC connections, retrying...")
			TranslationContainerServer.latestError = err.Error()
			PushC2Server.latestError = err.Error()
			time.Sleep(TranslationContainerServer.GetTimeout())
			continue
		} else {
			TranslationContainerServer.listening = true
			PushC2Server.listening = true
			TranslationContainerServer.latestError = ""
			PushC2Server.latestError = ""
			// create a new instance of a grpc server
			logging.LogInfo("gRPC Initialized", "connection", connectString)
			// tie the Servers to our new grpc server and our server struct
			// use the TCP port in listen to process requests for the grpc server translationContainerGRPCServer
			if err = grpcServer.Serve(listen); err != nil {
				logging.LogError(err, "Failed to listen for gRPC connections")
				TranslationContainerServer.latestError = err.Error()
			}
		}
	}
}
