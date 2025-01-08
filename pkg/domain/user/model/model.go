package model

import (
	"fmt"
	"math/rand"
	"net/mail"

	"golang.org/x/crypto/bcrypt"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

type User struct {
	ID                 int64
	Email              *mail.Address
	AuthorizationToken string
	Password           string
	PasswordHash       string
	IsActive           bool
	ActivationToken    string
	PasswordResetToken string
}

func (u *User) MatchPassword(password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
	if err != nil {
		return fmt.Errorf("password comparison failed: %w", err)
	}
	return nil
}

func (u *User) SetPasswordHash() error {
	if u.Password == "" {
		return fmt.Errorf("password has to be set")
	}
	pwdHash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("could not hash password: %w", err)
	}
	u.PasswordHash = string(pwdHash)
	fmt.Println(u.PasswordHash)
	return nil
}

func Parse(email string, options ...func(*User) error) (*User, error) {
	user := &User{}
	address, err := mail.ParseAddress(email)
	if err != nil {
		return nil, fmt.Errorf("could not parse email: %w", err)
	}
	user.Email = address
	for _, o := range options {
		err := o(user)
		if err != nil {
			return nil, fmt.Errorf("could not construct user: %w", err)
		}
	}
	return user, nil
}

func WithPassword(password string) func(*User) error {
	return func(u *User) error {
		return SetPassword(password, u)
	}
}

func WithAuthorizationToken() func(*User) error {
	return func(u *User) error {
		u.AuthorizationToken = generateRandomString(50)
		return nil
	}
}
func WithActivationToken() func(*User) error {
	return func(u *User) error {
		u.ActivationToken = generateRandomString(50)
		return nil
	}
}

func generateRandomString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func SetPassword(password string, u *User) error {
	if len(password) < 8 {
		return fmt.Errorf("password can't be empty less then 8 characters")
	}
	u.Password = password
	err := u.SetPasswordHash()
	if err != nil {
		return fmt.Errorf("couldn't set hash: %w", err)
	}
	return nil
}

func WithID(id int64) func(*User) error {
	return func(u *User) error {
		u.ID = id
		return nil
	}
}

func (u *User) GeneratePasswordResetToken() string {
	u.PasswordResetToken = generateRandomString(64)
	return u.PasswordResetToken
}
