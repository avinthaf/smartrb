package users

import (
	"database/sql"
	"fmt"
)

func createUser(email string, externalId string, db *sql.DB, mqCallback MqCallback) (User, error) {
	
	user := User{}
	
	query := "INSERT INTO users (email, external_id) VALUES ($1, $2) RETURNING id"
	
	err := db.QueryRow(query, email, externalId).Scan(&user.Id)
	if err != nil {
		return User{}, fmt.Errorf("error creating user: %v", err)
	}

	return User{
		Id:          user.Id,
		Email:       email,		
		ExternalId:  externalId,
	}, nil
}

func updateUser(userId string, firstName string, lastName string, db *sql.DB) (User, error) {
	user := User{}
	
	query := "UPDATE users SET first_name = $1, last_name = $2, updated_at = NOW() WHERE id = $3 RETURNING id, email, external_id, created_at, updated_at"
	
	err := db.QueryRow(query, firstName, lastName, userId).Scan(&user.Id, &user.Email, &user.ExternalId, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return User{}, fmt.Errorf("user not found")
		}
		return User{}, fmt.Errorf("error updating user: %v", err)
	}

	return user, nil
}

func getUserByExternalId(externalId string, db *sql.DB) (User, error) {
	user := User{}
	
	query := "SELECT id, email, external_id FROM users WHERE external_id = $1"
	
	err := db.QueryRow(query, externalId).Scan(&user.Id, &user.Email, &user.ExternalId)
	if err != nil {
		if err == sql.ErrNoRows {
			return User{}, fmt.Errorf("user not found")
		}
		return User{}, fmt.Errorf("error getting user by external id: %v", err)
	}

	return User{
		Id:          user.Id,
		Email:       user.Email,		
		ExternalId:  user.ExternalId,
	}, nil
}


