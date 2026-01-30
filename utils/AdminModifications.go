package utils

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserDetails struct {
	First string `json:"first"`
	Last  string `json:"last"`
	Email string `json:"email"`
	Phone string `json:"phone"`
	//Password string `json:"password"`
	Role string `json:"role"`
	Verified bool `json:"verified"`
}

func ReciveDetails(
	pool *pgxpool.Pool,
	ctx context.Context,
	userID int64,
) ([]UserDetails, error) {

	query := `
		SELECT
        name,
        surname,
        email,
		phone,
		user_type,
		is_verified
		FROM users
		WHERE user_id <> $1
        ORDER BY created_at DESC;
     `

	rows, err := pool.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var userdetails []UserDetails

	for rows.Next() {
		var user UserDetails

		err := rows.Scan(
			&user.First,
			&user.Last,
			&user.Email,
			&user.Phone,
			//&user.Password,
			&user.Role,
			&user.Verified,

		)
		if err != nil {
			return nil, err
		}

		userdetails = append(userdetails, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return userdetails, nil
}