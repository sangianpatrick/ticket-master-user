package customer

import "time"

type Customer struct {
	ID           int64      `json:"id"`
	Email        string     `json:"email"`
	Password     string     `json:"password"`
	PasswordSalt string     `json:"password_salt"`
	Name         string     `json:"name"`
	IsVerified   bool       `json:"is_verfied"`
	IsDeleted    bool       `json:"is_deleted"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    *time.Time `json:"updated_at"`
}
