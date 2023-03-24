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
	"split-rex-backend/types"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

var (
	testMetadata       = configs.Config.GetTestMetadata()
	testAuthController = NewAuthController(database.DBTesting.GetConnection(), testMetadata)
)

func CreateTestUser() error {
	db := database.DBTesting.GetConnection()

	userFac := factories.UserFactory{}
	userFac.UserB()
	user := entities.User{
		ID:       uuid.New(),
		Name:     userFac.Name,
		Username: userFac.Username,
		Email:    userFac.Email,
		Password: userFac.Password,
		Groups:   types.ArrayOfUUID{},
	}

	// if user already in db, then no need to create
	check := entities.User{}
	if err := db.Where(&entities.User{Username: user.Username}).Find(&check).Error; err != nil {
		return err
	}
	if check.Username != "" {
		return nil
	}

	if err := db.FirstOrCreate(&user).Error; err != nil {
		return err
	}

	return nil
}

func GetUserTestID() uuid.UUID {
	db := database.DBTesting.GetConnection()

	userFac := factories.UserFactory{}
	userFac.InitAuth()
	user := entities.User{
		Username: userFac.Username,
	}

	userRes := entities.User{}
	db.Where(&user).Find(&userRes)

	return userRes.ID
}

func TestAuth(t *testing.T) {
	db := database.DBTesting.GetConnection()

	// 1. Register account
	e := echo.New()

	user := factories.UserFactory{}
	user.Init(uuid.New())

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

func TestUpdate(t *testing.T) {
	e := echo.New()

	// get user id
	if err := CreateTestUser(); err != nil {
		t.Error(err)
	}
	id := GetUserTestID()

	// 1. update profile
	updateProfileRequest := requests.UpdateProfileRequest{
		Name:  "new",
		Color: 3,
	}

	updateRequest, _ := json.Marshal(updateProfileRequest)

	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(string(updateRequest)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("id", id)

	if assert.NoError(t, testAuthController.UpdateProfileController(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}

	// 2. get profile
	req = httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	c.Set("id", id)

	profileResponse := responses.TestResponse[responses.ProfileResponse]{}
	if assert.NoError(t, testAuthController.ProfileController(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)

		if err := json.Unmarshal(rec.Body.Bytes(), &profileResponse); err != nil {
			t.Error(err.Error())
		}
	}

	// check name and color
	assert.Equal(t, updateProfileRequest.Name, profileResponse.Data.Fullname)
	assert.Equal(t, updateProfileRequest.Color, profileResponse.Data.Color)

}
