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
	"split-rex-backend/entities/responses"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
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

	// 3. update profile
	updateProfileRequest := requests.UpdateProfileRequest{
		Name:     "new",
		Password: string("newpassword"),
		Color:    3,
	}

	updateRequest, _ := json.Marshal(updateProfileRequest)

	req = httptest.NewRequest(http.MethodPost, "/", strings.NewReader(string(updateRequest)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)

	if assert.NoError(t, testAuthController.UpdateProfileController(c)) {
		assert.Equal(t, http.StatusAccepted, rec.Code)
	}

	// 4. get profile
	req = httptest.NewRequest(http.MethodPost, "/", strings.NewReader(string(loginRequest)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)

	profileResponse := responses.TestResponse[string]{}
	if assert.NoError(t, testAuthController.ProfileController(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)
		if err := json.Unmarshal(rec.Body.Bytes(), &profileResponse); err != nil {
			t.Error(err.Error())
		}
	}

	profile := responses.ProfileResponse{}
	if err := json.Unmarshal([]byte(profileResponse.Data), &profile); err != nil {
		t.Error(err.Error())
	}

	// check name and color
	assert.Equal(t, profile.Fullname, updateProfileRequest.Name)
	assert.Equal(t, profile.Color, updateProfileRequest.Color)

	// check password
	userDb := entities.User{}
	conditionUsername := entities.User{Username: user.Username}
	if err := db.Where(&conditionUsername).Find(&userDb).Error; err != nil {
		t.Error(err.Error())
	}

	if err := bcrypt.CompareHashAndPassword(userDb.Password, []byte(updateProfileRequest.Password)); err != nil {
		t.Error(err.Error())
	}

	// 3. Delete registered account for reusability
	db.Where(&entities.User{
		Username: user.Username,
		Email:    user.Email,
		Name:     user.Name,
	}).Delete(&entities.User{})
}
