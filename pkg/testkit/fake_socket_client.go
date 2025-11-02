package testkit

import (
	"sync"

	"github.com/sprucelabsai/mercury-client-go/pkg/mercury"
	ioClient "github.com/zishang520/socket.io/clients/socket/v3"
	socketTypes "github.com/zishang520/socket.io/v3/pkg/types"
)

type FakeSocketClient struct {
	host         string
	opts         ioClient.OptionsInterface
	is_connected bool
}

var (
	lastFakeSocketMu sync.RWMutex
	lastFakeSocket   *FakeSocketClient
)

func FakeSocketConnect(host string, opts ioClient.OptionsInterface) (mercury.Socket, error) {
	client := &FakeSocketClient{
		host: host,
		opts: opts,
	}

	client.is_connected = true
	setLastFakeSocket(client)

	return client, nil
}

func (s *FakeSocketClient) Emit(event string, args ...any) error {
	return nil
}

func (s *FakeSocketClient) On(event string, listeners ...socketTypes.EventListener) error {
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
