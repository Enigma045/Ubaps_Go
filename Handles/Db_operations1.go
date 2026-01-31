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
