package whoami_anon_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/sprucelabsai-community/mercury-client-go/pkg/testkit"
)

func TestWhoAmIAnonymous(t *testing.T) {
	testkit.BeforeEachInternal(t)

	client := testkit.MakeClientWithTestHost(t)
	defer client.Disconnect()

	_, authType := testkit.EmitWhoAmI(t, client)
	require.Equal(t, "anonymous", authType, "Auth should be anonymous")
}
