package whoami_login_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/sprucelabsai/mercury-client-go/pkg/mercury/internal/helpers"
)

func TestWhoAmILoggedIn(t *testing.T) {
	helpers.LoadTestEnv(t)
	helpers.SetupSocketConnect(t)

	client, person, _ := helpers.LoginAsDemoPerson(t, "+1 555-555-5555")
	defer client.Disconnect()

	who, _ := helpers.EmitWhoAmI(t, client)
	require.Equal(t, person.Id, who.Id, "Person id should match")
}
