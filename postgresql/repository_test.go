package postgresql

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/quantonganh/ssr"
)

const (
	sqlInsertRepository = `INSERT INTO "repository" ("provider","full_name","description") VALUES ($1,$2,$3) RETURNING "id"`
	sqlSelectRepository = `SELECT * FROM "repository" WHERE id = $1 ORDER BY "repository"."id" LIMIT 1`
)

func TestRepositoryService(t *testing.T) {
	t.Run("create repo", testCreateRepo)
	t.Run("get repo", testGetRepo)
}

func testCreateRepo(t *testing.T) {
	sqlDB, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)
	defer sqlDB.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDB,
	}), &gorm.Config{
		SkipDefaultTransaction: true,
	})
	require.NoError(t, err)

	repo := &ssr.Repository{
		Provider:    "GitHub",
		FullName:    "quantonganh/ssr",
		Description: "Security scan result",
	}
	mock.ExpectQuery(sqlInsertRepository).
		WithArgs(repo.Provider, repo.FullName, repo.Description).
		WillReturnRows(sqlmock.NewRows([]string{"provider", "full_name", "description"}).AddRow(repo.Provider, repo.FullName, repo.Description))

	repoService := NewRepositoryService(gormDB)
	require.NoError(t, repoService.Create(repo))
}

func testGetRepo(t *testing.T) {
	sqlDB, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)
	defer sqlDB.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDB,
	}), &gorm.Config{})
	require.NoError(t, err)

	repo := &ssr.Repository{
		ID: 1,
		Provider:    "GitHub",
		FullName:    "quantonganh/ssr",
		Description: "Security scan result",
	}
	rows := sqlmock.NewRows([]string{"provider", "full_name", "description"}).AddRow(repo.Provider, repo.FullName, repo.Description)
	mock.ExpectQuery(sqlSelectRepository).WithArgs(repo.ID).WillReturnRows(rows)

	repoService := NewRepositoryService(gormDB)
	r, err := repoService.Get(1)
	require.NoError(t, err)
	assert.Equal(t, "GitHub", r.Provider)
	assert.Equal(t, "quantonganh/ssr", r.FullName)
	assert.Equal(t, "Security scan result", r.Description)
}
