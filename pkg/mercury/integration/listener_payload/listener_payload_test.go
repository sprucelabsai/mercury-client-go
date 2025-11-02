package listener_payload_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	mercury "github.com/sprucelabsai/mercury-client-go/pkg/mercury"
	"github.com/sprucelabsai/mercury-client-go/pkg/mercury/internal/helpers"
)

func TestListenerReceivesTargetAndPayload(t *testing.T) {
	helpers.LoadTestEnv(t)
	helpers.SetupSocketConnect(t)

	org, skill1Client, skill2Client, fqen := helpers.LoginCreateOrgSetupTwoSkills(t)
	defer skill1Client.Disconnect()
	defer skill2Client.Disconnect()

	var captured mercury.TargetAndPayload
	skill2Client.On(fqen, func(targetAndPayload mercury.TargetAndPayload) any {
		captured = targetAndPayload
		return map[string]any{
			"messages": []string{helpers.GenerateRandomID()},
		}
	})

	actual := mercury.TargetAndPayload{
		Target: map[string]any{
			"organizationId": org.Id,
		},
		Payload: map[string]any{
			"message": helpers.GenerateRandomID(),
		},
	}

	_, err := skill1Client.Emit(fqen, actual)
	require.NoError(t, err, "Emitting custom event should not return an error")

	require.Equal(t, actual.Target, captured.Target, "Targets should match")
	require.Equal(t, actual.Payload, captured.Payload, "Payloads should match")
}
