package dao

import (
	"context"
	"database/sql"
	"wallet/internal/model"
	"wallet/pkg/e"
	"wallet/pkg/logger"
)

type User struct {
	dao *Dao
}

func NewUser(dao *Dao) *User {
	return &User{dao: dao}
}

func (d *User) FindOneByID(ctx context.Context, userID int64) (*model.User, e.Codes) {
	var user model.User

	query := `SELECT "id", "email", "password", "first_name", "last_name", "gender", "status", "created_at" FROM "user" WHERE "id" = $1;`
	row := d.dao.WithContext(ctx).QueryRow(query, userID)

	err := row.Scan(&user.ID, &user.Email, &user.Password, &user.FirstName, &user.LastName, &user.Gender, &user.Status, &user.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, e.NotFound
		}
		logger.Errorf("failed to find user: %w", err)
		return nil, e.SQLErr
	}

	return &user, nil
}

func (d *User) FindOneByEmail(ctx context.Context, email string) (*model.User, e.Codes) {
	var user model.User

	query := `SELECT "id", "email", "password", "first_name", "last_name", "gender", "status", "created_at" FROM "user" WHERE "email" = $1;`
	row := d.dao.WithContext(ctx).QueryRow(query, email)

	err := row.Scan(&user.ID, &user.Email, &user.Password, &user.FirstName, &user.LastName, &user.Gender, &user.Status, &user.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, e.NotFound
		}
		logger.Errorf("failed to find user: %w", err)
		return nil, e.SQLErr
	}

	return &user, nil
}
