package postgresql

import (
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"

	"github.com/quantonganh/ssr"
)

func TestRepositoryService(t *testing.T) {
	t.Run("create repo", testCreateRepo)
	t.Run("get repo", testGetRepo)
}

func testCreateRepo(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := &ssr.Repository{
		Provider:    "GitHub",
		FullName:    "quantonganh/ssr",
		Description: "Security scan result",
	}
	mock.ExpectExec(regexp.QuoteMeta(sqlInsertRepository)).WithArgs(repo.Provider, repo.FullName, repo.Description).WillReturnResult(sqlmock.NewResult(1, 1))

	repoService := NewRepositoryService(db)
	require.NoError(t, repoService.Create(repo))
}

func testGetRepo(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := &ssr.Repository{
		ID: 1,
		Provider:    "GitHub",
		FullName:    "quantonganh/ssr",
		Description: "Security scan result",
	}
	rows := sqlmock.NewRows([]string{"provider", "full_name", "description"}).AddRow(repo.Provider, repo.FullName, repo.Description)
	mock.ExpectQuery(regexp.QuoteMeta(sqlSelectRepository)).WithArgs(repo.ID).WillReturnRows(rows)

	repoService := NewRepositoryService(db)
	r, err := repoService.Get(1)
	require.NoError(t, err)
	assert.Equal(t, "GitHub", r.Provider)
	assert.Equal(t, "quantonganh/ssr", r.FullName)
	assert.Equal(t, "Security scan result", r.Description)
}
