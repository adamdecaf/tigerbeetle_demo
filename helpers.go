package tigerbeetle_demo

import (
	"testing"

	tigerbeetle_go "github.com/tigerbeetle/tigerbeetle-go"
	"github.com/tigerbeetle/tigerbeetle-go/pkg/types"

	"github.com/stretchr/testify/require"
)

func setupClient(t *testing.T) tigerbeetle_go.Client {
	t.Helper()

	addresses := []string{"127.0.0.1:3002", "127.0.0.1:3001", "127.0.0.1:3003"}

	client, err := tigerbeetle_go.NewClient(types.ToUint128(0), addresses)
	require.NoError(t, err)

	t.Cleanup(func() {
		client.Close()
	})

	// ping the cluster
	require.NoError(t, client.Nop())

	return client
}
