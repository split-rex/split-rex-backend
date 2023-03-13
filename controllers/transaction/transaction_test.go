package controllers

import (
	"net/http"
	"net/http/httptest"
	"split-rex-backend/configs/database"
	"strings"
	"testing"

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

	if assert.NoError(t, testUserController.UserCreateTransaction(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)
	}
}

// func TestUserCreateTransaction(t *testing.T) {
// 	groupID, _ := uuid.Parse("6251ac85-e43d-4b88-8779-588099df5008")
// 	billOwnerID, _ := uuid.Parse("6251ac85-e43d-4b88-8779-588099df5008")

// 	date := time.Date(2023, 3, 3, 0, 0, 0, 0, time.UTC)

// transaction := &requests.UserCreateTransactionRequest{
// 	Name:        "Transaction 1",
// 	Description: "Transaction 1 Description",
// 	GroupID:     groupID,
// 	Date:        date,
// 	Subtotal:    100,
// 	Tax:         10,
// 	Service:     10,
// 	Total:       120,
// 	BillOwner:   billOwnerID,
// 	Items:       types.ArrayOfUUID{},
// }
// 	body, _ := json.Marshal(transaction)

// 	res, err := http.Post("http://localhost:8080/userCreateTransaction",
// 		"application/json", bytes.NewBuffer(body))
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	defer res.Body.Close()

// 	if res.StatusCode != http.StatusOK {
// 		t.Error("Expected status code 200, got ", res.StatusCode)
// 	}
// }
