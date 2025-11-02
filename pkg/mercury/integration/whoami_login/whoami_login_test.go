package whoami_login_test

import (
	"testing"

	"github.com/sprucelabsai-community/mercury-client-go/pkg/testkit"
	"github.com/stretchr/testify/require"
)

func TestWhoAmILoggedIn(t *testing.T) {
	testkit.BeforeEachInternal(t)

	client, person, _ := testkit.LoginAsDemoPerson(t, "+1 555-555-5555")
	defer client.Disconnect()

	who, _ := testkit.EmitWhoAmI(t, client)
	require.Equal(t, person.Id, who.Id, "Person id should match")
}
