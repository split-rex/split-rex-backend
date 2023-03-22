package controllers

import (
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
		response.Message = types.ERROR_BAD_REQUEST
		return c.JSON(http.StatusBadRequest, response)
	}

	// check if group exist
	group := entities.Group{}
	if err := db.Find(&group, request.GroupID).Error; err != nil {
		response.Message = types.ERROR_INTERNAL_SERVER
		return c.JSON(http.StatusInternalServerError, response)
	}
	if group.GroupID == uuid.Nil {
		response.Message = types.ERROR_GROUP_NOT_FOUND
		return c.JSON(http.StatusBadRequest, response)
	}

	// TRANSACTION
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

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
			return c.JSON(http.StatusBadRequest, response)
		}

		// check if userID exist in UserID1 payment table
		tempPayment := entities.Payment{}
		conditionPayment := entities.Payment{UserID1: userID, UserID2: payment.UserID, GroupID: request.GroupID}
		if err := db.Where(&conditionPayment).Find(&tempPayment).Error; err != nil {
			response.Message = types.ERROR_INTERNAL_SERVER
			return c.JSON(http.StatusInternalServerError, response)
		}
		// if payment not exist, create new payment
		if tempPayment.PaymentID == uuid.Nil {
			newPayment := entities.Payment{
				PaymentID:   uuid.New(),
				GroupID:     request.GroupID,
				UserID1:     userID,
				UserID2:     payment.UserID,
				TotalUnpaid: -payment.TotalUnpaid,
				TotalPaid:   0,
				Status:      "UNPAID",
			}
			if err := tx.Create(&newPayment).Error; err != nil {
				response.Message = types.ERROR_INTERNAL_SERVER
				return c.JSON(http.StatusInternalServerError, response)
			}
		} else {
			// if payment exist, update payment
			tempPayment.TotalUnpaid = tempPayment.TotalUnpaid - payment.TotalUnpaid
			if tempPayment.Status == "PAID" {
				tempPayment.Status = "UNPAID"
			}
			if err := tx.Save(&tempPayment).Error; err != nil {
				response.Message = types.ERROR_INTERNAL_SERVER
				return c.JSON(http.StatusInternalServerError, response)
			}
		}

		// check if userID exist in UserID2 payment table
		tempPayment = entities.Payment{}
		conditionPayment = entities.Payment{UserID1: payment.UserID, UserID2: userID, GroupID: request.GroupID}
		if err := db.Where(&conditionPayment).Find(&tempPayment).Error; err != nil {
			response.Message = types.ERROR_INTERNAL_SERVER
			return c.JSON(http.StatusInternalServerError, response)
		}
		// if payment not exist, create new payment
		if tempPayment.PaymentID == uuid.Nil {
			newPayment := entities.Payment{
				PaymentID:   uuid.New(),
				GroupID:     request.GroupID,
				UserID1:     payment.UserID,
				UserID2:     userID,
				TotalUnpaid: payment.TotalUnpaid,
				TotalPaid:   0,
				Status:      "UNPAID",
			}
			if err := tx.Create(&newPayment).Error; err != nil {
				response.Message = types.ERROR_INTERNAL_SERVER
				return c.JSON(http.StatusInternalServerError, response)
			}
		} else {
			// if payment exist, update payment
			tempPayment.TotalUnpaid = tempPayment.TotalUnpaid + payment.TotalUnpaid
			if tempPayment.Status == "PAID" {
				tempPayment.Status = "UNPAID"
			}
			if err := tx.Save(&tempPayment).Error; err != nil {
				response.Message = types.ERROR_INTERNAL_SERVER
				return c.JSON(http.StatusInternalServerError, response)
			}
		}
	}

	// commit transaction
	if err := tx.Commit().Error; err != nil {
		response.Message = types.ERROR_INTERNAL_SERVER
		return c.JSON(http.StatusInternalServerError, response)
	}

	response.Message = types.SUCCESS
	response.Data = "success"
	return c.JSON(http.StatusOK, response)

}

