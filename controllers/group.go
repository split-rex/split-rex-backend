package controllers

import (
	"net/http"
	"split-rex-backend/configs/database"
	"split-rex-backend/entities"
	"split-rex-backend/entities/requests"
	"split-rex-backend/types"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func UserCreateGroupController(c echo.Context) error {
	db := database.DB.GetConnection()
	response := entities.Response[string]{}

	request := requests.UserCreateGroupRequest{}
	if err := c.Bind(&request); err != nil {
		response.Message = types.ERROR_BAD_REQUEST
		return c.JSON(http.StatusBadRequest, response)
	}

	if err := db.Create(
		&entities.Group{
			GroupID:   uuid.New(),
			Name:      request.Name,
			StartDate: request.StartDate,
			EndDate:   request.EndDate,
		}).Error; err != nil {
		response.Message = types.ERROR_INTERNAL_SERVER
		return c.JSON(http.StatusInternalServerError, response)
	}

	response.Message = types.SUCCESS
	return c.JSON(http.StatusAccepted, response)
}

func EditGroupInfoController(c echo.Context) error {
	db := database.DB.GetConnection()
	response := entities.Response[string]{}

	request := requests.EditGroupInfoRequest{}
	if err := c.Bind(&request); err != nil {
		response.Message = types.ERROR_BAD_REQUEST
		return c.JSON(http.StatusBadRequest, response)
	}

	group := entities.Group{}
	condition := entities.Group{GroupID: request.GroupID}
	if err := db.Where(&condition).Find(&group).Error; err != nil {
		response.Message = types.ERROR_INTERNAL_SERVER
		return c.JSON(http.StatusInternalServerError, response)
	}

	group.Name = request.Name
	group.StartDate = request.StartDate
	group.EndDate = request.EndDate

	if err := db.Save(&group).Error; err != nil {
		response.Message = types.ERROR_INTERNAL_SERVER
		return c.JSON(http.StatusInternalServerError, response)
	}

	response.Message = types.SUCCESS
	return c.JSON(http.StatusAccepted, response)
}

func UserGroupsController(c echo.Context) error {
	// db := database.DB.GetConnection()
	response := entities.Response[string]{}
	return c.JSON(http.StatusAccepted, response)
}

func GroupDetailController(c echo.Context) error {
	// db := database.DB.GetConnection()
	// config := configs.Config.GetMetadata()
	response := entities.Response[string]{}
	return c.JSON(http.StatusAccepted, response)
}

func GroupTransactionsController(c echo.Context) error {
	// db := database.DB.GetConnection()
	// config := configs.Config.GetMetadata()
	response := entities.Response[string]{}
	return c.JSON(http.StatusAccepted, response)
}
