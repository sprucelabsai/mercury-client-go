package connect_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/sprucelabsai-community/mercury-client-go/pkg/mercury/internal/helpers"
)

func TestConnectAndDisconnect(t *testing.T) {
	helpers.LoadTestEnv(t)
	helpers.SetupSocketConnect(t)

	client := helpers.MakeClientWithTestHost(t)
	require.True(t, client.IsConnected(), "Client should be connected on construction")

	client.Disconnect()
	require.False(t, client.IsConnected(), "Client should be disconnected after disconnect")
}