func (h *transactionController) ResolveTransaction(c echo.Context) error {
	db := h.db
	response := entities.Response[string]{}

	request := requests.ResolveTransactionRequest{}
	if err := c.Bind(&request); err != nil {
		response.Message = types.ERROR_BAD_REQUEST
		return c.JSON(http.StatusBadRequest, response)
	}

	// check if group exist
	group := entities.Group{}
	if err := db.Find(&group, request.GroupID).Error; err != nil {
		response.Message = types.ERROR_INTERNAL_SERVER
		return c.JSON(http.StatusInternalServerError, response)
	}
	if group.GroupID == uuid.Nil {
		response.Message = types.ERROR_GROUP_NOT_FOUND
		return c.JSON(http.StatusBadRequest, response)
	}

	// get payment with groupID
	payments := []entities.Payment{}
	conditionPayment := entities.Payment{GroupID: request.GroupID, Status: "UNPAID"}
	if err := db.Where(&conditionPayment).Find(&payments).Error; err != nil {
		response.Message = types.ERROR_INTERNAL_SERVER
		return c.JSON(http.StatusInternalServerError, response)
	}

	// init map to store net balance of each user
	balance := make(map[uuid.UUID]float64)
	for _, member := range group.MemberID {
		balance[member] = 0
	}
	for _, payment := range payments {
		if payment.TotalUnpaid > 0 {
			balance[payment.UserID1] = balance[payment.UserID1] + payment.TotalUnpaid
			balance[payment.UserID2] = balance[payment.UserID2] - payment.TotalUnpaid
		}
	}

	// resolve payment
	i := 0
	j := 0
	updatePayment := []entities.Payment{}
	for i < len(group.MemberID) && j < len(group.MemberID) {
		if balance[group.MemberID[i]] <= 0 {
			i = i + 1
		} else if balance[group.MemberID[j]] >= 0 {
			j = j + 1
		} else {
			m := 0.0
			if balance[group.MemberID[i]] < -balance[group.MemberID[j]] {
				m = balance[group.MemberID[i]]
			} else {
				m = -balance[group.MemberID[j]]
			}
			balance[group.MemberID[i]] = balance[group.MemberID[i]] - m
			balance[group.MemberID[j]] = balance[group.MemberID[j]] + m

			// append updatePayment, where userID1 = group.MemberID[i] and userID2 = group.MemberID[j]
			tempPayment := entities.Payment{}
			conditionPayment := entities.Payment{UserID1: group.MemberID[i], UserID2: group.MemberID[j], GroupID: request.GroupID}
			if err := db.Where(&conditionPayment).Find(&tempPayment).Error; err != nil {
				response.Message = types.ERROR_INTERNAL_SERVER
				return c.JSON(http.StatusInternalServerError, response)
			}
			if tempPayment.PaymentID == uuid.Nil {
				newPayment := entities.Payment{
					PaymentID:   uuid.New(),
					GroupID:     request.GroupID,
					UserID1:     group.MemberID[i],
					UserID2:     group.MemberID[j],
					TotalUnpaid: m,
					TotalPaid:   0,
					Status:      "UNPAID",
				}
				updatePayment = append(updatePayment, newPayment)
			} else {
				tempPayment.TotalUnpaid = m
				tempPayment.TotalPaid = 0
				tempPayment.Status = "UNPAID"
				updatePayment = append(updatePayment, tempPayment)
			}

			// append updatePayment, where userID1 = group.MemberID[j] and userID2 = group.MemberID[i]
			tempPayment = entities.Payment{}
			conditionPayment = entities.Payment{UserID1: group.MemberID[j], UserID2: group.MemberID[i], GroupID: request.GroupID}
			if err := db.Where(&conditionPayment).Find(&tempPayment).Error; err != nil {
				response.Message = types.ERROR_INTERNAL_SERVER
				return c.JSON(http.StatusInternalServerError, response)
			}
			if tempPayment.PaymentID == uuid.Nil {
				newPayment := entities.Payment{
					PaymentID:   uuid.New(),
					GroupID:     request.GroupID,
					UserID1:     group.MemberID[j],
					UserID2:     group.MemberID[i],
					TotalUnpaid: -m,
					TotalPaid:   0,
					Status:      "UNPAID",
				}
				updatePayment = append(updatePayment, newPayment)
			} else {
				tempPayment.TotalUnpaid = -m
				tempPayment.TotalPaid = 0
				tempPayment.Status = "UNPAID"
				updatePayment = append(updatePayment, tempPayment)
			}
		}
	}
	// TRANSACTION
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// update payment table totalUnpaid = 0 and totalPaid = 0 and status = "Paid" for all payment with groupID
	for _, payment := range payments {
		payment.TotalUnpaid = 0
		payment.TotalPaid = 0
		payment.Status = "PAID"
		if err := tx.Save(&payment).Error; err != nil {
			response.Message = types.ERROR_INTERNAL_SERVER
			return c.JSON(http.StatusInternalServerError, response)
		}
	}

	// update payment table from updatePayment
	for _, payment := range updatePayment {
		if err := tx.Save(&payment).Error; err != nil {
			response.Message = types.ERROR_INTERNAL_SERVER
			return c.JSON(http.StatusInternalServerError, response)
		}
	}

	// commit transaction
	if err := tx.Commit().Error; err != nil {
		response.Message = types.ERROR_INTERNAL_SERVER
		return c.JSON(http.StatusInternalServerError, response)
	}

	response.Message = types.SUCCESS
	response.Data = "success"
	return c.JSON(http.StatusOK, response)
}
