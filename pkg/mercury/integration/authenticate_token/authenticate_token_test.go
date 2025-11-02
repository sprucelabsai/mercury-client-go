package authenticate_token_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	mercury "github.com/sprucelabsai-community/mercury-client-go/pkg/mercury"
	"github.com/sprucelabsai-community/mercury-client-go/pkg/testkit"
)

func TestAuthenticateWithExistingToken(t *testing.T) {
	testkit.BeforeEach(t)

	initialClient, person, token := testkit.LoginAsDemoPerson(t, "+1 555-555-5555")
	initialClient.Disconnect()

	client := testkit.MakeClientWithTestHost(t)
	defer client.Disconnect()

	_, err := client.Authenticate(mercury.AuthenticatePayload{Token: token})
	require.NoError(t, err, "Authenticating with token should not error")

	who, _ := testkit.EmitWhoAmI(t, client)
	require.Equal(t, person.Id, who.Id, "Person id should match")
}
