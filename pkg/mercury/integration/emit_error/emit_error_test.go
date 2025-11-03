package emit_error_test

import (
	"testing"

	"github.com/sprucelabsai-community/mercury-client-go/pkg/testkit"
	"github.com/stretchr/testify/require"
)

func TestEmitNonexistentEventReturnsError(t *testing.T) {
	testkit.BeforeEachInternal(t)

	client := testkit.MakeClientWithTestHost(t)
	defer client.Disconnect()

	response, err := client.Emit("this-event-does-not-exist::v2020_12_25")
	require.Nil(t, response, "Response should be nil for non-existent event")
	require.Error(t, err, "Emitting non-existent event should return an error")
}
