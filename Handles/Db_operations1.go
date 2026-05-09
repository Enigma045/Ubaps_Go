package Handles

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	return string(bytes), err
}

func CheckPassword(hash, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}

func CreateUser(
	name, surname, email, phone, password, userType string, tx pgx.Tx, is_verified bool,
) (int64,error) {
	password_hash, err := HashPassword(password)
	if err != nil {
		fmt.Errorf("Password hashing failed")
	}
	// Default role
	if userType == "" {
		userType = "student"
	}

	// Optional: validate in Go (extra safety)
	allowed := map[string]bool{
		"admin":           true,
		"student":         true,
		"dean_of_student": true,
		"dean_of_science": true,
		"dean_of_facult": true,
		"registrar":       true,
		"finance_office":  true,
	}

	if !allowed[userType] {
		return 0,fmt.Errorf("invalid user type")
	}

	var userID int64
	query := `
    INSERT INTO users
        (name, surname, email, phone, password_hash, user_type, is_active, is_verified)
    VALUES ($1,$2,$3,$4,$5,$6,true,$7)
    RETURNING user_id
`

	err = tx.QueryRow(
		context.Background(),
		query,
		name, surname, email, phone, password_hash, userType,is_verified,
	).Scan(&userID)

	if err != nil {
		if strings.Contains(err.Error(), "users_email_unique") {
			return 0,fmt.Errorf("email already exists")
		}
		return 0,err
	}
	fmt.Println("success2")
	return userID,nil
}

func UpdateUserProfile(ctx context.Context, tx pgx.Tx, userID int64, name, surname, email, phone string) error {
	query := `
		UPDATE users 
		SET name = $1, surname = $2, email = $3, phone = $4, updated_at = NOW()
		WHERE user_id = $5
	`
	_, err := tx.Exec(ctx, query, name, surname, email, phone, userID)
	return err
}

func UpdateUserPassword(ctx context.Context, tx pgx.Tx, userID int64, newPassword string) error {
	hash, err := HashPassword(newPassword)
	if err != nil {
		return err
	}
	query := `UPDATE users SET password_hash = $1, updated_at = NOW() WHERE user_id = $2`
	_, err = tx.Exec(ctx, query, hash, userID)
	return err
}

func AdminUpdateUser(ctx context.Context, tx pgx.Tx, userID int64, name, surname, email, phone, status string) error {
	isActive := status == "true"
	query := `
		UPDATE users 
		SET name = $1, surname = $2, email = $3, phone = $4, is_active = $5, updated_at = NOW()
		WHERE user_id = $6
	`
	_, err := tx.Exec(ctx, query, name, surname, email, phone, isActive, userID)
	return err
}
