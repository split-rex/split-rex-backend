package controllers

import (
	"net/http"
	"split-rex-backend/configs/database"
	"split-rex-backend/entities"
	"split-rex-backend/entities/responses"
	"split-rex-backend/types"

	"github.com/labstack/echo/v4"
)

func SearchUser(c echo.Context) error {
	username := c.QueryParam("username")
	db := database.DB.GetConnection()
	response := entities.Response[responses.FriendResponse]{}

	//check if username exist in user table
	user := []entities.User{}
	conditionUser := entities.User{Username: username}
	if err := db.Where(&conditionUser).Find(&user).Error; err != nil {
		response.Message = types.ERROR_INTERNAL_SERVER
		return c.JSON(http.StatusInternalServerError, response)
	}

	if len(user) > 0 {
		response.Message = types.SUCCESS
		friend := responses.FriendResponse{}
		friend.User_id = user[0].ID.String()
		friend.Username = user[0].Username
		friend.Fullname = user[0].Name
		response.Data = friend
	} else {
		response.Message = types.DATA_NOT_FOUND
	}
	return c.JSON(http.StatusOK, response)
}
