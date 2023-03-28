package controllers

import (
	"net/http"
	"split-rex-backend/entities"
	"split-rex-backend/entities/responses"
	"split-rex-backend/types"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func (h *paymentController) GetUnconfirmedPayment(c echo.Context) error {
	db := h.db
	response := entities.Response[[]responses.UnconfirmedPaymentResponse]{}

	//get user_id from token and group_id from param
	userID := c.Get("id").(uuid.UUID)
	groupID, _ := uuid.Parse(c.QueryParam("group_id"))

	// check if group id present
	group := entities.Group{}
	if err := db.Find(&group, groupID).Error; err != nil {
		response.Message = types.ERROR_INTERNAL_SERVER
		return c.JSON(http.StatusInternalServerError, response)
	}

	// get from payment table where GroupID = groupID and UserID1 = userID and Status = PENDING
	payments := []entities.Payment{}
	conditionPayment := entities.Payment{GroupID: groupID, UserID1: userID, Status: types.STATUS_PAYMENT_PENDING}
	if err := db.Where(conditionPayment).Find(&payments).Error; err != nil {
		response.Message = types.ERROR_INTERNAL_SERVER
		return c.JSON(http.StatusInternalServerError, response)
	}

	// move payments to response data
	unconfirmedTransaction := []responses.UnconfirmedPaymentResponse{}
	for _, payment := range payments {
		user := entities.User{}
		if err := db.Find(&user, payment.UserID2).Error; err != nil {
			response.Message = types.ERROR_INTERNAL_SERVER
			return c.JSON(http.StatusInternalServerError, response)
		}
		unconfirmedTransaction = append(unconfirmedTransaction, responses.UnconfirmedPaymentResponse{
			PaymentID: payment.PaymentID,
			UserID:    payment.UserID2,
			Name:      user.Name,
			Color:     user.Color,
			TotalPaid: payment.TotalPaid,
		})
	}

	response.Message = types.SUCCESS
	response.Data = unconfirmedTransaction
	return c.JSON(http.StatusOK, response)
}
