package controllers

import (
	"net/http"
	"split-rex-backend/entities"
	"split-rex-backend/entities/responses"
	"split-rex-backend/types"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func (con *friendController) UserFriendList(c echo.Context) error {
	db := con.db

	// get user id from context, cast to uuid
	user_id := c.Get("id").(uuid.UUID)

	// returning array of friendResponse struct
	response := entities.Response[[]responses.ProfileResponse]{}

	//check if user_id exist in friend table
	userInFriend := entities.Friend{}
	conditionSearchUser := entities.Friend{ID: user_id}
	if err := db.Where(&conditionSearchUser).Find(&userInFriend).Error; err != nil {
		response.Message = types.ERROR_INTERNAL_SERVER
		return c.JSON(http.StatusInternalServerError, response)
	}

	//check if userFriend empty
	if userInFriend.ID == uuid.Nil {
		response.Message = types.DATA_NOT_FOUND
		response.Data = []responses.ProfileResponse{}
		return c.JSON(http.StatusOK, response)
	}

	// get username and full name where friend_id (Friend table) exist for user
	friends := []responses.ProfileResponse{}
	for _, id := range userInFriend.Friend_id {
		user := entities.User{}
		friend := responses.ProfileResponse{}
		condition := entities.User{ID: id}
		// get id, username, and name from user table
		if err := db.Where(&condition).Select("id", "username", "name").Find(&user).Error; err != nil {
			response.Message = types.ERROR_INTERNAL_SERVER
			return c.JSON(http.StatusInternalServerError, response)
		}
		friend.User_id = user.ID.String()
		friend.Username = user.Username
		friend.Fullname = user.Name
		friends = append(friends, friend)
	}

	response.Data = friends

	// if no friends, data not found but status success
	if len(friends) <= 0 {
		response.Message = types.DATA_NOT_FOUND
		return c.JSON(http.StatusOK, response)
	}

	response.Message = types.SUCCESS
	return c.JSON(http.StatusOK, response)

}
