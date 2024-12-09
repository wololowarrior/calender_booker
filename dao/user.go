package dao

import (
	"database/sql"
	"errors"
	"fmt"

	"calendly_adventures/dao/query"
	"calendly_adventures/db"
	"calendly_adventures/models"
)

func Get(userID int) (*models.User, error) {
	var user models.User
	err := db.DB.QueryRow(query.GetUser, userID).Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user with ID %d not found", userID)
		}
		return nil, fmt.Errorf("failed to retrieve user: %w", err)
	}

	return &user, nil
}

func Create(user *models.User) (*models.User, error) {
	err := db.DB.QueryRow(query.InsertUser, user.Name, user.Email).Scan(&user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to insert user: %w", err)
	}
	return user, nil
}
