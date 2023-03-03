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

func UserCreateTransactionController(c echo.Context) error {
	db := database.DB.GetConnection()
	response := entities.Response[string]{}

	request := requests.UserCreateTransactionRequest{}
	if err := c.Bind(&request); err != nil {
		response.Message = err.Error()
		return c.JSON(http.StatusBadRequest, response)
	}

	// check if group id present
	group := entities.Group{}
	if err := db.Find(&group, request.GroupID).Error; err != nil {
		response.Message = err.Error()
		return c.JSON(http.StatusBadRequest, response)
	}

	// check if billowner present in user
	user := entities.User{}
	if err := db.Find(&user, request.BillOwner).Error; err != nil {
		response.Message = types.ERROR_BAD_REQUEST
		return c.JSON(http.StatusBadRequest, response)
	}

	// TODO: check if item exist in item table

	transaction := entities.Transaction{
		TransactionID: uuid.New(),
		Name:          request.Name,
		Description:   request.Description,
		GroupID:       request.GroupID,
		Date:          request.Date,
		Subtotal:      request.Subtotal,
		Tax:           request.Tax,
		Service:       request.Service,
		Total:         request.Total,
		BillOwner:     request.BillOwner,
		Items:         request.Items,
	}

	if err := db.Create(&transaction).Error; err != nil {
		response.Message = err.Error()
		return c.JSON(http.StatusInternalServerError, response)
	}

	response.Message = types.SUCCESS
	return c.JSON(http.StatusOK, response)
}
