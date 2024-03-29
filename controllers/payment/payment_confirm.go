package controllers

import (
	"math"
	"net/http"
	"split-rex-backend/entities"
	"split-rex-backend/entities/requests"
	"split-rex-backend/entities/responses"
	"split-rex-backend/types"
	"time"

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
		if payment.TotalPaid < 0 {
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
	}

	response.Message = types.SUCCESS
	response.Data = unconfirmedTransaction
	return c.JSON(http.StatusOK, response)
}

func (h *paymentController) ConfirmSettle(c echo.Context) error {
	db := h.db
	response := entities.Response[string]{}

	//get payment_id and total_paid from request body
	request := requests.ConfirmRequest{}
	if err := c.Bind(&request); err != nil {
		response.Message = types.ERROR_BAD_REQUEST
		return c.JSON(http.StatusBadRequest, response)
	}

	//get payment from payment_id
	payment := entities.Payment{}
	if err := db.Find(&payment, request.PaymentID).Error; err != nil {
		response.Message = types.ERROR_INTERNAL_SERVER
		return c.JSON(http.StatusInternalServerError, response)
	}

	//update payment table
	paid := math.Abs(payment.TotalPaid)
	payment.TotalUnpaid = payment.TotalUnpaid - payment.TotalPaid
	payment.TotalPaid = 0
	if payment.TotalUnpaid == 0 {
		payment.Status = types.STATUS_PAYMENT_PAID
	} else {
		payment.Status = types.STATUS_PAYMENT_UNPAID
	}

	if err := db.Save(&payment).Error; err != nil {
		response.Message = types.ERROR_INTERNAL_SERVER
		return c.JSON(http.StatusInternalServerError, response)
	}

	//update payment table where UserID2 = payment.UserID1 and UserID1 = payment.UserID2
	conditionPayment := entities.Payment{GroupID: payment.GroupID, UserID1: payment.UserID2, UserID2: payment.UserID1}
	payment2 := entities.Payment{}
	if err := db.Where(&conditionPayment).Find(&payment2).Error; err != nil {
		response.Message = types.ERROR_INTERNAL_SERVER
		return c.JSON(http.StatusInternalServerError, response)
	}
	payment2.TotalUnpaid = payment2.TotalUnpaid - payment2.TotalPaid
	payment2.TotalPaid = 0
	if payment2.TotalUnpaid == 0 {
		payment2.Status = types.STATUS_PAYMENT_PAID
	} else {
		payment2.Status = types.STATUS_PAYMENT_UNPAID
	}

	if err := db.Save(&payment2).Error; err != nil {
		response.Message = types.ERROR_INTERNAL_SERVER
		return c.JSON(http.StatusInternalServerError, response)
	}

	//get user from payment.UserID1
	user := entities.User{}
	if err := db.Find(&user, payment.UserID1).Error; err != nil {
		response.Message = types.ERROR_INTERNAL_SERVER
		return c.JSON(http.StatusInternalServerError, response)
	}
	// get group from payment.GroupID
	group := entities.Group{}
	if err := db.Find(&group, payment.GroupID).Error; err != nil {
		response.Message = types.ERROR_INTERNAL_SERVER
		return c.JSON(http.StatusInternalServerError, response)
	}

	// insert to activity and payment activity table
	newID := uuid.New()
	paymentActivity := entities.PaymentActivity{
		PaymentActivityID: newID,
		Name:              user.Name,
		Status:            "CONFIRMED",
		Amount:            paid,
		GroupName:         group.Name,
	}
	if err := db.Create(&paymentActivity).Error; err != nil {
		response.Message = types.ERROR_INTERNAL_SERVER
		return c.JSON(http.StatusInternalServerError, response)
	}

	activity := entities.Activity{
		ActivityID:   uuid.New(),
		ActivityType: "PAYMENT",
		UserID:       payment.UserID2,
		Date:         time.Now(),
		RedirectID:   payment.GroupID,
		DetailID:     newID,
	}
	if err := db.Create(&activity).Error; err != nil {
		response.Message = types.ERROR_INTERNAL_SERVER
		return c.JSON(http.StatusInternalServerError, response)
	}

	// insert to group activity table
	groupActivity := entities.GroupActivity{
		ActivityID: uuid.New(),
		GroupID:    payment.GroupID,
		UserID1:    payment.UserID2,
		UserID2:    payment.UserID1,
		Amount:     paid,
		Date:       time.Now(),
	}
	if err := db.Create(&groupActivity).Error; err != nil {
		response.Message = types.ERROR_INTERNAL_SERVER
		return c.JSON(http.StatusInternalServerError, response)
	}

	response.Message = types.SUCCESS
	return c.JSON(http.StatusOK, response)

}

