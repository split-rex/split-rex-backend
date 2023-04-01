package controllers

import (
	"net/http"
	"split-rex-backend/entities"
	"split-rex-backend/entities/requests"
	"split-rex-backend/types"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func (con *authController) AddPaymentInfo(c echo.Context) error {
	db := con.db
	response := entities.Response[string]{}

	addPaymentInfoRequest := requests.AddPaymentInfoRequest{}
	if err := c.Bind(&addPaymentInfoRequest); err != nil {
		response.Message = types.ERROR_BAD_REQUEST
		return c.JSON(http.StatusBadRequest, response)
	}

	// get id from context, cast to uuid
	user_id := c.Get("id").(uuid.UUID)

	// search the id of the user
	user := entities.User{}
	conditionUser := entities.User{ID: user_id}
	if err := db.Where(&conditionUser).Find(&user).Error; err != nil {
		response.Message = types.ERROR_INTERNAL_SERVER
		return c.JSON(http.StatusInternalServerError, response)
	}
	
	// error if user not found
	if user.ID == uuid.Nil {
		response.Message = types.ERROR_BAD_REQUEST
		return c.JSON(http.StatusInternalServerError, response)
	}

	// get all the user's payment info
	userPaymentInfo := user.PaymentInfo

	numberAndNameInfo := map[int]string{
		int(addPaymentInfoRequest.Account_number): addPaymentInfoRequest.Account_name,
	}
	// val is the value of "addPaymentInfoRequest.Payment_method" from the map if it exists, or a "zero value" if it doesn't.
	val, ok := userPaymentInfo[addPaymentInfoRequest.Payment_method]
	if ok {
		val = append(val, numberAndNameInfo)
		userPaymentInfo[addPaymentInfoRequest.Payment_method] = val
	} else {
		paymentInfo := []map[int]string{numberAndNameInfo}
		userPaymentInfo[addPaymentInfoRequest.Payment_method] = paymentInfo
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
