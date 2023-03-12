package controllers

import (
	"net/http"
	"split-rex-backend/entities"
	"split-rex-backend/entities/responses"
	"split-rex-backend/types"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func (con *authController) ProfileController(c echo.Context) error {
	db := con.db
	response := entities.Response[responses.ProfileResponse]{}

	user_id := c.Get("id").(uuid.UUID)

	user := entities.User{}
	condition := entities.User{ID: user_id}
	if err := db.Where(&condition).Select("id", "username", "name").Find(&user).Error; err != nil {
		response.Message = types.ERROR_INTERNAL_SERVER
		return c.JSON(http.StatusInternalServerError, response)
	}
	response.Message = types.SUCCESS

	response.Data.User_id = user.ID.String()
	response.Data.Username = user.Username
	response.Data.Fullname = user.Name

	return c.JSON(http.StatusOK, response)
}
