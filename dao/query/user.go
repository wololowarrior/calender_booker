package query

var GetUser = `
		SELECT id, name, email, created_at
		FROM users
		WHERE id = $1
	`

var InsertUser = `INSERT INTO users (name, email) VALUES ($1, $2) returning id`
