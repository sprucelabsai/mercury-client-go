package testkit

import (
	"fmt"

	"github.com/sprucelabsai-community/mercury-client-go/pkg/mercury"
)

// MakeFakeClient creates a Mercury client wired to the in-memory FakeSocketClient.
func MakeFakeClient(opts ...mercury.MercuryClientOptions) (*FakeSocketClient, mercury.MercuryClient, error) {
	mercury.SetConnect(FakeSocketConnect)

	client, err := mercury.MakeMercuryClient(opts...)
	if err != nil {
		return nil, nil, err
	}

	fake := LastFakeSocket()
	if fake == nil {
		return nil, nil, fmt.Errorf("fake socket not captured")
	}

	return fake, client, nil
}

// ResetConnect restores the client connect function to the default implementation.
func ResetConnect() {
	mercury.SetConnect(nil)
	setLastFakeSocket(nil)
}
