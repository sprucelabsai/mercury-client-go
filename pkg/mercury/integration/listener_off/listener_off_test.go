package listener_off_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	mercury "github.com/sprucelabsai-community/mercury-client-go/pkg/mercury"
	"github.com/sprucelabsai-community/mercury-client-go/pkg/mercury/internal/helpers"
)

func TestListenerOffStopsEvents(t *testing.T) {
	helpers.LoadTestEnv(t)
	helpers.SetupSocketConnect(t)

	org, skill1Client, skill2Client, fqen := helpers.LoginCreateOrgSetupTwoSkills(t)
	defer skill1Client.Disconnect()
	defer skill2Client.Disconnect()

	hitCount := 0
	skill2Client.On(fqen, func(targetAndPayload mercury.TargetAndPayload) any {
		hitCount++
		return map[string]any{
			"messages": []string{helpers.GenerateRandomID()},
		}
	})

	helpers.EmitSkillEvent(t, skill1Client, fqen, org.Id, helpers.GenerateRandomID())

	skill2Client.Off(fqen)

	helpers.EmitSkillEvent(t, skill1Client, fqen, org.Id, helpers.GenerateRandomID())

	require.Equal(t, 1, hitCount, "Hit count should be 1 after turning off listener")
}
