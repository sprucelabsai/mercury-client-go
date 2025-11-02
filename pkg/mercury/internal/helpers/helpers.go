package helpers

import (
	"fmt"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	spruce "github.com/sprucelabsai-community/spruce-core-schemas/v41/pkg/schemas"
	schemas "github.com/sprucelabsai-community/spruce-core-schemas/v41/pkg/schemas/spruce/v2020_07_22"
	mercury "github.com/sprucelabsai/mercury-client-go/pkg/mercury"
	"github.com/stretchr/testify/require"
	ioClient "github.com/zishang520/socket.io/clients/socket/v3"
)

// EventContract represents an event contract payload.
type EventContract map[string]any

// SetupSocketConnect configures the mercury client to use the socket.io default connect.
func SetupSocketConnect(t *testing.T) {
	t.Helper()
	mercury.SetConnect(func(url string, opts ioClient.OptionsInterface) (mercury.Socket, error) {
		if opts == nil {
			opts = ioClient.DefaultOptions()
		}
		opts.SetForceNew(true)
		opts.SetMultiplex(false)
		opts.SetReconnection(true)

		socket, err := ioClient.Connect(url, opts)
		if err != nil {
			return nil, err
		}
		return mercury.NewSocketIOClient(socket), nil
	})

	t.Cleanup(func() {
		mercury.SetConnect(nil)
	})
}

// LoadTestEnv hydrates environment variables required by integration tests.
func LoadTestEnv(t *testing.T) {
	t.Helper()
	_ = godotenv.Load("../../../../.env", "../.env", ".env")
}

// MakeClientWithTestHost constructs a mercury client using the TEST_HOST env var.
func MakeClientWithTestHost(t *testing.T, opts ...mercury.MercuryClientOptions) mercury.MercuryClient {
	t.Helper()
	host := os.Getenv("TEST_HOST")

	require.NotEmpty(t, host, "TEST_HOST environment variable must be set for tests")
	fmt.Println("Making client with test host " + host)

	client, err := mercury.MakeMercuryClient(append(opts, mercury.MercuryClientOptions{Host: host})...)
	require.NoError(t, err, "Making Mercury client with test host should not return an error")

	fmt.Println("Made client with test host: ", host)

	return client
}

// GenerateRandomID produces a UUID string for tests.
func GenerateRandomID() string {
	return uuid.NewString()
}

// EmitWhoAmI returns the person and auth type from the whoami event.
func EmitWhoAmI(t *testing.T, client mercury.MercuryClient) (*spruce.Person, string) {
	t.Helper()
	auth, err := client.Emit("whoami::v2020_12_25")
	require.NoError(t, err, "Emit whoami should not return an error")
	require.NotNil(t, auth, "Emit whoami should return a response")
	require.Equal(t, 1, len(auth), "Emit whoami should return one response")
	first := auth[0]

	authMap := first["auth"].(map[string]any)

	authType := first["type"].(string)
	if authType == "anonymous" {
		return nil, "anonymous"
	}

	authPerson := authMap["person"].(map[string]any)
	person, err := schemas.MakePerson(authPerson)
	require.NoError(t, err, "Making person from whoami response should not return an error")
	require.NotNil(t, person, "Person from whoami should not be nil")

	return person, authType
}

// Login performs the PIN flow for the provided phone number and returns the person and token.
func Login(client mercury.MercuryClient, phone string) (*spruce.Person, string) {
	requestPinResponse, _ := client.Emit("request-pin::v2020_12_25", mercury.TargetAndPayload{
		Payload: map[string]any{
			"phone": phone,
		},
	})
	first := requestPinResponse[0]
	challenge := first["challenge"].(string)

	confirmPinResponse, _ := client.Emit("confirm-pin::v2020_12_25", mercury.TargetAndPayload{
		Payload: map[string]any{
			"challenge": challenge,
			"pin":       "0000",
		},
	})

	first = confirmPinResponse[0]
	token := first["token"].(string)
	personValues := first["person"].(map[string]any)
	person, _ := schemas.MakePerson(personValues)

	return person, token
}

// LoginAsDemoPerson logs in a demo person and returns the client, person, and auth token.
func LoginAsDemoPerson(t *testing.T, phone string) (mercury.MercuryClient, *spruce.Person, string) {
	t.Helper()
	fmt.Println("Logging in as demo person with phone:", phone)
	client := MakeClientWithTestHost(t)
	person, token := Login(client, phone)
	return client, person, token
}

// LoginAsSkill authenticates a skill client via API key.
func LoginAsSkill(t *testing.T, skill *spruce.Skill) (mercury.MercuryClient, error) {
	t.Helper()
	client := MakeClientWithTestHost(t)
	_, err := client.Authenticate(mercury.AuthenticatePayload{
		SkillId: skill.Id,
		ApiKey:  skill.ApiKey,
	})

	return client, err
}

// InstallSkill installs a skill in the given organization.
func InstallSkill(client mercury.MercuryClient, orgID string, skillID string) error {
	_, err := client.Emit("install-skill::v2020_12_25", mercury.TargetAndPayload{
		Target: map[string]any{
			"organizationId": orgID,
		},
		Payload: map[string]any{
			"skillId": skillID,
		},
	})
	return err
}

// SeedRandomSkill registers a random skill for tests.
func SeedRandomSkill(client mercury.MercuryClient) (*spruce.Skill, error) {
	skillName := fmt.Sprintf("Test Skill %s", uuid.NewString())
	results, err := client.Emit("register-skill::v2020_12_25", mercury.TargetAndPayload{
		Payload: map[string]any{
			"name": skillName,
		},
	})

	if err != nil {
		return nil, err
	}

	first := results[0]
	skillValues, ok := first["skill"].(map[string]any)

	if !ok {
		return nil, fmt.Errorf("skill field not found in response")
	}

	skill, err := schemas.MakeSkill(skillValues)
	if err != nil {
		return nil, err
	}

	return skill, nil
}

