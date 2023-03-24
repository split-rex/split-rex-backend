package controllers

import (
	"encoding/json"
	"fmt"
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

	// init new group
	group := factories.GroupFactory{}
	group.GroupA()

	// make new transaction factory
	transaction := factories.TransactionFactory{
		TransactionID: uuid.New(),
		GroupID:       group.GroupID,
		BillOwner:     group.MemberID[0],
	}
	transaction.Init()

	// create new item with consumer from members in group
	item := factories.ItemFactory{}
	item.Consumer = append(item.Consumer, group.MemberID...)
	item.Init()

	// change into item requests
	itemRequest := requests.ItemRequest{
		Name:     item.Name,
		Quantity: item.Quantity,
		Price:    item.Price,
		Consumer: item.Consumer,
	}

	// then append item to transaction
	transaction.Items = append(transaction.Items, itemRequest)

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
		fmt.Println(rec.Body.String())
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
