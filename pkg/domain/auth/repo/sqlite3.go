package repo

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"time"

	model "github.com/pawelkuk/pictura-certamine/pkg/domain/auth/model"
)

type SQLiteRepo struct {
	DB *sql.DB
}

func (r *SQLiteRepo) Create(ctx context.Context, session *model.Session) error {
	_, err := r.DB.ExecContext(ctx,
		`insert into session(user_id, token, expiry) values(?, ?, ?)`,
		session.UserID, session.Token.Value, session.Expiry.Format(time.RFC3339),
	)
	if err != nil {
		return fmt.Errorf("could not create session: %w", err)
	}
	return nil
}
func (r *SQLiteRepo) Read(ctx context.Context, session *model.Session) error {
	row := r.DB.QueryRowContext(ctx,
		"select token, expiry, user_id from session where token = ?", session.Token.Value)
	if row.Err() != nil {
		return fmt.Errorf("could not query session with token=%s: %w", session.Token.Value, row.Err())
	}
	var token, expiry string
	var userid int64
	err := row.Scan(&token, &expiry, &userid)
	if err != nil {
		return fmt.Errorf("could not scan row: %w", err)
	}
	session.Token.Value = token
	expiryTime, err := time.Parse(time.RFC3339, expiry)
	if err != nil {
		return fmt.Errorf("could not parse expiry time: %w", err)
	}
	session.Expiry = expiryTime
	session.UserID = userid
	return nil
}
func (r *SQLiteRepo) Update(ctx context.Context, session *model.Session) error {
	_, err := r.DB.ExecContext(ctx,
		`update session set expiry = ?, user_id = ? where token = ?`,
		session.Expiry.Format(time.RFC3339), session.UserID, session.Token.Value,
	)
	if err != nil {
		return fmt.Errorf("could not update session: %w", err)
	}
	return nil
}
func (r *SQLiteRepo) Delete(ctx context.Context, session *model.Session) error {
	_, err := r.DB.ExecContext(ctx, `delete from session where token = ?`, session.Token.Value)
	if err != nil {
		return fmt.Errorf("could not delete session: %w", err)
	}
	return nil
}

func (r *SQLiteRepo) Query(ctx context.Context, filter model.QueryFilter) ([]model.Session, error) {
	q := `
	SELECT
		token, user_id, expiry
	FROM session
	`
	buf := bytes.NewBufferString(q)
	args := []any{}
	applyFilter(filter, &args, buf)
	rows, err := r.DB.QueryContext(ctx, buf.String(), args...)
	if err != nil {
		return nil, fmt.Errorf("could not execute query: %w", err)
	}
	sessions := []model.Session{}
	for rows.Next() {
		s := &model.Session{}
		var token, expiry string
		var userid int64
		err := rows.Scan(&token, &userid, &expiry)
		if err != nil {
			return nil, fmt.Errorf("could not scan row: %w", err)
		}
		s.Token.Value = token
		expiryTime, err := time.Parse(time.RFC3339, expiry)
		if err != nil {
			return nil, fmt.Errorf("could not parse expiry time: %w", err)
		}
		s.Expiry = expiryTime
		s.UserID = userid
		sessions = append(sessions, *s)
	}
	err = rows.Close()
	if err != nil {
		return nil, fmt.Errorf("could not close rows: %w", err)
	}
	return sessions, nil
}