// SeedRandomOrg creates a random organization and returns it.
func SeedRandomOrg(t *testing.T, client mercury.MercuryClient) *spruce.Organization {
	t.Helper()
	orgName := fmt.Sprintf("Test Org %s", GenerateRandomID())
	results, err := client.Emit("create-organization::v2020_12_25", mercury.TargetAndPayload{
		Payload: map[string]any{
			"name": orgName,
		},
	})

	require.NoError(t, err, "Seeding organization should not return an error")

	fmt.Println("Create organization results:", results)
	first := results[0]

	orgValues, ok := first["organization"].(map[string]any)
	require.True(t, ok, "Organization field should be present in response")
	org, err := schemas.MakeOrganization(orgValues)
	require.NoError(t, err, "Making organization from response should not return an error")

	return org
}

// EmitSkillEvent triggers an event from skill1 to skill2 and returns the responses.
func EmitSkillEvent(t *testing.T, skillClient mercury.MercuryClient, fqen string, orgID string, message string) []mercury.ResponsePayload {
	t.Helper()
	results, err := skillClient.Emit(fqen, mercury.TargetAndPayload{
		Target: map[string]any{
			"organizationId": orgID,
		},
		Payload: map[string]any{
			"message": message,
		},
	})

	require.NoError(t, err, "Emitting event should not return an error")
	require.NotNil(t, results, "Results should not be nil")

	return results
}

// RegisterEvents registers the given contract and returns the FQEN.
func RegisterEvents(t *testing.T, client mercury.MercuryClient, eventContract EventContract) string {
	t.Helper()
	results, err := client.Emit("register-events::v2020_12_25", mercury.TargetAndPayload{
		Payload: map[string]any{
			"contract": eventContract,
		},
	})

	require.NoError(t, err, "Registering events should not return an error")
	require.NotNil(t, results, "Registering events should return results")

	first := results[0]
	fqenValues, ok := first["fqens"].([]any)
	require.True(t, ok, "FQENS slice should be present in response")
	require.NotEmpty(t, fqenValues, "FQENS slice should not be empty")

	fqen, ok := fqenValues[0].(string)
	require.True(t, ok, "First FQEN entry should be a string")
	return fqen
}

// GenerateWillSendVipEventSignature creates a sample contract for tests.
func GenerateWillSendVipEventSignature(slug ...string) EventContract {
	namespace := ""
	if len(slug) > 0 && slug[0] != "" {
		namespace = slug[0] + "."
	}

	return EventContract{
		"eventSignatures": map[string]any{
			fmt.Sprintf("%swill-send-vip::v1", namespace): map[string]any{
				"emitPayloadSchema": map[string]any{
					"id": "willSendVipTargetAndPayload",
					"fields": map[string]any{
						"target": map[string]any{
							"type":       "schema",
							"isRequired": true,
							"options": map[string]any{
								"schema": map[string]any{
									"id": "willSendVipTarget",
									"fields": map[string]any{
										"organizationId": map[string]any{
											"type": "text",
										},
									},
								},
							},
						},
						"payload": map[string]any{
							"type":       "schema",
							"isRequired": true,
							"options": map[string]any{
								"schema": map[string]any{
									"id": "willSendVipPayload",
									"fields": map[string]any{
										"message": map[string]any{
											"type": "text",
										},
									},
								},
							},
						},
					},
				},
				"responsePayloadSchema": map[string]any{
					"id": "testEventResponsePayload",
					"fields": map[string]any{
						"messages": map[string]any{
							"type":       "text",
							"isArray":    true,
							"isRequired": true,
						},
					},
				},
			},
		},
	}
}

// RegisterTestContract registers the default test contract and returns the FQEN.
func RegisterTestContract(t *testing.T, client mercury.MercuryClient) string {
	t.Helper()
	eventContract := GenerateWillSendVipEventSignature()
	require.NotNil(t, eventContract, "Expected event contract to be generated")

	fqen := RegisterEvents(t, client, eventContract)
	fmt.Println("Registered FQEN:", fqen)
	return fqen
}

// LoginCreateOrgSetupTwoSkills authenticates as a person, creates an org, installs two skills, and registers a contract.
func LoginCreateOrgSetupTwoSkills(t *testing.T) (*spruce.Organization, mercury.MercuryClient, mercury.MercuryClient, string) {
	t.Helper()
	client, _, _ := LoginAsDemoPerson(t, "+1 555-555-5555")
	org := SeedRandomOrg(t, client)
	skill1Client := seedSkillInstallToOrgAndLoginAsSkill(t, client, org)
	skill2Client := seedSkillInstallToOrgAndLoginAsSkill(t, client, org)
	fqen := RegisterTestContract(t, skill1Client)

	client.Disconnect()
	return org, skill1Client, skill2Client, fqen
}

func seedSkillInstallToOrgAndLoginAsSkill(t *testing.T, personClient mercury.MercuryClient, org *spruce.Organization) mercury.MercuryClient {
	t.Helper()
	skill, err := SeedRandomSkill(personClient)
	require.NoError(t, err, "Seeding skill should not return an error")
	fmt.Println("Seeded skill:", skill)
	err = InstallSkill(personClient, org.Id, skill.Id)
	require.NoError(t, err, "Installing skill should not return an error")

	skillClient, err := LoginAsSkill(t, skill)
	require.NoError(t, err, "Logging in as skill should not return an error")
	return skillClient
}
