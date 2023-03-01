package controllers

import (
	"net/http"
	"split-rex-backend/configs/database"
	"split-rex-backend/entities"
	"split-rex-backend/entities/responses"
	"split-rex-backend/types"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func FriendRequestSent(c echo.Context) error {
	user_id := c.Get("id").(uuid.UUID)
	db := database.DB.GetConnection()
	response := entities.Response[[]responses.FriendResponse]{}

	//check if user_id exist in friend table
	userFriend := entities.Friend{}
	userExist := true
	conditionFriend := entities.Friend{ID: user_id}
	if err := db.Where(&conditionFriend).Find(&userFriend).Error; err != nil {
		response.Message = types.ERROR_INTERNAL_SERVER
		return c.JSON(http.StatusInternalServerError, response)
	}
	//check if userFriend empty
	if userFriend.ID == uuid.Nil {
		userExist = false
	}

	if userExist {
		// get username and full name where user_id in Req_sent
		users := []responses.FriendResponse{}
		for _, id := range userFriend.Req_sent {
			user := entities.User{}
			friend := responses.FriendResponse{}
			condition := entities.User{ID: id}
			if err := db.Where(&condition).Select("id", "username", "name").Find(&user).Error; err != nil {
				response.Message = types.ERROR_INTERNAL_SERVER
				return c.JSON(http.StatusInternalServerError, response)
			}
			friend.User_id = user.ID.String()
			friend.Username = user.Username
			friend.Fullname = user.Name
			users = append(users, friend)
		}
		response.Message = types.SUCCESS
		if len(users) > 0 {
			response.Data = users
		} else {
			response.Message = types.DATA_NOT_FOUND
		}
		return c.JSON(http.StatusOK, response)

	} else {
		response.Message = types.DATA_NOT_FOUND
		response.Data = []responses.FriendResponse{}
		return c.JSON(http.StatusOK, response)
	}
}
