package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/simplebank/util"
	"github.com/stretchr/testify/require"
)

func createRandomTransfer(t *testing.T) Transfer {
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	arg := CreateTransferParams{
		FromAccountID: account1.ID,
		ToAccountID:   account2.ID,
		Amount:        util.RandomMoney(),
	}
	transfer, err := testQueries.CreateTransfer(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, transfer)

	require.Equal(t, arg.FromAccountID, transfer.FromAccountID)
	require.Equal(t, arg.ToAccountID, transfer.ToAccountID)
	require.Equal(t, arg.Amount, transfer.Amount)

	require.NotZero(t, transfer.ID)
	require.NotZero(t, transfer.CreatedAt)

	return transfer
}

func TestCreateTransfer(t *testing.T) {
	createRandomTransfer(t)
}

func TestGetTransfer(t *testing.T) {
	account1 := createRandomTransfer(t)
	account2, err := testQueries.GetTransfer(context.Background(), account1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, account2)

	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, account1.FromAccountID, account2.FromAccountID)
	require.Equal(t, account1.ToAccountID, account2.ToAccountID)
	require.Equal(t, account1.Amount, account2.Amount)
	require.WithinDuration(t, account1.CreatedAt, account2.CreatedAt, time.Second)
}

func TestUpdateTransfer(t *testing.T) {
	account1 := createRandomTransfer(t)

	arg := UpdateTransferParams{
		ID:          account1.ID,
		ToAccountID: util.RandomMoney(),
	}

	account2, err := testQueries.UpdateTransfer(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account2)

	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, account1.FromAccountID, account2.FromAccountID)
	require.Equal(t, arg.ToAccountID, account2.ToAccountID)
	require.Equal(t, account1.Amount, account2.Amount)
	require.WithinDuration(t, account1.CreatedAt, account2.CreatedAt, time.Second)
}

func TestDeleteTransfer(t *testing.T) {
	account1 := createRandomTransfer(t)
	err := testQueries.DeleteTransfer(context.Background(), account1.ID)
	require.NoError(t, err)

	account2, err := testQueries.GetTransfer(context.Background(), account1.ID)
	require.Error(t, err)
	require.Error(t, err, sql.ErrNoRows.Error())
	require.Empty(t, account2)
}

func TestListTransfers(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomTransfer(t)
	}

	arg := ListAccounsParams{
		Limit:  5,
		Offset: 5,
	}
	accounts, err := testQueries.ListAccouns(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, accounts, 5)
	for _, account := range accounts {
		require.NotEmpty(t, account)
	}
}
