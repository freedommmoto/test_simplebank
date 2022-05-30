package db

import (
	"context"
	"testing"
	"time"

	tool "github.com/freedommmoto/test_simplebank/tool"
	"github.com/stretchr/testify/assert"
)

func randomMakeEntries(t *testing.T) Entry {

	customer := randomMakeCustomer(t)
	arg := CreateEntriesParams{
		CustomerID: customer.ID,
		Amount:     tool.RandomMoney(),
	}

	Entry, err := testQueries.CreateEntries(context.Background(), arg)

	assert.NoError(t, err)
	assert.NotEmpty(t, Entry)

	assert.Equal(t, arg.CustomerID, Entry.CustomerID)
	assert.Equal(t, arg.Amount, Entry.Amount)

	assert.NotZero(t, Entry.ID)
	assert.NotZero(t, Entry.CreatedAt)
	return Entry
}

// run test|dubug testing
func TestCreateEntries(t *testing.T) {
	randomMakeEntries(t)
}

func TestGetEntrie(t *testing.T) {
	Entries := randomMakeEntries(t)
	EntriesList, err := testQueries.ListEntries(context.Background(), Entries.ID)
	assert.NoError(t, err)
	assert.NotEmpty(t, EntriesList)

	assert.Equal(t, EntriesList.CustomerID, Entries.CustomerID)
	assert.Equal(t, EntriesList.Amount, Entries.Amount)
	//t.Logf("list of EntriesList is %v is used Entries ID %v for select", EntriesList, Entries.ID)
	assert.WithinDuration(t, EntriesList.CreatedAt, Entries.CreatedAt, time.Second)
}

func TestGetEntrieByCustomerID(t *testing.T) {
	Entries := randomMakeEntries(t)
	EntriesList, err := testQueries.ListEntriesByCustomerID(context.Background(), Entries.CustomerID)

	assert.NoError(t, err)
	for _, Entry := range EntriesList {
		assert.NotEmpty(t, Entry.ID)
	}
}
