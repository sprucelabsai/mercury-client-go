package connect_test

import (
	"testing"

	"github.com/sprucelabsai-community/mercury-client-go/pkg/testkit"
	"github.com/stretchr/testify/require"
)

func TestConnectAndDisconnect(t *testing.T) {
	testkit.BeforeEachInternal(t)

	client := testkit.MakeClientWithTestHost(t)
	require.True(t, client.IsConnected(), "Client should be connected on construction")

	client.Disconnect()
	require.False(t, client.IsConnected(), "Client should be disconnected after disconnect")
}
