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

func getUserByExternalId(externalId string, db *sql.DB) (User, error) {
	user := User{}
	
	query := "SELECT id, email, external_id FROM users WHERE external_id = $1"
	
	err := db.QueryRow(query, externalId).Scan(&user.Id, &user.Email, &user.ExternalId)
	if err != nil {
		return User{}, fmt.Errorf("error getting user by external id: %v", err)
	}

	return User{
		Id:          user.Id,
		Email:       user.Email,		
		ExternalId:  user.ExternalId,
	}, nil
}


