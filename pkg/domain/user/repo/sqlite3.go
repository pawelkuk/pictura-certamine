package repo

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"net/mail"

	model "github.com/pawelkuk/pictura-certamine/pkg/domain/user/model"
)

type SQLiteRepo struct {
	DB *sql.DB
}

func (r *SQLiteRepo) Create(ctx context.Context, user *model.User) error {
	result, err := r.DB.ExecContext(ctx,
		`insert into user(email, password, authorization_token, is_active, activation_token, password_reset_token) values(?, ?, ?, ?, ?, ?)
		returning id`,
		user.Email.Address, user.PasswordHash, user.AuthorizationToken, user.IsActive, user.ActivationToken, user.PasswordResetToken,
	)
	if err != nil {
		return fmt.Errorf("could not create user: %w", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("could not get inserted id: %w", err)
	}
	user.ID = id
	return nil
}
func (r *SQLiteRepo) Read(ctx context.Context, user *model.User) error {
	row := r.DB.QueryRowContext(ctx,
		"select id, authorization_token, password, email, is_active, activation_token, password_reset_token from user where id = ?", user.ID)
	if row.Err() != nil {
		return fmt.Errorf("could not query user with id=%d: %w", user.ID, row.Err())
	}
	var authorizationToken, email, password, activationToken, passwordResetToken string
	var id int64
	var isActive bool
	err := row.Scan(&id, &authorizationToken, &password, &email, &isActive, &activationToken, &passwordResetToken)
	if err != nil {
		return fmt.Errorf("could not scan row: %w", err)
	}
	user.ID = id
	user.AuthorizationToken = authorizationToken
	user.ActivationToken = activationToken
	user.PasswordResetToken = passwordResetToken
	user.PasswordHash = password
	user.IsActive = isActive
	address, err := mail.ParseAddress(email)
	if err != nil {
		return fmt.Errorf("could not parse email: %w", err)
	}
	user.Email = address
	return nil
}
func (r *SQLiteRepo) Update(ctx context.Context, user *model.User) error {
	_, err := r.DB.ExecContext(ctx,
		`update user set authorization_token = ?, password = ?, email = ?, is_active = ?, activation_token = ?, password_reset_token = ? where id = ?`,
		user.AuthorizationToken, user.PasswordHash, user.Email.Address, user.IsActive, user.ActivationToken, user.PasswordResetToken, user.ID,
	)
	if err != nil {
		return fmt.Errorf("could not update user: %w", err)
	}
	return nil
}
func (r *SQLiteRepo) Delete(ctx context.Context, user *model.User) error {
	_, err := r.DB.ExecContext(ctx, `delete from user where id = ?`, user.ID)
	if err != nil {
		return fmt.Errorf("could not delete user: %w", err)
	}
	return nil
}

func (r *SQLiteRepo) Query(ctx context.Context, filter model.QueryFilter) ([]model.User, error) {
	q := `
	SELECT
		id, authorization_token, password, email, is_active, activation_token, password_reset_token
	FROM user
	`
	buf := bytes.NewBufferString(q)
	args := []any{}
	applyFilter(filter, &args, buf)
	rows, err := r.DB.QueryContext(ctx, buf.String(), args...)
	if err != nil {
		return nil, fmt.Errorf("could not execute query: %w", err)
	}
	users := []model.User{}
	for rows.Next() {
		u := &model.User{}
		var authorizationToken, email, password, activationToken, passwordResetToken string
		var id int64
		var isActive bool
		err := rows.Scan(&id, &authorizationToken, &password, &email, &isActive, &activationToken, &passwordResetToken)
		if err != nil {
			return nil, fmt.Errorf("could not scan row: %w", err)
		}
		u.AuthorizationToken = authorizationToken
		u.PasswordHash = password
		u.IsActive = isActive
		u.ActivationToken = activationToken
		u.PasswordResetToken = passwordResetToken
		address, err := mail.ParseAddress(email)
		if err != nil {
			return nil, fmt.Errorf("could not parse email: %w", err)
		}
		u.ID = id
		u.Email = address
		users = append(users, *u)
	}
	err = rows.Close()
	if err != nil {
		return nil, fmt.Errorf("could not close rows: %w", err)
	}
	return users, nil
}
