package storer

import (
	"context"
	"fmt"
)

func (ps *PostgresStorer) CreateUser(ctx context.Context, u *Users) (*Users, error) {
	res, err := ps.db.PrepareNamedContext(ctx,
		"INSERT INTO users (email, password_hash, is_admin) VALUES (:email, :password_hash, :is_admin) RETURNING id")
	if err != nil {
		return nil, fmt.Errorf("error inserting user: %w", err)
	}

	err = res.QueryRow(u).Scan(&u.ID)

	if err != nil {
		return nil, fmt.Errorf("error getting last insert ID: %w", err)
	}

	return u, nil
}

func (ps *PostgresStorer) GetUser(ctx context.Context, email string) (*Users, error) {
	var u Users
	err := ps.db.GetContext(ctx, &u,
		"SELECT id, email, password_hash, is_admin FROM users WHERE email=$1", email)
	if err != nil {
		ps.logger.Errorf("error getting user: %w", err)
		return nil, fmt.Errorf("error getting user: %w", err)
	}

	return &u, nil
}

func (ps *PostgresStorer) ListUsers(ctx context.Context, list List) ([]Users, error) {

	query := fmt.Sprintf(
		"SELECT id, email, password_hash, is_admin, created_at, updated_at FROM users ORDER BY %s %s LIMIT $1 OFFSET $2",
		list.Sort,
		list.Order,
	)

	var users []Users
	err := ps.db.SelectContext(ctx, &users, query, list.Limit, list.Offset)

	if err != nil {
		return nil, fmt.Errorf("error listing users: %w", err)
	}

	return users, nil
}

func (ps *PostgresStorer) UpdateUser(ctx context.Context, u *Users) (*Users, error) {
	_, err := ps.db.NamedExecContext(ctx,
		"UPDATE users SET email=:email, password_hash=:password_hash, is_admin=:is_admin, updated_at=:updated_at WHERE email=:email", u)
	if err != nil {
		return nil, fmt.Errorf("error updating user: %w", err)
	}

	return u, nil
}

func (ps *PostgresStorer) DeleteUser(ctx context.Context, id int) error {
	_, err := ps.db.ExecContext(ctx, "DELETE FROM users WHERE id=?", id)
	if err != nil {
		return fmt.Errorf("error deleting user: %w", err)
	}

	return nil
}
