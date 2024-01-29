package customer

import "time"

type Customer struct {
	ID         int64      `json:"id"`
	Email      string     `json:"email"`
	Password   string     `json:"password"`
	Name       string     `json:"name"`
	IsVerified bool       `json:"isVerfied"`
	IsDeleted  bool       `json:"isDeleted"`
	CreatedAt  time.Time  `json:"createdAt"`
	UpdatedAt  *time.Time `json:"updatedAt"`
}
