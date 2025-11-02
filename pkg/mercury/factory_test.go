package mercury_test

import (
	"testing"
	"time"

	mercury "github.com/sprucelabsai-community/mercury-client-go/pkg/mercury"
	"github.com/sprucelabsai-community/mercury-client-go/pkg/testkit"
	"github.com/stretchr/testify/require"
)

func TestFactory(t *testing.T) {

	t.Run("can get back client", func(t *testing.T) {
		testkit.BeforeEachInternal(t)
		factory := mercury.Factory{}
		require.NotNil(t, factory, "factory should not be nil")
	})

	t.Run("sets expected default options when none are provided", func(t *testing.T) {
		testkit.BeforeEachInternal(t)
		fake, _, err := testkit.MakeFakeClient()
		require.NoError(t, err)
		opts := fake.GetOptions()
		require.NotNil(t, opts, "Options should not be nil")
		require.Equal(t, 10*time.Second, opts.Timeout(), "Timeout should be 10 seconds by default")
		require.True(t, opts.Reconnection(), "Reconnection should be true by default")
	})

	t.Run("can set reconnect to false", func(t *testing.T) {
		testkit.BeforeEachInternal(t)
		fake, _, err := testkit.MakeFakeClient(mercury.MercuryClientOptions{ShouldRetryConnect: false, Host: "http://waka-waka"})
		require.NoError(t, err)
		opts := fake.GetOptions()
		require.NotNil(t, opts, "Options should not be nil")
		require.True(t, opts.Reconnection(), "Reconnection should be false when set to false")
		require.Equal(t, "http://waka-waka", fake.GetHost(), "Host should be set to http://waka-waka")
	})

	t.Run("returns error with bad url 1", func(t *testing.T) {
		testkit.BeforeEachInternal(t)
		_, err := mercury.MakeMercuryClient(mercury.MercuryClientOptions{Host: "aoeuao://bad-url", ShouldRetryConnect: false})
		require.Error(t, err, "Bad url should have returned an error")
	})

	t.Run("returns error with bad url 2", func(t *testing.T) {
		testkit.BeforeEachInternal(t)
		_, err := mercury.MakeMercuryClient(mercury.MercuryClientOptions{Host: "enon://aoeu333another-bad-uaoeuaoeurl", ShouldRetryConnect: false})
		require.Error(t, err, "Bad url should have returned an error")
	})

	t.Run("IsConnected calls method on socket client", func(t *testing.T) {
		testkit.BeforeEachInternal(t)
		fake, client, err := testkit.MakeFakeClient()
		require.NoError(t, err)

		fake.SetConnected(true)
		require.True(t, client.IsConnected(), "Client should be connected when socket client is connected")

		fake.SetConnected(false)
		require.False(t, client.IsConnected(), "Client should not be connected when socket client is not connected")

	})

	t.Run("defaults host is https_//mercury.spruce.ai", func(t *testing.T) {
		testkit.BeforeEachInternal(t)
		fake, _, err := testkit.MakeFakeClient()
		require.NoError(t, err)
		require.Equal(t, "https://mercury.spruce.ai", fake.GetHost(), "Default host should be https://mercury.spruce.ai")
	})

}
