package interests

import "database/sql"

func getInterestsByUserIdService(userId string, db *sql.DB) ([]Interest, error) {
	return getInterestsByUserId(userId, db)
}

func createInterestService(userId string, categoryId string, db *sql.DB) error {
	return createInterest(userId, categoryId, db)
}

func createInterestsService(reqs []CreateInterestRequest, db *sql.DB) error {
	return createInterests(reqs, db)
}