func (h *paymentController) DenySettle(c echo.Context) error {
	db := h.db
	response := entities.Response[string]{}

	//get payment_id and total_paid from request body
	request := requests.ConfirmRequest{}
	if err := c.Bind(&request); err != nil {
		response.Message = types.ERROR_BAD_REQUEST
		return c.JSON(http.StatusBadRequest, response)
	}

	//get payment from payment_id
	payment := entities.Payment{}
	if err := db.Find(&payment, request.PaymentID).Error; err != nil {
		response.Message = types.ERROR_INTERNAL_SERVER
		return c.JSON(http.StatusInternalServerError, response)
	}

	//update payment table
	paid := math.Abs(payment.TotalPaid)
	payment.TotalPaid = 0
	if payment.TotalUnpaid == 0 {
		payment.Status = types.STATUS_PAYMENT_PAID
	} else {
		payment.Status = types.STATUS_PAYMENT_UNPAID
	}

	if err := db.Save(&payment).Error; err != nil {
		response.Message = types.ERROR_INTERNAL_SERVER
		return c.JSON(http.StatusInternalServerError, response)
	}

	//update payment table where UserID2 = payment.UserID1 and UserID1 = payment.UserID2
	conditionPayment := entities.Payment{GroupID: payment.GroupID, UserID1: payment.UserID2, UserID2: payment.UserID1}
	payment2 := entities.Payment{}
	if err := db.Where(&conditionPayment).Find(&payment2).Error; err != nil {
		response.Message = types.ERROR_INTERNAL_SERVER
		return c.JSON(http.StatusInternalServerError, response)
	}
	payment2.TotalPaid = 0
	if payment2.TotalUnpaid == 0 {
		payment2.Status = types.STATUS_PAYMENT_PAID
	} else {
		payment2.Status = types.STATUS_PAYMENT_UNPAID
	}

	if err := db.Save(&payment2).Error; err != nil {
		response.Message = types.ERROR_INTERNAL_SERVER
		return c.JSON(http.StatusInternalServerError, response)
	}

	//get user from payment.UserID1
	user := entities.User{}
	if err := db.Find(&user, payment.UserID1).Error; err != nil {
		response.Message = types.ERROR_INTERNAL_SERVER
		return c.JSON(http.StatusInternalServerError, response)
	}
	// get group from payment.GroupID
	group := entities.Group{}
	if err := db.Find(&group, payment.GroupID).Error; err != nil {
		response.Message = types.ERROR_INTERNAL_SERVER
		return c.JSON(http.StatusInternalServerError, response)
	}

	// insert to activity and payment activity table
	newID := uuid.New()
	paymentActivity := entities.PaymentActivity{
		PaymentActivityID: newID,
		Name:              user.Name,
		Status:            "DENIED",
		Amount:            paid,
		GroupName:         group.Name,
	}
	if err := db.Create(&paymentActivity).Error; err != nil {
		response.Message = types.ERROR_INTERNAL_SERVER
		return c.JSON(http.StatusInternalServerError, response)
	}

	activity := entities.Activity{
		ActivityID:   uuid.New(),
		ActivityType: "PAYMENT",
		UserID:       payment.UserID2,
		Date:         time.Now(),
		RedirectID:   payment.GroupID,
		DetailID:     newID,
	}
	if err := db.Create(&activity).Error; err != nil {
		response.Message = types.ERROR_INTERNAL_SERVER
		return c.JSON(http.StatusInternalServerError, response)
	}

	response.Message = types.SUCCESS
	return c.JSON(http.StatusOK, response)

}
