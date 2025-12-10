package interests

import "database/sql"

func getInterestsByUserId(userId string, db *sql.DB) ([]Interest, error) {
    query := "SELECT user_id, category_id FROM interests WHERE user_id = $1"
    
    rows, err := db.Query(query, userId)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    var interests []Interest
    for rows.Next() {
        var interest Interest
        if err := rows.Scan(&interest.UserId, &interest.CategoryId); err != nil {
            return nil, err
        }
        interests = append(interests, interest)
    }
    
    if err = rows.Err(); err != nil {
        return nil, err
    }
    
    return interests, nil
}

func createInterest(userId string, categoryId string, db *sql.DB) error {
	
	query := "INSERT INTO interests (user_id, category_id) VALUES ($1, $2)"
	
	_, err := db.Exec(query, userId, categoryId)
	if err != nil {
		return err
	}

	return nil
}

func createInterests(reqs []CreateInterestRequest, db *sql.DB) error {
	for _, req := range reqs {
		if err := createInterest(req.UserId, req.CategoryId, db); err != nil {
			return err
		}
	}
	return nil
}
