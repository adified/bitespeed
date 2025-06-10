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

// buildResponse constructs the final JSON from a list of contacts.
func buildResponse(contacts []db.User) userCreatedRes {
	if len(contacts) == 0 {
		return userCreatedRes{}
	}

	var primaryContact db.User
	secondaryIDs := make([]int64, 0)
	emailSet := make(map[string]struct{})
	phoneSet := make(map[string]struct{})

	for _, c := range contacts {
		if c.LinkPrecedence == "primary" {
			primaryContact = c
		} else {
			secondaryIDs = append(secondaryIDs, c.ID)
		}
		if c.Email != "" {
			emailSet[c.Email] = struct{}{}
		}
		if c.PhoneNumber != "" {
			phoneSet[c.PhoneNumber] = struct{}{}
		}
	}

	// Ensure primary contact's info is first in the list.
	emails := []string{}
	if primaryContact.Email != "" {
		emails = append(emails, primaryContact.Email)
		delete(emailSet, primaryContact.Email)
	}
	for email := range emailSet {
		emails = append(emails, email)
	}

	phones := []string{}
	if primaryContact.PhoneNumber != "" {
		phones = append(phones, primaryContact.PhoneNumber)
		delete(phoneSet, primaryContact.PhoneNumber)
	}
	for phone := range phoneSet {
		phones = append(phones, phone)
	}

	return userCreatedRes{
		PrimaryContatctId:   primaryContact.ID,
		Emails:              emails,
		PhoneNumbers:        phones,
		SecondaryContactIds: secondaryIDs,
	}
}
