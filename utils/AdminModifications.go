package utils

import (
	"context"

	"github.com/jackc/pgx/v5"
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
	limit, offset int,
) ([]UserDetails, int, error) {

	countQuery := `SELECT COUNT(*) FROM users WHERE user_id <> $1`
	var total int
	err := pool.QueryRow(ctx, countQuery, userID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

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
        ORDER BY created_at DESC
        LIMIT $2 OFFSET $3;
     `

	rows, err := pool.Query(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, 0, err
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
			&user.Role,
			&user.Verified,
		)
		if err != nil {
			return nil, 0, err
		}
		userdetails = append(userdetails, user)
	}
	return userdetails, total, nil
}

func DeleteUser(tx pgx.Tx,
	            ctx context.Context,
				 userid int64) (error) {

	query := `DELETE FROM
	users WHERE
	user_id = $1`

	_,err := tx.Exec(ctx,query,userid)
	if err != nil {
    return err
	}

	return nil
}