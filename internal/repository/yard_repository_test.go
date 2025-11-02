package repository

import (
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestYardRepository_GetByCode(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewYardRepository(db)

	t.Run("success", func(t *testing.T) {
		now := time.Now()

		rows := sqlmock.NewRows([]string{"id", "code", "name", "description", "created_at", "updated_at"}).
			AddRow(1, "YRD1", "Yard 1", "Main yard", now, now)

		mock.ExpectQuery("SELECT id, code, name, description, created_at, updated_at FROM yards WHERE code = \\$1").
			WithArgs("YRD1").
			WillReturnRows(rows)

		yard, err := repo.GetByCode("YRD1")
		assert.NoError(t, err)
		assert.NotNil(t, yard)
		assert.Equal(t, "YRD1", yard.Code)
		assert.Equal(t, "Yard 1", yard.Name)
	})

	t.Run("not found", func(t *testing.T) {
		mock.ExpectQuery("SELECT id, code, name, description, created_at, updated_at FROM yards WHERE code = \\$1").
			WithArgs("INVALID").
			WillReturnError(sql.ErrNoRows)

		yard, err := repo.GetByCode("INVALID")
		assert.Error(t, err)
		assert.Nil(t, yard)
	})

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestYardRepository_GetAll(t *testing.T) {
	time := time.Now()
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewYardRepository(db)

	rows := sqlmock.NewRows([]string{"id", "code", "name", "description", "created_at", "updated_at"}).
		AddRow(1, "YRD1", "Yard 1", "Main yard", time, time).
		AddRow(2, "YRD2", "Yard 2", "Secondary yard", time, time)

	mock.ExpectQuery("SELECT id, code, name, description, created_at, updated_at FROM yards ORDER BY code").
		WillReturnRows(rows)

	yards, err := repo.GetAll()
	assert.NoError(t, err)
	assert.Len(t, yards, 2)
	assert.Equal(t, "YRD1", yards[0].Code)
	assert.Equal(t, "YRD2", yards[1].Code)

	assert.NoError(t, mock.ExpectationsWereMet())
}
