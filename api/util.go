package api

import db "github.com/adified/bitespeed/db/sqlc"

// Helper functions to check for duplicates
func containsDuplicateEmail(users []db.User, email string) bool {
	for _, u := range users {
		if u.Email == email {
			return true
		}
	}
	return false
}

func containsDuplicatePhone(users []db.User, phone string) bool {
	for _, u := range users {
		if u.PhoneNumber == phone {
			return true
		}
	}
	return false
}
