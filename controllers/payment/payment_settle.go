package controllers

import (
	"net/http"
	"split-rex-backend/entities"
	"split-rex-backend/entities/requests"
	"split-rex-backend/entities/responses"
	"split-rex-backend/types"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func (h *paymentController) GetUnsettledPayment(c echo.Context) error {
	db := h.db
	response := entities.Response[[]responses.UnsettledPaymentResponse]{}

	//get user_id from token and group_id from param
	userID := c.Get("id").(uuid.UUID)
	groupID, _ := uuid.Parse(c.QueryParam("group_id"))

	// check if group id present
	group := entities.Group{}
	if err := db.Find(&group, groupID).Error; err != nil {
		response.Message = types.ERROR_INTERNAL_SERVER
		return c.JSON(http.StatusInternalServerError, response)
	}

	// get from payment table where GroupID = groupID and UserID1 = userID and Status = UNPAID
	payments := []entities.Payment{}
	conditionPayment := entities.Payment{GroupID: groupID, UserID1: userID, Status: types.STATUS_PAYMENT_UNPAID}
	if err := db.Where(conditionPayment).Find(&payments).Error; err != nil {
		response.Message = types.ERROR_INTERNAL_SERVER
		return c.JSON(http.StatusInternalServerError, response)
	}

	// move payments to response data
	unsettledTransaction := []responses.UnsettledPaymentResponse{}
	for _, payment := range payments {
		user := entities.User{}
		if err := db.Find(&user, payment.UserID2).Error; err != nil {
			response.Message = types.ERROR_INTERNAL_SERVER
			return c.JSON(http.StatusInternalServerError, response)
		}
		unsettledTransaction = append(unsettledTransaction, responses.UnsettledPaymentResponse{
			PaymentID:   payment.PaymentID,
			UserID:      payment.UserID2,
			Name:        user.Name,
			Color:       user.Color,
			TotalUnpaid: payment.TotalUnpaid,
			TotalPaid:   payment.TotalPaid,
			Status:      payment.Status,
		})
	}

	response.Message = types.SUCCESS
	response.Data = unsettledTransaction
	return c.JSON(http.StatusOK, response)
}

func (h *paymentController) SettlePaymentOwed(c echo.Context) error {
	db := h.db
	response := entities.Response[string]{}

	//get payment_id and total_paid from request body
	request := requests.SettleRequest{}
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

	//check if total_paid > total_unpaid
	if (request.TotalPaid + payment.TotalPaid) > payment.TotalUnpaid {
		response.Message = types.ERROR_TOO_MUCH_PAYMENT
		return c.JSON(http.StatusBadRequest, response)
	}

	//update payment table
	payment.TotalPaid = payment.TotalPaid + request.TotalPaid
	payment.Status = types.STATUS_PAYMENT_PENDING
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
	payment2.TotalPaid = payment2.TotalPaid - request.TotalPaid
	payment2.Status = types.STATUS_PAYMENT_PENDING

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
		Status:            "UNCONFIRMED",
		Amount:            request.TotalPaid,
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

func (h *paymentController) SettlePaymentLent(c echo.Context) error {
	db := h.db
	response := entities.Response[string]{}

	//get payment_id and total_paid from request body
	request := requests.SettleRequest{}
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

	//check if total_paid - total_unpaid < request.TotalPaid
	if (payment.TotalPaid - payment.TotalUnpaid) < request.TotalPaid {
		response.Message = types.ERROR_TOO_MUCH_PAYMENT
		return c.JSON(http.StatusBadRequest, response)
	}

	//update payment table
	payment.TotalUnpaid = payment.TotalUnpaid + request.TotalPaid
	if payment.TotalUnpaid == 0 {
		payment.Status = types.STATUS_PAYMENT_PAID
	} else {
		payment.Status = types.STATUS_PAYMENT_PENDING
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
	payment2.TotalUnpaid = payment2.TotalUnpaid - request.TotalPaid
	if payment2.TotalUnpaid == 0 {
		payment2.Status = types.STATUS_PAYMENT_PAID
	} else {
		payment2.Status = types.STATUS_PAYMENT_PENDING
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
		Amount:            request.TotalPaid,
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
		Amount:     request.TotalPaid,
		Date:       time.Now(),
	}
	if err := db.Create(&groupActivity).Error; err != nil {
		response.Message = types.ERROR_INTERNAL_SERVER
		return c.JSON(http.StatusInternalServerError, response)
	}

	response.Message = types.SUCCESS
	return c.JSON(http.StatusOK, response)
}
