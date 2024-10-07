package nats

import (
	"errors"
	"testing"

	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestDefaultNATSConnector_Connect(t *testing.T) {
	// Start a NATS server
	ns, url := startNATSServer(t)
	defer ns.Shutdown()

	connector := &DefaultNATSConnector{}

	// Test successful connection
	conn, err := connector.Connect(url)
	require.NoError(t, err)
	assert.NotNil(t, conn)
	assert.Implements(t, (*ConnInterface)(nil), conn)

	// Close the connection
	conn.Close()

	// Test connection failure
	_, err = connector.Connect("nats://invalid-url:4222")
	assert.Error(t, err)
}

func TestDefaultJetStreamCreator_New(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("Successful JetStream creation", func(t *testing.T) {
		// Start a NATS server
		ns, url := startNATSServer(t)
		defer ns.Shutdown()

		// Create a real NATS connection
		nc, err := nats.Connect(url)
		require.NoError(t, err)
		defer nc.Close()

		creator := &DefaultJetStreamCreator{}

		// Test successful JetStream creation
		js, err := creator.New(nc)
		require.NoError(t, err)
		assert.NotNil(t, js)
	})

	t.Run("JetStream creation failure", func(t *testing.T) {
		// Create a mock NATS connection
		mockConn := NewMockConnInterface(ctrl)

		// Make the mock connection return an error when JetStream() is called
		mockConn.EXPECT().JetStream().Return(nil, errors.New("JetStream creation failed"))

		creator := &DefaultJetStreamCreator{}

		// Test JetStream creation failure
		_, err := creator.New(mockConn)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "JetStream creation failed")
	})
}
