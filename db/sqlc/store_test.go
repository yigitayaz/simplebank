package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(TestDB)

	fromAccount := createRandomAccount(t)
	toAccount := createRandomAccount(t)
	fmt.Println(">> before:", fromAccount.Balance, toAccount.Balance)

	numberOfTransactions := 5

	amount := int64(10)

	errs := make(chan error)
	results := make(chan TransferTxResult)

	for i := 0; i < numberOfTransactions; i++ {
		txName := fmt.Sprintf("tx %d", i+1)
		go func() {
			ctx := context.WithValue(context.Background(), txKey, txName)
			result, err := store.TransferTx(ctx, TransferTxParams{
				FromAccountID: fromAccount.ID,
				ToAccountID:   toAccount.ID,
				Amount:        amount,
			})
			errs <- err
			results <- result
		}()
	}

	existed := make(map[int]bool)
	for i := 0; i < numberOfTransactions; i++ {
		err := <-errs
		require.NoError(t, err)

		result := <-results
		require.NotEmpty(t, result)

		transfer := result.Transfer
		require.Equal(t, fromAccount.ID, transfer.FromAccountID)
		require.Equal(t, toAccount.ID, transfer.ToAccountID)
		require.Equal(t, amount, transfer.Amount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		_, err = store.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		fromEntry := result.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, fromAccount.ID, fromEntry.AccountID)
		require.Equal(t, -amount, fromEntry.Amount)
		require.NotZero(t, fromAccount.ID)
		require.NotZero(t, fromAccount.CreatedAt)

		_, err = store.GetEntry(context.Background(), fromEntry.ID)
		require.NoError(t, err)

		toEntry := result.ToEntry
		require.NotEmpty(t, toEntry)
		require.Equal(t, toAccount.ID, toEntry.AccountID)
		require.Equal(t, amount, toEntry.Amount)
		require.NotZero(t, toAccount.ID)
		require.NotZero(t, toAccount.CreatedAt)

		_, err = store.GetEntry(context.Background(), toEntry.ID)
		require.NoError(t, err)

		// check accounts

		testFromAccount := result.FromAccount

		require.NotEmpty(t, fromAccount)
		require.Equal(t, fromAccount.ID, testFromAccount.ID)

		testToAccount := result.ToAccount

		require.NotEmpty(t, testToAccount)
		require.Equal(t, toAccount.ID, testToAccount.ID)

		fmt.Println(">> tx:", testFromAccount.Balance, testToAccount.Balance)

		// checl accounts balance

		diffFrom := fromAccount.Balance - testFromAccount.Balance
		diffTo := testToAccount.Balance - toAccount.Balance

		require.Equal(t, diffFrom, diffTo)

		require.True(t, diffFrom > 0)
		require.True(t, diffFrom%amount == 0)

		k := int(diffFrom / amount)
		require.True(t, k >= 1 && k <= numberOfTransactions)
		require.NotContains(t, existed, k)
		existed[k] = true

	}

	updatedAccountFrom, err := TestQueries.GetAccount(context.Background(), fromAccount.ID)
	require.NoError(t, err)

	updatedAccountTo, err := TestQueries.GetAccount(context.Background(), toAccount.ID)
	require.NoError(t, err)
	fmt.Println(">> after:", updatedAccountFrom.Balance, updatedAccountTo.Balance)
	require.Equal(t, fromAccount.Balance-int64(numberOfTransactions)*amount, updatedAccountFrom.Balance)
	require.Equal(t, toAccount.Balance+int64(numberOfTransactions)*amount, updatedAccountTo.Balance)

}
