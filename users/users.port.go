package users

import "database/sql"

type UpdateUserRequest struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

type MqCallback func(topic string, routingKey string, message string)

func CreateUser(email string, externalId string, db *sql.DB, mqCallback MqCallback) (User, error) {
	return createUserService(email, externalId, db, mqCallback)
}

func UpdateUser(userId string, firstName string, lastName string, db *sql.DB) (User, error) {
	return updateUserService(userId, firstName, lastName, db)
}

func GetUserByExternalId(externalId string, db *sql.DB) (User, error) {
	return getUserByExternalIdService(externalId, db)
}
