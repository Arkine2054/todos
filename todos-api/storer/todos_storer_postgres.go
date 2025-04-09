package storer

import (
	"context"
	"fmt"
)

func (ps *PostgresStorer) CreateTodos(ctx context.Context, t *ToDos) (*ToDos, error) {
	res, err := ps.db.PrepareNamedContext(ctx,
		"INSERT INTO todos (title, description, user_id) VALUES (:title, :description, :user_id) RETURNING id")

	if err != nil {
		return nil, fmt.Errorf("error inserting todos: %w", err)
	}
	err = res.QueryRow(t).Scan(&t.ID)

	if err != nil {
		return nil, fmt.Errorf("error getting last insert ID: %w", err)
	}
	return t, nil
}

func (ps *PostgresStorer) GetTodos(ctx context.Context, userID, todosID int) (*ToDos, error) {
	var t ToDos
	err := ps.db.GetContext(ctx, &t,
		"SELECT id, user_id, title, completed, description FROM todos WHERE user_id=$1 AND id=$2", userID, todosID)

	if err != nil {
		return nil, fmt.Errorf("error getting todos: %w", err)
	}

	return &t, nil
}
func (ps *PostgresStorer) ListUserTodos(ctx context.Context, userID int, list List) ([]ToDos, error) {

	query := fmt.Sprintf(
		"SELECT id, title, description, completed, created_at, updated_at FROM todos WHERE user_id=$1 AND title ILIKE $2 ORDER BY %s %s LIMIT $3 OFFSET $4",
		list.Sort, list.Order)

	var todos []ToDos
	err := ps.db.SelectContext(ctx, &todos, query, userID, list.Title, list.Limit, list.Offset)
	if err != nil {
		return nil, fmt.Errorf("error listing todos: %w", err)
	}

	return todos, nil
}

func (ps *PostgresStorer) ListTodos(ctx context.Context) ([]ToDos, error) {
	var todos []ToDos

	err := ps.db.SelectContext(ctx, &todos,
		"SELECT id, title, description, completed, created_at, updated_at FROM todos")
	if err != nil {
		return nil, fmt.Errorf("error listing todos: %w", err)
	}
	return todos, nil
}

func (ps *PostgresStorer) UpdateTodos(ctx context.Context, t *ToDos) (*ToDos, error) {
	_, err := ps.db.NamedExecContext(ctx,
		"UPDATE todos SET title=:title, description=:description, completed=:completed, updated_at=:updated_at WHERE id=:id", t)
	if err != nil {
		return nil, fmt.Errorf("error updating todos: %w", err)
	}

	return t, nil
}

func (ps *PostgresStorer) DeleteTodos(ctx context.Context, id int) error {
	_, err := ps.db.ExecContext(ctx,
		"DELETE FROM todos WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("error deleting todos: %w", err)
	}
	return nil
}
