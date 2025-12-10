package users

import "database/sql"

func createUserService(email string, externalId string, db *sql.DB, mqCallback MqCallback) (User, error) {
	return createUser(email, externalId, db, mqCallback)
}

func getUserByExternalIdService(externalId string, db *sql.DB) (User, error) {
	return getUserByExternalId(externalId, db)
}
