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

func SearchUser(c echo.Context) error {
	username := c.QueryParam("username")
	db := database.DB.GetConnection()
	response := entities.Response[responses.ProfileResponse]{}

	//check if username exist in user table
	user := []entities.User{}
	conditionUser := entities.User{Username: username}
	if err := db.Where(&conditionUser).Find(&user).Error; err != nil {
		response.Message = types.ERROR_INTERNAL_SERVER
		return c.JSON(http.StatusInternalServerError, response)
	}

	if len(user) > 0 {
		response.Message = types.SUCCESS
		friend := responses.ProfileResponse{}
		friend.User_id = user[0].ID.String()
		friend.Username = user[0].Username
		friend.Fullname = user[0].Name
		response.Data = friend
	} else {
		response.Message = types.DATA_NOT_FOUND
	}
	return c.JSON(http.StatusOK, response)
}

func SearchUserToAdd(c echo.Context) error {
	username := c.QueryParam("username")
	db := database.DB.GetConnection()
	response := entities.Response[responses.ProfileResponse]{}

	//check if username exist in user table
	user := entities.User{}
	conditionUser := entities.User{Username: username}
	if err := db.Where(&conditionUser).Find(&user).Error; err != nil {
		response.Message = types.ERROR_INTERNAL_SERVER
		return c.JSON(http.StatusInternalServerError, response)
	}
	if user.Username == "" {
		response.Message = types.ERROR_USER_NOT_FOUND
		return c.JSON(http.StatusBadRequest, response)
	}

	// check if userid same with user_id in jwt token
	currentUserId := c.Get("id").(uuid.UUID)
	if user.ID == currentUserId {
		response.Message = types.ERROR_CANNOT_ADD_SELF
		return c.JSON(http.StatusBadRequest, response)
	}

	// check if user_id in current user friends
	currentUser := entities.Friend{}
	if err := db.Find(&currentUser, currentUserId).Error; err != nil {
		response.Message = types.ERROR_INTERNAL_SERVER
		return c.JSON(http.StatusInternalServerError, response)
	}

	for _, friend_id := range currentUser.Friend_id {
		if friend_id == user.ID {
			response.Message = types.ERROR_ALREADY_FRIEND
			return c.JSON(http.StatusBadRequest, response)
		}
	}

	response.Message = types.SUCCESS
	response.Data = responses.ProfileResponse{
		User_id:  user.ID.String(),
		Username: user.Username,
		Fullname: user.Name,
	}
	return c.JSON(http.StatusOK, response)
}
