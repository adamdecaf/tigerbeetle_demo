package tigerbeetle_demo

import (
	"math"
	"math/big"
	"math/rand/v2"
	"testing"

	tigerbeetle_go "github.com/tigerbeetle/tigerbeetle-go"
	"github.com/tigerbeetle/tigerbeetle-go/pkg/types"

	"github.com/stretchr/testify/require"
)

func TestComplex(t *testing.T) {
	client := setupClient(t)

	accounts := makeRandomAcccounts(t, client)
	require.Len(t, accounts, numbAccounts)

	transfers := makeRandomTransfers(t, client, accounts)
	require.Len(t, transfers, numbTransfers)

	findMinMaxAccountBalances(t, client, accounts)
}

const (
	numbAccounts  = 5000
	numbTransfers = 250_000

	batchSize = 1000

	ledger = 1
)

var (
	AccountAsset     uint16 = 100 // 101 could be Bank Account, 102 could be Money Market
	AccountLiability uint16 = 200
	AccountEquity    uint16 = 300
	AccountIncome    uint16 = 400
	AccountExpense   uint16 = 500

	TransferPurchase uint16 = 1000 // can represent anything
	TransferRefund   uint16 = 2000
	TransferFx       uint16 = 3000
)

var (
	accountTypes = []uint16{AccountAsset, AccountLiability, AccountEquity, AccountIncome, AccountExpense}

	transferTypes = []uint16{TransferPurchase, TransferRefund, TransferFx}
)

func makeRandomAcccounts(t *testing.T, client tigerbeetle_go.Client) []types.Account {
	t.Helper()

	out := make([]types.Account, numbAccounts)
	for i := 0; i < numbAccounts; i++ {
		// Create a new account and assign it to out[i]
		out[i] = types.Account{
			ID:     types.ID(),
			Ledger: ledger,
			Code:   randomItem(accountTypes),
		}

		// Process accounts in batches
		if (i+1)%batchSize == 0 || i == numbAccounts-1 {
			start := i - (i % batchSize) // Start of the current batch
			if i == numbAccounts-1 && numbAccounts%batchSize != 0 {
				start = (numbAccounts / batchSize) * batchSize // Handle partial last batch
			}

			res, err := client.CreateAccounts(out[start : i+1])
			require.NoError(t, err)

			for j, resp := range res {
				if resp.Result != types.AccountOK {
					t.Errorf("Account %d failed: %s", start+j, resp.Result.String())
				}
			}
		}
	}

	return out
}

func makeRandomTransfers(t *testing.T, client tigerbeetle_go.Client, accounts []types.Account) []types.Transfer {
	t.Helper()

	out := make([]types.Transfer, numbTransfers)
	for i := 0; i < numbTransfers; i++ {
		// Select random debit and credit accounts (ensure they are different)
		var debitAccount, creditAccount types.Account
		for {
			debitAccount = randomItem(accounts)
			creditAccount = randomItem(accounts)

			if debitAccount != creditAccount {
				break
			}
		}

		// Create a new transfer and assign it to out[i]
		out[i] = types.Transfer{
			ID:              types.ID(),
			DebitAccountID:  debitAccount.ID,
			CreditAccountID: creditAccount.ID,
			Ledger:          ledger,
			Code:            randomItem(transferTypes),
			Amount:          types.ToUint128(rand.Uint64N(1000) + 1),
		}

		// Process transfers in batches
		if (i+1)%batchSize == 0 || i == numbTransfers-1 {
			start := i - (i % batchSize) // Start of the current batch
			if i == numbTransfers-1 && numbTransfers%batchSize != 0 {
				start = (numbTransfers / batchSize) * batchSize // Handle partial last batch
			}

			res, err := client.CreateTransfers(out[start : i+1])
			require.NoError(t, err)

			for j, resp := range res {
				if resp.Result != types.TransferOK {
					t.Errorf("Transfer %d failed: %s", start+j, resp.Result.String())
				}
			}
		}
	}

	return out
}

// AccountBalance holds an account and its calculated balance
type AccountBalance struct {
	Account types.Account
	Balance *big.Int
}

// findMinMaxAccountBalances finds the accounts with the largest and smallest balances
func findMinMaxAccountBalances(t *testing.T, client tigerbeetle_go.Client, accounts []types.Account) {
	t.Helper()

	// Lookup all accounts to get their balance details
	accountIDs := make([]types.Uint128, len(accounts))
	for i, acc := range accounts {
		accountIDs[i] = acc.ID
	}
	balances, err := client.LookupAccounts(accountIDs)
	require.NoError(t, err)
	require.Len(t, balances, len(accounts), "Expected all accounts to be found")

	// Initialize max and min balances
	maxBalance := big.NewInt(math.MinInt64)
	minBalance := big.NewInt(math.MaxInt64)
	var maxAcc, minAcc types.Account

	// For big.Int arithmetic
	var balance, debits, credits big.Int

	// Calculate balances based on account type
	for _, acc := range balances {
		// Convert DebitsPosted and CreditsPosted to big.Int (assuming they are uint64)
		// If they are Uint128, use acc.DebitsPosted.BigInt() and acc.CreditsPosted.BigInt()
		debits = acc.DebitsPosted.BigInt()
		credits = acc.CreditsPosted.BigInt()

		// Calculate balance based on account type
		switch acc.Code {
		case AccountAsset:
			// Asset: balance = debits_posted - credits_posted
			balance.Sub(&debits, &credits)

		case AccountLiability, AccountEquity, AccountIncome, AccountExpense:
			// Liability/Equity/Income/Expense: balance = credits_posted - debits_posted
			balance.Sub(&credits, &debits)

		default:
			t.Fatalf("Unknown account type code: %d", acc.Code)
		}

		// Update max and min balances
		if balance.Cmp(maxBalance) > 0 {
			maxBalance.Set(&balance)
			maxAcc = acc
		}
		if balance.Cmp(minBalance) < 0 {
			minBalance.Set(&balance)
			minAcc = acc
		}
	}

	// Convert account IDs to big.Int for return if needed
	maxID := maxAcc.ID.BigInt()
	minID := minAcc.ID.BigInt()

	// Log IDs as big.Int for clarity
	t.Logf("Max balance account ID (as big.Int): %s (balance: %v)", maxID.String(), maxBalance.String())
	t.Logf("Min balance account ID (as big.Int): %s (balance: %v)", minID.String(), minBalance.String())
}

func randomItem[T any](items []T) T {
	idx := rand.Int32N(int32(len(items)))
	return items[idx]
}
