package testkit

import (
	"fmt"
	"testing"

	"github.com/sprucelabsai-community/mercury-client-go/pkg/mercury"
	"github.com/stretchr/testify/require"
)

func TestFactory(t *testing.T) {
	t.Run("can easily create test client", func(t *testing.T) {
		BeforeEach(t)
		_, err := mercury.MakeMercuryClient()
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
		client, err := mercury.MakeMercuryClient()
		require.NoError(t, err, "Should not have error creating client")

		_, err = client.Emit("unregistered-event", mercury.TargetAndPayload{})
		require.Error(t, err, "Emitting to unregistered event should return an error")
	})

	t.Run("fake sockets can emit to each other", func(t *testing.T) {
		BeforeEach(t)
		client1, err := mercury.MakeMercuryClient()
		require.NoError(t, err)

		client2, err := mercury.MakeMercuryClient()
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

		client, err := mercury.MakeMercuryClient()
		require.NoError(t, err, "Should not have error creating client")

		_, err = client.Emit("unregistered-event", mercury.TargetAndPayload{})
		require.Error(t, err, "Emitting to unregistered event should return an error")
	})

	t.Run("only emits to the last listener set for an event", func(t *testing.T) {
		BeforeEach(t)

		client, err := mercury.MakeMercuryClient()
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
}
