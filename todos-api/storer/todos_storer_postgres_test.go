package storer

import (
	"context"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestPostgresStorer_CreateTodos(t *testing.T) {
	td := &ToDos{
		Title:       "test title",
		Description: "test description",
	}

	tcs := []struct {
		name string
		test func(*testing.T, *PostgresStorer, sqlmock.Sqlmock)
	}{
		{
			name: "success",
			test: func(t *testing.T, st *PostgresStorer, mock sqlmock.Sqlmock) {
				mock.ExpectPrepare("INSERT INTO todos (title, description, user_id) VALUES (?, ?, ?) RETURNING id").
					ExpectQuery().WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(td.ID))
				ctd, err := st.CreateTodos(context.Background(), td)
				require.NoError(t, err)
				require.Equal(t, td.ID, ctd.ID)
				err = mock.ExpectationsWereMet()
				require.NoError(t, err)
			},
		},
		{
			name: "failed inserting todos",
			test: func(t *testing.T, st *PostgresStorer, mock sqlmock.Sqlmock) {
				mock.ExpectPrepare("INSERT INTO todos (title, description, user_id) VALUES (?, ?, ?) RETURNING id").
					WillReturnError(fmt.Errorf("error inserting todos"))
				_, err := st.CreateTodos(context.Background(), td)
				require.Error(t, err)
				err = mock.ExpectationsWereMet()
				require.NoError(t, err)
			},
		},
		{
			name: "failed getting last insert ID",
			test: func(t *testing.T, st *PostgresStorer, mock sqlmock.Sqlmock) {
				mock.ExpectPrepare("INSERT INTO todos (title, description, user_id) VALUES (?, ?, ?) RETURNING id").
					WillReturnError(fmt.Errorf("error getting last insert ID"))
				_, err := st.CreateTodos(context.Background(), td)
				require.Error(t, err)
				err = mock.ExpectationsWereMet()
				require.NoError(t, err)
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			withTestDB(t, func(db *sqlx.DB, mock sqlmock.Sqlmock) {
				st := NewPostgresStorer(db)
				tc.test(t, st, mock)
			})
		})
	}
}

func TestPostgresStorer_GetTodos(t *testing.T) {

	td := &ToDos{
		Title:       "test title",
		Description: "test description",
	}

	tcs := []struct {
		name string
		test func(*testing.T, *PostgresStorer, sqlmock.Sqlmock)
	}{
		{
			name: "success",
			test: func(t *testing.T, st *PostgresStorer, mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "user_id", "title", "completed", "description"}).
					AddRow(1, 1, td.Title, td.Completed, td.Description)

				mock.ExpectQuery("SELECT id, user_id, title, completed, description FROM todos WHERE user_id=$1 AND id=$2").WithArgs(1, 1).WillReturnRows(rows)

				gt, err := st.GetTodos(context.Background(), 1, 1)
				require.NoError(t, err)
				require.Equal(t, 1, gt.UserID, gt.ID)

				err = mock.ExpectationsWereMet()
				require.NoError(t, err)
			},
		}, {
			name: "failed getting todos",
			test: func(t *testing.T, st *PostgresStorer, mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT id, user_id, title, completed, description FROM todos WHERE user_id=$1 AND id=$2").
					WithArgs(1, 1).WillReturnError(fmt.Errorf("error getting todos"))
				_, err := st.GetTodos(context.Background(), 1, 1)
				require.Error(t, err)

				err = mock.ExpectationsWereMet()
				require.NoError(t, err)
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			withTestDB(t, func(db *sqlx.DB, mock sqlmock.Sqlmock) {
				st := NewPostgresStorer(db)
				tc.test(t, st, mock)
			})
		})
	}
}

