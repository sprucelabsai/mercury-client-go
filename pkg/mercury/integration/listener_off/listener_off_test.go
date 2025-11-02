package listener_off_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	mercury "github.com/sprucelabsai-community/mercury-client-go/pkg/mercury"
	"github.com/sprucelabsai-community/mercury-client-go/pkg/testkit"
)

func TestListenerOffStopsEvents(t *testing.T) {
	testkit.BeforeEach(t)

	org, skill1Client, skill2Client, fqen := testkit.LoginCreateOrgSetupTwoSkills(t)
	defer skill1Client.Disconnect()
	defer skill2Client.Disconnect()

	hitCount := 0
	skill2Client.On(fqen, func(targetAndPayload mercury.TargetAndPayload) any {
		hitCount++
		return map[string]any{
			"messages": []string{testkit.GenerateRandomID()},
		}
	})

	testkit.EmitSkillEvent(t, skill1Client, fqen, org.Id, testkit.GenerateRandomID())
	skill2Client.Off(fqen)

	testkit.EmitSkillEvent(t, skill1Client, fqen, org.Id, testkit.GenerateRandomID())
	require.Equal(t, 1, hitCount, "Hit count should be 1 after turning off listener")
}
