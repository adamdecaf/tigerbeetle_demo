package tigerbeetle_demo

import (
	"log"
	"strings"
	"testing"
	"time"

	"github.com/tigerbeetle/tigerbeetle-go/pkg/types"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestBasic(t *testing.T) {
	client := setupClient(t)

	// Create two accounts
	accountID1 := convertUUID(t, uuid.NewString())
	accountID2 := convertUUID(t, uuid.NewString())

	res, err := client.CreateAccounts([]types.Account{
		{
			ID:     accountID1,
			Ledger: 1,
			Code:   1,
		},
		{
			ID:     accountID2,
			Ledger: 1,
			Code:   1,
		},
	})
	require.NoError(t, err)

	for _, err := range res {
		log.Fatalf("Error creating account %d: %s", err.Index, err.Result)
	}

	transferID := uint64(time.Now().UTC().UnixMilli())
	transferRes, err := client.CreateTransfers([]types.Transfer{
		{
			ID:              types.ToUint128(transferID),
			DebitAccountID:  accountID1,
			CreditAccountID: accountID2,
			Amount:          types.ToUint128(10),
			Ledger:          1,
			Code:            1,
		},
	})
	require.NoError(t, err)

	for _, err := range transferRes {
		log.Fatalf("Error creating transfer: %s", err.Result)
	}

	// Check the sums for both accounts
	accounts, err := client.LookupAccounts([]types.Uint128{types.ToUint128(1), types.ToUint128(2)})
	require.NoError(t, err)
	require.Len(t, accounts, 2)

	for _, account := range accounts {
		if account.ID == types.ToUint128(1) {
			require.Equal(t, types.ToUint128(10), account.DebitsPosted, "account 1 debits")
			require.Equal(t, types.ToUint128(0), account.CreditsPosted, "account 1 credits")
		} else if account.ID == types.ToUint128(2) {
			require.Equal(t, types.ToUint128(0), account.DebitsPosted, "account 2 debits")
			require.Equal(t, types.ToUint128(10), account.CreditsPosted, "account 2 credits")
		} else {
			log.Fatalf("Unexpected account")
		}
	}
}

func convertUUID(t *testing.T, input string) types.Uint128 {
	t.Helper()

	id, err := types.HexStringToUint128(strings.ReplaceAll(input, "-", ""))
	require.NoError(t, err)

	return id

}
