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

	addPaymentInfoRequest := requests.PaymentInfo{}
	if err := c.Bind(&addPaymentInfoRequest); err != nil {
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

	// val is the value of "addPaymentInfoRequest.Payment_method" from the map if it exists, or a "zero value" if it doesn't.
	_, ok := userPaymentInfo[addPaymentInfoRequest.Payment_method]
	if !ok {
		userPaymentInfo[addPaymentInfoRequest.Payment_method] = make(map[int]string)
	}

	_, ok = userPaymentInfo[addPaymentInfoRequest.Payment_method][int(addPaymentInfoRequest.Account_number)]
	if ok {
		response.Message = types.ERROR_PAYMENT_INFO_ALREADY_EXISTED
		return c.JSON(http.StatusBadRequest, response)
	}

	userPaymentInfo[addPaymentInfoRequest.Payment_method][int(addPaymentInfoRequest.Account_number)] = addPaymentInfoRequest.Account_name

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
