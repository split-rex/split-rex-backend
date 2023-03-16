package controllers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"split-rex-backend/configs"
	"split-rex-backend/configs/database"
	"split-rex-backend/entities"
	"split-rex-backend/entities/factories"
	"split-rex-backend/entities/requests"
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

	user := factories.UserFactory{}
	user.Init()
	
	registerRequest, _ := json.Marshal(requests.RegisterRequest{
		Name:     user.Name,
		Email:    user.Email,
		Username: user.Username,
		Password: string(user.Password),
	})

	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(string(registerRequest)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if assert.NoError(t, testAuthController.RegisterController(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)
	}

	// 2. Login to registered account
	loginRequest, _ := json.Marshal(requests.LoginRequest{
		Email:    user.Email,
		Password: string(user.Password),
	})

	req = httptest.NewRequest(http.MethodPost, "/", strings.NewReader(string(loginRequest)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)

	if assert.NoError(t, testAuthController.LoginController(c)) {
		assert.Equal(t, http.StatusAccepted, rec.Code)
	}

	// 3. Delete registered account for reusability
	db.Where(&entities.User{
		Username: user.Username,
		Email:    user.Email,
		Name:     user.Name,
	}).Delete(&entities.User{})
}
