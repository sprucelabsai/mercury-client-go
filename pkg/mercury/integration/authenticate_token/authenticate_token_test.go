package authenticate_token_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	mercury "github.com/sprucelabsai/mercury-client-go/pkg/mercury"
	"github.com/sprucelabsai/mercury-client-go/pkg/mercury/internal/helpers"
)

func TestAuthenticateWithExistingToken(t *testing.T) {
	helpers.LoadTestEnv(t)
	helpers.SetupSocketConnect(t)

	initialClient, person, token := helpers.LoginAsDemoPerson(t, "+1 555-555-5555")
	initialClient.Disconnect()

	client := helpers.MakeClientWithTestHost(t)
	defer client.Disconnect()

	_, err := client.Authenticate(mercury.AuthenticatePayload{Token: token})
	require.NoError(t, err, "Authenticating with token should not error")

	who, _ := helpers.EmitWhoAmI(t, client)
	require.Equal(t, person.Id, who.Id, "Person id should match")
}
