package controllers

import (
	"net/http"
	"split-rex-backend/entities"
	"split-rex-backend/entities/requests"
	"split-rex-backend/types"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func (con *authController) EditPaymentInfo(c echo.Context) error {
	db := con.db
	response := entities.Response[string]{}

	editPaymentInfo := requests.EditPaymentInfoRequest{}
	if err := c.Bind(&editPaymentInfo); err != nil {
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

	// DELETE PROCESS FOR OLD PAYMENT INFO
	// check if payment method available
	payment, ok := userPaymentInfo[editPaymentInfo.Old_payment_method]
	if !ok {
		response.Message = types.ERROR_BAD_REQUEST
		return c.JSON(http.StatusInternalServerError, response)
	}

	// check if account number available
	_, ok = payment[int(editPaymentInfo.Old_account_number)]
	if !ok {
		response.Message = types.ERROR_BAD_REQUEST
		return c.JSON(http.StatusInternalServerError, response)
	}

	// delete account number and account name
	delete(userPaymentInfo[editPaymentInfo.Old_payment_method], int(editPaymentInfo.Old_account_number))

	// if account number in payment method = 0, delete payment method
	if len(userPaymentInfo[editPaymentInfo.Old_payment_method]) < 1 {
		delete(userPaymentInfo, editPaymentInfo.Old_payment_method)
	}

	// ADD PROCESS FOR NEW PAYMENT INFO
	// val is the value of "addPaymentInfoRequest.Payment_method" from the map if it exists, or a "zero value" if it doesn't.
	_, ok = userPaymentInfo[editPaymentInfo.New_payment_method]
	if !ok {
		userPaymentInfo[editPaymentInfo.New_payment_method] = make(map[int]string)
	}

	_, ok = userPaymentInfo[editPaymentInfo.New_payment_method][int(editPaymentInfo.New_account_number)]
	if ok {
		// change back to original payment info
		userPaymentInfo = user.PaymentInfo
		response.Message = types.ERROR_PAYMENT_INFO_ALREADY_EXISTED
		return c.JSON(http.StatusBadRequest, response)
	}

	userPaymentInfo[editPaymentInfo.New_payment_method][int(editPaymentInfo.New_account_number)] = editPaymentInfo.New_account_name

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
