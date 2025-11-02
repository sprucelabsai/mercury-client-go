package testkit

import (
	"fmt"
	"sync"

	"github.com/sprucelabsai-community/mercury-client-go/pkg/mercury"
	ioClient "github.com/zishang520/socket.io/clients/socket/v3"
	socketTypes "github.com/zishang520/socket.io/v3/pkg/types"
)

type FakeSocketClient struct {
	host         string
	opts         ioClient.OptionsInterface
	is_connected bool
	listeners    []FakedListener
}

type FakedListener struct {
	fqen string
	cb   socketTypes.EventListener
}

type SocketIoEmitCallback = func([]any, error)

var (
	lastFakeSocketMu sync.RWMutex
	lastFakeSocket   *FakeSocketClient
)

func FakeSocketConnect(host string, opts ioClient.OptionsInterface) (mercury.Socket, error) {
	if lastFakeSocket != nil {
		return lastFakeSocket, nil
	}

	client := &FakeSocketClient{
		host: host,
		opts: opts,
	}

	client.is_connected = true
	setLastFakeSocket(client)

	client.On("register-listeners::v2020_12_25", func(args ...any) {
		cb := PluckCallback(args)
		if cb != nil {
			cb([]any{}, nil)
		}
	})

	return client, nil
}

func (s *FakeSocketClient) Emit(event string, args ...any) error {
	cb := PluckCallback(args)

	for _, listener := range s.listeners {
		if listener.fqen == event {
			argsWithBridge := args[:len(args)-1]
			bridge := func(responseArgs []any, _ error) {
				mapped := wrapInAggregateResponse(mercury.ResponsePayload{})
				fmt.Print("tmest", mapped)
				if cb != nil {
					cb([]any{mapped}, nil)
				}
			}

			argsWithBridge = append(argsWithBridge, bridge)
			listener.cb(argsWithBridge...)

			return nil
		}
	}

	if cb != nil {
		cb(nil, fmt.Errorf("no listener registered for event %s", event))
	}

	return nil
}

func wrapInAggregateResponse(payload mercury.ResponsePayload) mercury.MercuryAggregateResponse {
	return mercury.MercuryAggregateResponse{
		TotalContracts: 1,
		TotalResponses: 1,
		TotalErrors:    0,
		Responses: []mercury.MercurySingleResponse{
			{
				ResponderRef: "fake-responder-1",
				Errors:       []any{},
				Payload:      payload,
			},
		},
	}
}

func (s *FakeSocketClient) On(event string, listeners ...socketTypes.EventListener) error {
	if len(listeners) > 0 {

		var filteredListeners []FakedListener
		for _, existing := range s.listeners {
			if existing.fqen != event {
				filteredListeners = append(filteredListeners, existing)
			}
		}
		s.listeners = filteredListeners

		s.listeners = append(s.listeners, FakedListener{
			fqen: event,
			cb:   listeners[0],
		})
	}
	return nil
}

func (s *FakeSocketClient) Connected() bool {
	return s.is_connected
}

func (s *FakeSocketClient) Disconnect() mercury.Socket {
	return s
}

func (s *FakeSocketClient) SetConnected(connected bool) {
	s.is_connected = connected
}

func (s *FakeSocketClient) GetOptions() ioClient.OptionsInterface {
	return s.opts
}

func (s *FakeSocketClient) GetHost() string {
	return s.host
}

func (s *FakeSocketClient) Off(event string, listener socketTypes.EventListener) bool {
	return false
}

func LastFakeSocket() *FakeSocketClient {
	lastFakeSocketMu.RLock()
	defer lastFakeSocketMu.RUnlock()
	return lastFakeSocket
}

func setLastFakeSocket(fake *FakeSocketClient) {
	lastFakeSocketMu.Lock()
	defer lastFakeSocketMu.Unlock()
	lastFakeSocket = fake
}
