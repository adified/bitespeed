package api

import (
	"context"
	"log"
	"net/http"
	"sort"

	db "github.com/adified/bitespeed/db/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
)

type Userreq struct {
	Email       string `json:"email" binding:"omitempty,email"`
	PhoneNumber string `json:"phoneNumber" binding:"omitempty,numeric"`
}

type userCreatedRes struct {
	PrimaryContatctId   int64    `json:"primaryContatctId"`
	Emails              []string `json:"emails"`
	PhoneNumbers        []string `json:"phoneNumbers"`
	SecondaryContactIds []int64  `json:"secondaryContactIds"`
}

// CheckifExists is the main handler function
func (server *Server) CheckifExists() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req Userreq
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
			return
		}

		// request is invalid if both fields are missing or are zero
		if req.Email == "" && req.PhoneNumber == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Either email or phoneNumber must be provided"})
			return
		}

		arg := db.FindUsersByEmailOrPhoneParams{}
		if req.Email != "" {
			arg.Email = req.Email
		}
		if req.PhoneNumber != "" {
			arg.PhoneNumber = req.PhoneNumber
		}

		// use the server's default querier for the initial find operation
		existingUsers, err := server.Querier.FindUsersByEmailOrPhone(c.Request.Context(), arg)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database query failed"})
			return
		}

		// case 1: no user exists so create a new primary contact
		if len(existingUsers) == 0 {
			createUserReq := db.CreatePrimaryUserParams{
				Email:       req.Email,
				PhoneNumber: req.PhoneNumber,
			}
			user, err := server.Querier.CreatePrimaryUser(context.Background(), createUserReq)
			if err != nil {
				log.Printf("Error creating primary user: %v", err)
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "User can't be created"})
				return
			}

			// buildResponse helper to format the response
			res := buildResponse([]db.User{user})
			c.JSON(http.StatusOK, gin.H{"contact": res})
			return
		}

		// case 2: one or more users already exist
		handleExistingContacts(c, server, req, existingUsers)
	}
}

// handleExistingContacts handels cases where a related contact already exists
func handleExistingContacts(c *gin.Context, server *Server, req Userreq, existingUsers []db.User) {
	// start a transaction
	tx, err := server.pool.Begin(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to start transaction"})
		return
	}
	defer tx.Rollback(c.Request.Context()) // Ensure rollback on any error

	// a new querier that operates within this transaction
	qtx := server.Querier.WithTx(tx)

	// main login begins

	// find all unique primary contacts involved
	primaryUsersMap := make(map[int64]db.User)
	for _, user := range existingUsers {
		if user.LinkPrecedence == "primary" {
			primaryUsersMap[user.ID] = user
		} else if user.LinkedID.Valid {
			// If we found a secondary, find its primary to ensure we're processing the whole chain
			primary, err := qtx.GetUserByID(c.Request.Context(), user.LinkedID.Int64)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch linked primary user"})
				return
			}
			primaryUsersMap[primary.ID] = primary
		}
	}

	// determine the oldest primary user. This will be the ultimate primary
	var primaries []db.User
	for _, p := range primaryUsersMap {
		primaries = append(primaries, p)
	}
	sort.Slice(primaries, func(i, j int) bool {
		return primaries[i].CreatedAt.Time.Before(primaries[j].CreatedAt.Time)
	})

	ultimatePrimary := primaries[0]

	// if there are other primary users then we must modify them to become secondary to the ultimate primary
	for i := 1; i < len(primaries); i++ {
		userToUpdate := primaries[i]
		err := qtx.UpdateUserToSecondary(c.Request.Context(), db.UpdateUserToSecondaryParams{
			LinkedID: pgtype.Int8{
				Int64: ultimatePrimary.ID,
				Valid: true,
			},
			ID: userToUpdate.ID,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to link primary contacts"})
			return
		}
	}

	// get all contacts to check for new info
	allConnectedUsers, err := qtx.GetUsersByPrimaryID(c.Request.Context(), ultimatePrimary.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get all associated contacts"})
		return
	}

	// check if the request contains new info
	hasNewInfo := false
	if req.Email != "" && !containsDuplicateEmail(allConnectedUsers, req.Email) {
		hasNewInfo = true
	}
	if req.PhoneNumber != "" && !containsDuplicatePhone(allConnectedUsers, req.PhoneNumber) {
		hasNewInfo = true
	}

	// if there's new info then  create a new secondary entry in db
	if hasNewInfo {
		newSecondary, err := qtx.CreateSecondaryUser(c.Request.Context(), db.CreateSecondaryUserParams{
			Email:       req.Email,
			PhoneNumber: req.PhoneNumber,
			LinkedID: pgtype.Int8{
				Int64: ultimatePrimary.ID,
				Valid: true,
			},
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create secondary user entry"})
			return
		}
		// add the new contact to our list for the final response
		allConnectedUsers = append(allConnectedUsers, newSecondary)
	}

	// if we've reached here all db actions must have been successful so now we commit the transaction
	if err := tx.Commit(c.Request.Context()); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to commit transction"})
		return
	}

	// return the final response
	finalResponse := buildResponse(allConnectedUsers)
	c.JSON(http.StatusOK, gin.H{"contact": finalResponse})
}
