package domain

import "time"

type Customer struct {
	ID              int64
	Email           string
	PasswordHash    string
	Firstname       string
	Lastname        string
	Phone           string
	IsActive        bool
	EmailVerifiedAt *time.Time
	LastLoginAt     *time.Time
	CreatedAt       time.Time
	UpdatedAt       time.Time
}
