package controllers

import (
	"fmt"
	"net/http"
	"split-rex-backend/entities"
	"split-rex-backend/entities/requests"
	"split-rex-backend/types"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func (h *transactionController) UserCreateTransaction(c echo.Context) error {
	db := h.db
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
	response.Data = transaction.TransactionID.String()
	return c.JSON(http.StatusCreated, response)
}

func (h *transactionController) UpdatePayment(c echo.Context) error {
	db := h.db
	response := entities.Response[string]{}
	userID := c.Get("id").(uuid.UUID)

	request := requests.UpdatePaymentRequest{}
	if err := c.Bind(&request); err != nil {
		response.Message = err.Error()
		fmt.Println(1)
		return c.JSON(http.StatusBadRequest, response)
	}

	// check if group exist
	group := entities.Group{}
	if err := db.Find(&group, request.GroupID).Error; err != nil {
		response.Message = err.Error()
		fmt.Println(2)
		return c.JSON(http.StatusBadRequest, response)
	}

	// insert into payment table
	for _, payment := range request.ListPayment {
		// check if user exist in group
		memberExist := false
		for _, member := range group.MemberID {
			if member == payment.UserID {
				memberExist = true
				break
			}
		}
		if !memberExist {
			response.Message = types.ERROR_BAD_REQUEST
			fmt.Println(3)
			return c.JSON(http.StatusBadRequest, response)
		}

		// check if userID exist in UserID1 payment table
		tempPayment := entities.Payment{}
		conditionPayment := entities.Payment{UserID1: userID, UserID2: payment.UserID, GroupID: request.GroupID}
		if err := db.Where(&conditionPayment).Find(&tempPayment).Error; err != nil {
			response.Message = err.Error()
			fmt.Println(4)
			return c.JSON(http.StatusBadRequest, response)
		}
		// if payment not exist, create new payment
		if tempPayment.PaymentID == uuid.Nil {
			newPayment := entities.Payment{
				PaymentID:   uuid.New(),
				GroupID:     request.GroupID,
				UserID1:     userID,
				UserID2:     payment.UserID,
				TotalUnpaid: payment.TotalUnpaid,
				TotalPaid:   0,
				Status:      "UNPAID",
			}
			if err := db.Create(&newPayment).Error; err != nil {
				response.Message = err.Error()
				fmt.Println(5)
				return c.JSON(http.StatusInternalServerError, response)
			}
		} else {
			// if payment exist, update payment
			tempPayment.TotalUnpaid = tempPayment.TotalUnpaid + payment.TotalUnpaid
			if tempPayment.Status == "PAID" {
				tempPayment.Status = "UNPAID"
			}
			if err := db.Save(&tempPayment).Error; err != nil {
				response.Message = err.Error()
				fmt.Println(6)
				return c.JSON(http.StatusInternalServerError, response)
			}
		}

		// check if userID exist in UserID2 payment table
		tempPayment = entities.Payment{}
		conditionPayment = entities.Payment{UserID1: payment.UserID, UserID2: userID, GroupID: request.GroupID}
		if err := db.Where(&conditionPayment).Find(&tempPayment).Error; err != nil {
			response.Message = err.Error()
			fmt.Println(7)
			return c.JSON(http.StatusBadRequest, response)
		}
		// if payment not exist, create new payment
		if tempPayment.PaymentID == uuid.Nil {
			newPayment := entities.Payment{
				PaymentID:   uuid.New(),
				GroupID:     request.GroupID,
				UserID1:     payment.UserID,
				UserID2:     userID,
				TotalUnpaid: -payment.TotalUnpaid,
				TotalPaid:   0,
				Status:      "UNPAID",
			}
			if err := db.Create(&newPayment).Error; err != nil {
				response.Message = err.Error()
				fmt.Println(8)
				return c.JSON(http.StatusInternalServerError, response)
			}
		} else {
			// if payment exist, update payment
			tempPayment.TotalUnpaid = tempPayment.TotalUnpaid - payment.TotalUnpaid
			if tempPayment.Status == "PAID" {
				tempPayment.Status = "UNPAID"
			}
			if err := db.Save(&tempPayment).Error; err != nil {
				response.Message = err.Error()
				fmt.Println(9)
				return c.JSON(http.StatusInternalServerError, response)
			}
		}
	}

	response.Message = types.SUCCESS
	response.Data = "success"
	return c.JSON(http.StatusOK, response)

}
