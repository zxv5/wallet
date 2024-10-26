package model

import (
	"time"
)

const TableNameUser = "user"

// User mapped from table <user>
type User struct {
	ID        int64     `json:"id"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Gender    int64     `json:"gender"`
	Status    int64     `json:"status"` // 0=active 1=inactive 2=suspens
	Deleted   int64     `json:"deleted"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedAt time.Time `json:"created_at"`
}

// TableName User's table name
func (*User) TableName() string {
	return TableNameUser
}
