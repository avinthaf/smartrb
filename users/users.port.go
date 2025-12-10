package users

import "database/sql"

type MqCallback func(topic string, routingKey string, message string)

func CreateUser(email string, externalId string, db *sql.DB, mqCallback MqCallback) (User, error) {
	return createUserService(email, externalId, db, mqCallback)
}

func GetUserByExternalId(externalId string, db *sql.DB) (User, error) {
	return getUserByExternalIdService(externalId, db)
}
