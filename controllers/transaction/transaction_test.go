package controllers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"split-rex-backend/configs/database"
	"split-rex-backend/entities"
	"split-rex-backend/entities/factories"
	"split-rex-backend/entities/requests"
	"split-rex-backend/entities/responses"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

var (
	testUserController = NewTransactionController(database.DBTesting.GetConnection())
)

func TestUserCreateTransaction(t *testing.T) {
	e := echo.New()

	group := factories.GroupFactory{}
	group.GroupA()

	transactionID := uuid.New()

	item := factories.ItemFactory{
		ItemID:        uuid.New(),
		TransactionID: transactionID,
	}
	item.Consumer = append(item.Consumer, group.MemberID...)
	item.Init()

	transaction := factories.TransactionFactory{
		TransactionID: transactionID,
		GroupID:       group.GroupID,
		BillOwner:     group.MemberID[0],
	}
	transaction.Init()

	transaction.Items = append(transaction.Items, item.ItemID)

	requestsNewTransaction := requests.UserCreateTransactionRequest{
		Name:        transaction.Name,
		Description: transaction.Description,
		GroupID:     transaction.GroupID,
		Date:        transaction.Date,
		Subtotal:    transaction.Subtotal,
		Tax:         transaction.Tax,
		Service:     transaction.Service,
		Total:       transaction.Total,
		BillOwner:   transaction.BillOwner,
		Items:       transaction.Items,
	}

	UserCreateTransactionRequest, _ := json.Marshal(requestsNewTransaction)

	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(string(UserCreateTransactionRequest)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	transactionResp := responses.TestResponse[string]{}
	if assert.NoError(t, testUserController.UserCreateTransaction(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)

		if err := json.Unmarshal(rec.Body.Bytes(), &transactionResp); err != nil {
			t.Error(err.Error())
		}
	}

	db := database.DBTesting.GetConnection()
	if err := db.Where(&entities.Transaction{
		TransactionID: uuid.MustParse(transactionResp.Data),
	}).Delete(&entities.Transaction{}).Error; err != nil {
		t.Error(err.Error())
	}

	// TODO: delete created item
}
