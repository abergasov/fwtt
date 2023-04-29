package validator_test

import (
	"fwtt/internal/entites"
	"fwtt/internal/testhelpers"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestRepository_SaveChallenges(t *testing.T) {
	container := testhelpers.GetClean(t)
	challenges := []*entites.Challenge{
		{
			ValidTill: time.Now().Add(time.Hour).Unix(),
			Challenge: uuid.NewString(),
		},
		{
			ValidTill: time.Now().Add(time.Hour).Unix(),
			Challenge: uuid.NewString(),
		},
		{
			ValidTill: time.Now().Add(time.Hour).Unix(),
			Challenge: uuid.NewString(),
		},
		{
			ValidTill: time.Now().Add(time.Hour).Unix(),
			Challenge: uuid.NewString(),
		},
	}

	require.NoError(t, container.RepoValidator.SaveChallenges(challenges))
	require.NoError(t, container.RepoValidator.DropChallenges([]string{
		challenges[1].Challenge,
		challenges[2].Challenge,
	}))

	loadedChallenges, err := container.RepoValidator.LoadChallenges()
	require.NoError(t, err)
	require.Len(t, loadedChallenges, 2)
	require.Contains(t, []string{challenges[0].Challenge, challenges[3].Challenge}, loadedChallenges[0].Challenge)
	require.Contains(t, []string{challenges[0].Challenge, challenges[3].Challenge}, loadedChallenges[1].Challenge)
}
