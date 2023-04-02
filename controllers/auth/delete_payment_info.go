package controllers

import (
	"net/http"
	"split-rex-backend/entities"
	"split-rex-backend/entities/requests"
	"split-rex-backend/types"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func (con *authController) DeletePaymentInfo(c echo.Context) error {
	db := con.db
	response := entities.Response[string]{}

	deletePaymentInfo := requests.PaymentInfo{}
	if err := c.Bind(&deletePaymentInfo); err != nil {
		response.Message = types.ERROR_INTERNAL_SERVER
		return c.JSON(http.StatusBadRequest, response)
	}

	// get id from context, cast to uuid
	user_id := c.Get("id").(uuid.UUID)

	// search the id of the user
	user := entities.User{}
	if err := db.Find(&user, user_id).Error; err != nil {
		response.Message = types.ERROR_INTERNAL_SERVER
		return c.JSON(http.StatusInternalServerError, response)
	}

	// error if user not found
	if user.ID == uuid.Nil {
		response.Message = types.ERROR_BAD_REQUEST
		return c.JSON(http.StatusInternalServerError, response)
	}

	// get all the user's payment info
	userPaymentInfo := types.PaymentInfo{}
	if len(user.PaymentInfo) > 0 {
		userPaymentInfo = user.PaymentInfo
	}

	// check if payment method available
	payment, ok := userPaymentInfo[deletePaymentInfo.Payment_method]
	if !ok {
		response.Message = types.ERROR_BAD_REQUEST
		return c.JSON(http.StatusInternalServerError, response)
	}

	// check if account number available
	_, ok = payment[int(deletePaymentInfo.Account_number)]
	if !ok {
		response.Message = types.ERROR_BAD_REQUEST
		return c.JSON(http.StatusInternalServerError, response)
	}

	// delete account number and account name
	delete(userPaymentInfo[deletePaymentInfo.Payment_method], int(deletePaymentInfo.Account_number))

	// if account number in payment method = 0, delete payment method
	if len(userPaymentInfo[deletePaymentInfo.Payment_method]) < 1 {
		delete(userPaymentInfo, deletePaymentInfo.Payment_method)
	}

	// update user
	if err := db.Model(&user).Updates(entities.User{
		PaymentInfo: userPaymentInfo,
	}).Error; err != nil {
		response.Message = types.ERROR_INTERNAL_SERVER
		return c.JSON(http.StatusInternalServerError, response)
	}

	response.Message = types.SUCCESS
	return c.JSON(http.StatusOK, response)
}
