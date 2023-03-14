package controllers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"split-rex-backend/configs/database"
	"split-rex-backend/entities"
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
	transactionJson := `{
		"name": "New Transaction",
		"description": "New Transaction Description",
		"group_id" : "0b865d7f-e40e-4440-905e-eccf2caaa6ed",
		"date" : "2023-02-07T17:19:20.968831+07:00",
		"subtotal" : 1000.0,
		"tax" : 100.0,
		"service" : 100.0,
		"total" : 1200.0,
		"bill_owner" : "6251ac85-e43d-4b88-8779-588099df5008",
		"items" : ["6251ac85-e43d-4b88-8779-588099df5008"]
	}`

	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(transactionJson))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	transaction := responses.TestResponse[string]{}
	if assert.NoError(t, testUserController.UserCreateTransaction(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)

		if err := json.Unmarshal(rec.Body.Bytes(), &transaction); err != nil {
			t.Error(err.Error())
		}
	}

	db := database.DBTesting.GetConnection()
	if err := db.Where(&entities.Transaction{
		TransactionID: uuid.MustParse(transaction.Data),
	}).Delete(&entities.Transaction{}).Error; err != nil {
		t.Error(err.Error())
	}
}