func TestPostgresStorer_ListTodos(t *testing.T) {

	td := &ToDos{
		Title:       "test title",
		Description: "test description",
	}

	tcs := []struct {
		name string
		test func(*testing.T, *PostgresStorer, sqlmock.Sqlmock)
	}{
		{
			name: "success",
			test: func(t *testing.T, st *PostgresStorer, mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "title", "description", "completed", "created_at", "updated_at"}).
					AddRow(1, td.Title, td.Description, td.Completed, td.CreatedAt, td.UpdatedAt)
				mock.ExpectQuery("SELECT id, title, description, completed, created_at, updated_at FROM todos").
					WillReturnRows(rows)

				todos, err := st.ListTodos(context.Background())
				require.NoError(t, err)
				require.Len(t, todos, 1)

				err = mock.ExpectationsWereMet()
				require.NoError(t, err)
			},
		}, {
			name: "failed querying todos",
			test: func(t *testing.T, st *PostgresStorer, mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT id, title, description, completed, created_at, updated_at FROM todos").
					WillReturnError(fmt.Errorf("error querying todos"))

				_, err := st.ListTodos(context.Background())
				require.Error(t, err)

				err = mock.ExpectationsWereMet()
				require.NoError(t, err)
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			withTestDB(t, func(db *sqlx.DB, mock sqlmock.Sqlmock) {
				st := NewPostgresStorer(db)
				tc.test(t, st, mock)
			})
		})
	}
}

func TestPostgresStorer_UpdateTodos(t *testing.T) {

	td := &ToDos{
		ID:          1,
		Title:       "test title",
		Description: "test description",
	}
	ntd := &ToDos{
		ID:          1,
		Title:       "new test title",
		Description: "test description",
	}

	tcs := []struct {
		name string
		test func(*testing.T, *PostgresStorer, sqlmock.Sqlmock)
	}{
		{
			name: "success",
			test: func(t *testing.T, st *PostgresStorer, mock sqlmock.Sqlmock) {
				mock.ExpectPrepare("INSERT INTO todos (title, description, user_id) VALUES (?, ?, ?) RETURNING id").
					ExpectQuery().WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(td.ID))
				ctd, err := st.CreateTodos(context.Background(), td)
				require.NoError(t, err)
				require.Equal(t, td.ID, ctd.ID)

				mock.ExpectExec("UPDATE todos SET title=?, description=?, completed=?, updated_at=? WHERE id=?").
					WillReturnResult(sqlmock.NewResult(1, 1))
				up, err := st.UpdateTodos(context.Background(), ntd)
				require.NoError(t, err)
				require.Equal(t, 1, up.ID)
				require.Equal(t, ntd.Title, up.Title)

				err = mock.ExpectationsWereMet()
				require.NoError(t, err)
			},
		},
		{
			name: "failed updating todos",
			test: func(t *testing.T, st *PostgresStorer, mock sqlmock.Sqlmock) {
				mock.ExpectExec("UPDATE todos SET title=?, description=?, completed=?, updated_at=? WHERE id=?").
					WillReturnError(fmt.Errorf("errpr updating todos"))
				_, err := st.UpdateTodos(context.Background(), td)
				require.Error(t, err)

				err = mock.ExpectationsWereMet()
				require.NoError(t, err)
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			withTestDB(t, func(db *sqlx.DB, mock sqlmock.Sqlmock) {
				st := NewPostgresStorer(db)
				tc.test(t, st, mock)
			})
		})
	}
}

func TestPostgresStorer_DeleteTodos(t *testing.T) {
	tcs := []struct {
		name string
		test func(*testing.T, *PostgresStorer, sqlmock.Sqlmock)
	}{
		{
			name: "success",
			test: func(t *testing.T, st *PostgresStorer, mock sqlmock.Sqlmock) {
				mock.ExpectExec("DELETE FROM todos WHERE id = $1").
					WithArgs(1).WillReturnResult(sqlmock.NewResult(1, 1))
				err := st.DeleteTodos(context.Background(), 1)
				require.NoError(t, err)

				err = mock.ExpectationsWereMet()
				require.NoError(t, err)
			},
		},
		{
			name: "failed deleting todos",
			test: func(t *testing.T, st *PostgresStorer, mock sqlmock.Sqlmock) {
				mock.ExpectExec("DELETE FROM todos WHERE id = $1").
					WithArgs(1).WillReturnError(fmt.Errorf("error deleting todos"))
				err := st.DeleteTodos(context.Background(), 1)
				require.Error(t, err)

				err = mock.ExpectationsWereMet()
				require.NoError(t, err)
			},
		},
	}
	for _, tc := range tcs {
		withTestDB(t, func(db *sqlx.DB, mock sqlmock.Sqlmock) {
			st := NewPostgresStorer(db)
			tc.test(t, st, mock)
		})
	}

}
