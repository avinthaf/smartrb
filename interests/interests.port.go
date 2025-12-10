package interests

import "database/sql"

type CreateInterestRequest struct {
	UserId     string `json:"user_id"`
	CategoryId string `json:"category_id"`
}

func GetInterestsByUserId(userId string, db *sql.DB) ([]Interest, error) {
	return getInterestsByUserIdService(userId, db)
}

func CreateInterest(req CreateInterestRequest, db *sql.DB) error {
	return createInterestService(req.UserId, req.CategoryId, db)
}

func CreateInterests(reqs []CreateInterestRequest, db *sql.DB) error {
	return createInterestsService(reqs, db)
}