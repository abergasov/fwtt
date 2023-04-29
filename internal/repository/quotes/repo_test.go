package quotes_test

import (
	"fwtt/internal/testhelpers"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRepository_LoadQuotes(t *testing.T) {
	container := testhelpers.GetClean(t)
	quotes, err := container.RepoQuotes.LoadQuotes()
	require.NoError(t, err)
	require.Len(t, quotes, 10)
}
