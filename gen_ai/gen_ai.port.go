package gen_ai

import "database/sql"

func CreateAIContent(prompt string, db *sql.DB) (string, error) {
	return createAIContentService(prompt, db)
}