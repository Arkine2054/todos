package storer

import (
	"context"
	"fmt"
)

func (ps *PostgresStorer) CreateSession(ctx context.Context, s *Session) (*Session, error) {
	res, err := ps.db.PrepareNamedContext(ctx,
		"INSERT INTO sessions (id, user_email, refresh_token, is_revoked, created_at, expires_at ) VALUES (:id, :user_email, :refresh_token, :is_revoked, :created_at, :expires_at) RETURNING id")

	if err != nil {
		return nil, fmt.Errorf("error inserting todos: %w", err)
	}
	err = res.QueryRow(s).Scan(&s.ID)
	if err != nil {
		return nil, fmt.Errorf("error inserting session: %w", err)
	}

	return s, nil
}

func (ps *PostgresStorer) GetSession(ctx context.Context, id string) (*Session, error) {
	var s Session
	err := ps.db.GetContext(ctx, &s,
		"SELECT * FROM sessions WHERE id=$1", id)
	if err != nil {
		return nil, fmt.Errorf("error getting session: %w", err)
	}

	return &s, nil
}

func (ps *PostgresStorer) RevokeSession(ctx context.Context, id string) error {
	_, err := ps.db.NamedExecContext(ctx,
		"UPDATE sessions SET is_revoked=true WHERE id=:id", map[string]interface{}{"id": id})
	if err != nil {
		return fmt.Errorf("error revoking session: %w", err)
	}

	return nil
}

func (ps *PostgresStorer) DeleteSession(ctx context.Context, id string) error {
	_, err := ps.db.ExecContext(ctx,
		"DELETE FROM sessions WHERE id=?", id)
	if err != nil {
		return fmt.Errorf("error deleting session: %w", err)
	}

	return nil
}
