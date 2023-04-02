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

	// create items
	items, err := createItems(db, request.Items)
	if err != nil {
		response.Message = err.Error()
		return c.JSON(http.StatusInternalServerError, response)
	}

	transactionID := uuid.New()
	transaction := entities.Transaction{
		TransactionID: transactionID,
		Name:          request.Name,
		Description:   request.Description,
		GroupID:       request.GroupID,
		Date:          request.Date,
		Subtotal:      request.Subtotal,
		Tax:           request.Tax,
		Service:       request.Service,
		Total:         request.Total,
		BillOwner:     request.BillOwner,
		Items:         items,
	}

	if err := db.Create(&transaction).Error; err != nil {
		response.Message = err.Error()
		return c.JSON(http.StatusInternalServerError, response)
	}

	// get set of consumers
	users := make(map[uuid.UUID]bool)
	for _, item := range request.Items {
		for _, consumer := range item.Consumer {
			users[consumer] = true
		}
	}

	// get name of billowner
	billOwner := entities.User{}
	if err := db.Find(&billOwner, request.BillOwner).Error; err != nil {
		response.Message = err.Error()
		return c.JSON(http.StatusInternalServerError, response)
	}

	// create activity and transaction activity
	newID := uuid.New()
	transactionActivity := entities.TransactionActivity{
		TransactionActivityID: newID,
		Name:                  billOwner.Name,
		GroupName:             group.Name,
	}
	if err := db.Create(&transactionActivity).Error; err != nil {
		response.Message = err.Error()
		return c.JSON(http.StatusInternalServerError, response)
	}
	for userID := range users {
		activity := entities.Activity{
			ActivityID:   uuid.New(),
			ActivityType: "TRANSACTION",
			UserID:       userID,
			Date:         time.Now(),
			RedirectID:   transactionID,
			DetailID:     newID,
		}
		if err := db.Create(&activity).Error; err != nil {
			response.Message = err.Error()
			return c.JSON(http.StatusInternalServerError, response)
		}
	}

	response.Message = types.SUCCESS
	response.Data = transaction.TransactionID.String()
	return c.JSON(http.StatusCreated, response)
}

func (h *transactionController) GetTransactionDetail(c echo.Context) error {
	db := h.db
	response := entities.Response[responses.TransactionDetailResponse]{}

	// get transaction id from param
	transactionID, _ := uuid.Parse(c.QueryParam("transaction_id"))

	// get transaction from transaction_id
	transaction := entities.Transaction{}
	if err := db.Find(&transaction, transactionID).Error; err != nil {
		response.Message = err.Error()
		return c.JSON(http.StatusBadRequest, response)
	}

	// get group from group_id
	group := entities.Group{}
	if err := db.Find(&group, transaction.GroupID).Error; err != nil {
		response.Message = err.Error()
		return c.JSON(http.StatusBadRequest, response)
	}

	// get items of transaction
	itemResponse := []responses.ItemDetailResponse{}
	for _, itemID := range transaction.Items {
		//get consumers of transaction
		consumerResponse := []responses.ConsumerDetailResponse{}
		item := entities.Item{}
		if err := db.Find(&item, itemID).Error; err != nil {
			response.Message = err.Error()
			return c.JSON(http.StatusBadRequest, response)
		}
		for _, consumerID := range item.Consumer {
			consumer := entities.User{}
			if err := db.Find(&consumer, consumerID).Error; err != nil {
				response.Message = err.Error()
				return c.JSON(http.StatusBadRequest, response)
			}
			consumerResponse = append(consumerResponse, responses.ConsumerDetailResponse{
				UserID: consumer.ID,
				Name:   consumer.Name,
				Color:  consumer.Color,
			})
		}

		itemResponse = append(itemResponse, responses.ItemDetailResponse{
			ItemID:     item.ItemID,
			Name:       item.Name,
			Quantity:   item.Quantity,
			Price:      item.Price,
			TotalPrice: float64(item.Quantity) * item.Price,
			Consumer:   consumerResponse,
		})
	}

	// return response
	response.Message = types.SUCCESS
	response.Data = responses.TransactionDetailResponse{
		Name:        transaction.Name,
		Description: transaction.Description,
		GroupID:     transaction.GroupID,
		GroupName:   group.Name,
		Date:        transaction.Date,
		Items:       itemResponse,
		Subtotal:    transaction.Subtotal,
		Tax:         transaction.Tax,
		Service:     transaction.Service,
		Total:       transaction.Total,
	}

	return c.JSON(http.StatusOK, response)
}
