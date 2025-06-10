package api

import (
	"context"
	"log"
	"net/http"

	db "github.com/adified/bitespeed/db/sqlc"
	"github.com/gin-gonic/gin"
)

type Userreq struct {
	Email       string `josn:"email" binding:"omitempty,email"`
	PhoneNumber string `json:"phoneNumber" binding:"omitempty,numeric"`
}

func CheckifExists(server *Server) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req Userreq
		if err := c.ShouldBindJSON(&req); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": "invalid input",
			})
			return
		}

		arg := db.FindUsersByEmailOrPhoneParams{
			Email:       req.Email,
			PhoneNumber: req.PhoneNumber,
		}

		Users, err := server.Querier.FindUsersByEmailOrPhone(context.Background(), arg)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "Status Internal Server Error",
			})
			return
		}

		createUserReq := db.CreatePrimaryUserParams{
			Email:       req.Email,
			PhoneNumber: req.PhoneNumber,
		}
		if len(Users) == 0 {
			// fmt.Println("User not found, so create a new entry")
			user, err := server.Querier.CreatePrimaryUser(context.Background(), createUserReq)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"error": "User can't be created",
				})
				return
			}
			c.JSON(http.StatusCreated, gin.H{
				"user": user,
			})
		}

		c.Next()
		status := c.Writer.Status()
		log.Println(status)
		return
	}
}
