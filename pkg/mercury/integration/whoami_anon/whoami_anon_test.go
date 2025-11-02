package whoami_anon_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/sprucelabsai/mercury-client-go/pkg/mercury/internal/helpers"
)

func TestWhoAmIAnonymous(t *testing.T) {
	helpers.LoadTestEnv(t)
	helpers.SetupSocketConnect(t)

	client := helpers.MakeClientWithTestHost(t)
	defer client.Disconnect()

	_, authType := helpers.EmitWhoAmI(t, client)
	require.Equal(t, "anonymous", authType, "Auth should be anonymous")
}
