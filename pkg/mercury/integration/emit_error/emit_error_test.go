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

	_, err := client.Emit("this-event-does-not-exist::v2020_12_25")
	require.Error(t, err, "Emitting non-existent event should return an error")
}
