package db

import (
	"context"
	"database/sql"
	"simplebank/util"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func createRandomAccount(t *testing.T) Account {
	arg := CreateAccountParams{
		Owner:    util.RandomOwner(),
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}

	account, err := TestQueries.CreateAccount(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)

	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)

	return account
}

func TestCreateAccount(t *testing.T) {
	createRandomAccount(t)

}

func TestGetAccount(t *testing.T) {
	account := createRandomAccount(t)
	dbAccount, err := TestQueries.GetAccount(context.Background(), account.ID)

	require.NoError(t, err)
	require.NotEmpty(t, dbAccount)

	require.Equal(t, account.ID, dbAccount.ID)
	require.Equal(t, account.Owner, dbAccount.Owner)
	require.Equal(t, account.Balance, dbAccount.Balance)
	require.Equal(t, account.Currency, dbAccount.Currency)

	require.WithinDuration(t, account.CreatedAt, dbAccount.CreatedAt, time.Second)

}

func TestUpdateAccount(t *testing.T) {
	testAccount := createRandomAccount(t)
	arg := UpdateAccountParams{
		ID:      testAccount.ID,
		Balance: util.RandomMoney(),
	}

	dbAccount, err := TestQueries.UpdateAccount(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, dbAccount)

	require.Equal(t, testAccount.ID, dbAccount.ID)
	require.Equal(t, testAccount.Owner, dbAccount.Owner)
	require.Equal(t, arg.Balance, dbAccount.Balance)
	require.Equal(t, testAccount.Currency, dbAccount.Currency)

	require.WithinDuration(t, testAccount.CreatedAt, dbAccount.CreatedAt, time.Second)
}

func TestDeleteAccount(t *testing.T) {
	testAccount := createRandomAccount(t)
	_, err := TestQueries.DeleteAccount(context.Background(), testAccount.ID)

	require.NoError(t, err)

	dbAccount, err := TestQueries.GetAccount(context.Background(), testAccount.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, dbAccount)
}

func TestListAccounts(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomAccount(t)
	}

	arg := ListAccountsParams{
		Limit:  5,
		Offset: 5,
	}

	accounts, err := TestQueries.ListAccounts(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, accounts, 5)

	for _, account := range accounts {
		require.NotEmpty(t, account)
	}
}
