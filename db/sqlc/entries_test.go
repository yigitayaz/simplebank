package db

import (
	"context"
	"simplebank/util"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func createRandomEntry(t *testing.T, acc Account) Entry {
	arg := CreateEntryParams{
		AccountID: acc.ID,
		Amount:    util.RandomMoney(),
	}

	entry, err := TestQueries.CreateEntry(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, entry)

	require.Equal(t, arg.AccountID, entry.AccountID)
	require.Equal(t, arg.Amount, entry.Amount)

	require.NotZero(t, entry.ID)
	require.NotZero(t, entry.CreatedAt)

	return entry

}

func TestCreateEntry(t *testing.T) {
	account := createRandomAccount(t)
	createRandomEntry(t, account)
}

func TestGetEntry(t *testing.T) {
	account := createRandomAccount(t)
	testEntry := createRandomEntry(t, account)
	dbEntry, err := TestQueries.GetEntry(context.Background(), testEntry.ID)

	require.NoError(t, err)
	require.NotEmpty(t, dbEntry)

	require.Equal(t, testEntry.ID, dbEntry.ID)
	require.Equal(t, testEntry.AccountID, dbEntry.AccountID)
	require.Equal(t, testEntry.Amount, dbEntry.Amount)
	require.WithinDuration(t, testEntry.CreatedAt, dbEntry.CreatedAt, time.Second)

}

func TestListEntries(t *testing.T) {
	account := createRandomAccount(t)

	for i := 0; i < 10; i++ {
		createRandomEntry(t, account)
	}

	arg := ListEntriesParams{
		AccountID: account.ID,
		Limit:     5,
		Offset:    5,
	}

	entries, err := TestQueries.ListEntries(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, entries, 5)

	for _, entry := range entries {
		require.NotEmpty(t, entry)
		require.Equal(t, arg.AccountID, entry.AccountID)
	}
}
