package api

import (
	"context"
	"log"
	"net/http"

	db "github.com/adified/bitespeed/db/sqlc"
	"github.com/gin-gonic/gin"
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
	}
}
