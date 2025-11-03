package testkit

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/sprucelabsai-community/mercury-client-go/pkg/mercury"
	"github.com/sprucelabsai-community/spruce-core-schemas/v41/pkg/schemas"
	"github.com/stretchr/testify/require"
)

func TestFactory(t *testing.T) {
	t.Run("can easily create test client", func(t *testing.T) {
		BeforeEach(t)
		_, err := mercury.NewMercuryClient()
		require.NoError(t, err, "Created socket just fine")
	})

	t.Run("can build aggregate response with single response", func(t *testing.T) {
		actual := BuildAggregateResponse([]mercury.ResponsePayload{
			{"Hello": "World"},
		})

		expected := mercury.MercuryAggregateResponse{
			TotalContracts: 1,
			TotalResponses: 1,
			TotalErrors:    0,
			Responses: []mercury.MercurySingleResponse{
				{
					ResponderRef: "fake-responder-1",
					Errors:       []any{},
					Payload:      map[string]any{"Hello": "World"},
				},
			},
		}

		require.Equal(t, expected, actual, "Aggregate response should match expected")

	})

	t.Run("can build aggregate response with different response payload", func(t *testing.T) {
		actual := BuildAggregateResponse([]mercury.ResponsePayload{
			{"Foo": "Bar", "Number": 42},
		})

		expected := mercury.MercuryAggregateResponse{
			TotalContracts: 1,
			TotalResponses: 1,
			TotalErrors:    0,
			Responses: []mercury.MercurySingleResponse{
				{
					ResponderRef: "fake-responder-1",
					Errors:       []any{},
					Payload:      map[string]any{"Foo": "Bar", "Number": 42},
				},
			},
		}

		require.Equal(t, expected, actual, "Aggregate response should match expected")
	})

	t.Run("can build aggregate with 2 response payloads", func(t *testing.T) {
		actual := BuildAggregateResponse([]mercury.ResponsePayload{
			{"First": "Response"},
			{"Second": "Response"},
		})

		expected := mercury.MercuryAggregateResponse{
			TotalContracts: 2,
			TotalResponses: 2,
			TotalErrors:    0,
			Responses: []mercury.MercurySingleResponse{
				{
					ResponderRef: "fake-responder-1",
					Errors:       []any{},
					Payload:      map[string]any{"First": "Response"},
				},
				{
					ResponderRef: "fake-responder-2",
					Errors:       []any{},
					Payload:      map[string]any{"Second": "Response"},
				},
			},
		}

		require.Equal(t, expected, actual, "Aggregate response should match expected")
	})

	t.Run("returns error if no listener for event is set on fake", func(t *testing.T) {
		BeforeEach(t)
		client, err := mercury.NewMercuryClient()
		require.NoError(t, err, "Should not have error creating client")

		_, err = client.Emit("unregistered-event", mercury.TargetAndPayload{})
		require.Error(t, err, "Emitting event without local listener should return an error")
	})

	t.Run("fake sockets can emit to each other", func(t *testing.T) {
		BeforeEach(t)
		client1, err := mercury.NewMercuryClient()
		require.NoError(t, err)

		client2, err := mercury.NewMercuryClient()
		require.NoError(t, err)

		var wasHit bool
		client1.On("test-event", func(targetAndPayload mercury.TargetAndPayload) any {
			wasHit = true
			fmt.Println("Client 1 received event with payload:")
			return nil
		})

		client2.Emit("test-event")
		require.True(t, wasHit, "Client 1 should have received the event emitted by Client 2")
	})

	t.Run("returns error if no listener for event is set on fake", func(t *testing.T) {
		BeforeEach(t)

		client, err := mercury.NewMercuryClient()
		require.NoError(t, err, "Should not have error creating client")

		_, err = client.Emit("unregistered-event", mercury.TargetAndPayload{})
		require.Error(t, err, "Emitting to unregistered event should return an error")
	})

	t.Run("only emits to the last listener set for an event", func(t *testing.T) {
		BeforeEach(t)

		client, err := mercury.NewMercuryClient()
		require.NoError(t, err, "Should not have error creating client")

		var firstHit bool
		client.On("another.event::v100", func(targetAndPayload mercury.TargetAndPayload) any {
			firstHit = true
			return nil
		})

		var secondHit bool
		client.On("another.event::v100", func(targetAndPayload mercury.TargetAndPayload) any {
			secondHit = true
			return nil
		})

		client.Emit("another.event::v100", mercury.TargetAndPayload{})

		require.False(t, firstHit, "First listener should not have been hit")
		require.True(t, secondHit, "Second listener should have been hit")
	})

	t.Run("can return map response 1 from listener", func(t *testing.T) {
		BeforeEach(t)

		client, err := mercury.NewMercuryClient()
		require.NoError(t, err, "Should not have error creating client")

		emitAndAssertResponsePassedBack(t, client, mercury.ResponsePayload{
			"Key": "Value",
		})

		emitAndAssertResponsePassedBack(t, client, mercury.ResponsePayload{
			"AnotherKey": "watermelon",
			"Cheese":     "gouda",
		})
	})

	t.Run("returning a core location maps to camelCase keys", func(t *testing.T) {
		BeforeEach(t)

		location := schemas.Location{
			Id:   GenerateRandomId(),
			Name: GenerateRandomId(),
		}

		expected, mapErr := StructToMap(location)
		require.NoError(t, mapErr, "Should be able to convert location struct to map")

		client, _ := mercury.NewMercuryClient()
		client.On("get-location::v100", func(_ mercury.TargetAndPayload) any {
			return map[string]any{
				"location": location,
			}
		})

		responses, err := client.Emit("get-location::v100")
		require.NoError(t, err, "Emitting to registered event should not return an error")
		require.Equal(t, expected, responses[0]["location"], "Did not pass back map of location")

	})
}

func emitAndAssertResponsePassedBack(t *testing.T, client mercury.MercuryClient, responsePayload mercury.ResponsePayload) {
	client.On("map.response.event::v1", func(targetAndPayload mercury.TargetAndPayload) any {
		return responsePayload
	})

	responses, err := client.Emit("map.response.event::v1")
	require.NoError(t, err, "Emitting to registered event should not return an error")
	require.Len(t, responses, 1, "Response should have one payload")
	require.Equal(t, responsePayload, responses[0], "Response payload should match expected")
}

func StructToMap[T any](value T) (map[string]any, error) {
	data, err := json.Marshal(value)
	if err != nil {
		return nil, err
	}

	var out map[string]any
	if err := json.Unmarshal(data, &out); err != nil {
		return nil, err
	}
	return out, nil
}
