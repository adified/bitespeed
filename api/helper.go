package api

import (
	"context"
	"log"
	"net/http"
	"strconv"

	db "github.com/adified/bitespeed/db/sqlc"
	"github.com/gin-gonic/gin"
)

type Userreq struct {
	Email       string `josn:"email" binding:"omitempty,email"`
	PhoneNumber int64  `json:"phoneNumber" binding:"omitempty,numeric"`
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
			PhoneNumber: strconv.FormatInt(req.PhoneNumber, 10),
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
			PhoneNumber: strconv.FormatInt(req.PhoneNumber, 10),
		}

		if len(Users) == 0 {
			// fmt.Println("User not found, so create a new entry")
			user, err := server.Querier.CreatePrimaryUser(context.Background(), createUserReq)
			if err != nil {
				log.Fatal(err)
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"error": "User can't be created",
				})
				return
			}

			// {
			// 	"contact":{
			// 		"primaryContatctId": number,
			// 		"emails": string[], // first element being email of primary contact
			// 		"phoneNumbers": string[], // first element being phoneNumber of primary contact
			// 		"secondaryContactIds": number[] // Array of all Contact IDs that are "secondary" to the primary contact
			// 	}
			// }
			type userCreatedRes struct {
				PrimaryContatctId   int64    `json:"primaryContatctId"`
				Emails              []string `json:"emails"`
				PhoneNumbers        []string `json:"phoneNumbers"`
				SecondaryContactIds []int64  `json:"secondaryContactIds"`
			}

			res := userCreatedRes{
				PrimaryContatctId:   user.ID,
				Emails:              []string{user.Email},
				PhoneNumbers:        []string{user.PhoneNumber},
				SecondaryContactIds: make([]int64, 0),
			}
			c.JSON(http.StatusCreated, gin.H{
				"contact": res,
			})
		} else {

		}

		c.Next()
		status := c.Writer.Status()
		log.Println(status)
		return
	}
}
