package controllers

import (
	"net/http"
	"net/http/httptest"
	"split-rex-backend/configs"
	"split-rex-backend/configs/database"
	"split-rex-backend/entities"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

var (
	testMetadata       = configs.Config.GetTestMetadata()
	testAuthController = NewAuthController(database.DBTesting.GetConnection(), testMetadata)
)

func TestAuth(t *testing.T) {
	db := database.DBTesting.GetConnection()

	// 1. Register account
	e := echo.New()
	registerJson := `{
		"name": "testing",
		"email": "testing@gmail.com",
		"username": "testing",
		"password": "testing"
	}`

	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(registerJson))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if assert.NoError(t, testAuthController.RegisterController(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)
	}

	// 2. Login to registered account
	loginJson := `{
		"email": "testing@gmail.com",
		"password": "testing"
	}`

	req = httptest.NewRequest(http.MethodPost, "/", strings.NewReader(loginJson))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)

	if assert.NoError(t, testAuthController.LoginController(c)) {
		assert.Equal(t, http.StatusAccepted, rec.Code)
	}

	// 3. Delete registered account for reusability
	user := entities.User{}
	db.Where(&entities.User{
		Username: "testing",
		Email:    "testing@gmail.com",
		Name:     "testing",
	}).Delete(&user)
}
