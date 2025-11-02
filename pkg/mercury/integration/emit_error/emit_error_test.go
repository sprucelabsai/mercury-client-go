package emit_error_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/sprucelabsai/mercury-client-go/pkg/mercury/internal/helpers"
)

func TestEmitNonexistentEventReturnsError(t *testing.T) {
	helpers.LoadTestEnv(t)
	helpers.SetupSocketConnect(t)

	client := helpers.MakeClientWithTestHost(t)
	defer client.Disconnect()

	_, err := client.Emit("this-event-does-not-exist::v2020_12_25")
	require.Error(t, err, "Emitting non-existent event should return an error")
}
