package skill_to_skill_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	mercury "github.com/sprucelabsai-community/mercury-client-go/pkg/mercury"
	"github.com/sprucelabsai-community/mercury-client-go/pkg/testkit"
)

func TestSkillToSkillEmission(t *testing.T) {
	testkit.BeforeEach(t)

	org, skill1Client, skill2Client, fqen := testkit.LoginCreateOrgSetupTwoSkills(t)
	defer skill1Client.Disconnect()
	defer skill2Client.Disconnect()

	wasHit := false
	messages := []string{testkit.GenerateRandomID(), testkit.GenerateRandomID(), testkit.GenerateRandomID()}

	skill2Client.On(fqen, func(targetAndPayload mercury.TargetAndPayload) any {
		wasHit = true
		return map[string]any{
			"messages": messages,
		}
	})

	results := testkit.EmitSkillEvent(t, skill1Client, fqen, org.Id, testkit.GenerateRandomID())

	require.True(t, wasHit, "Event handler should have been hit")
	require.Equal(t, 1, len(results), "There should be one result")

	first := results[0]
	returnedMessages, ok := first["messages"].([]any)

	require.True(t, ok, "Messages field should be present in response")

	require.Equal(t, len(messages), len(returnedMessages), "Returned messages length should match sent messages length")
	require.Equal(t, messages[0], returnedMessages[0], "Returned message should match sent message")
}
