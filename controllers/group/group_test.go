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
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

var (
	testGroupController = NewGroupController(database.DBTesting.GetConnection())
)

func TestUserCreateGroup(t *testing.T) {
	db := database.DBTesting.GetConnection()

	e := echo.New()

	userA := factories.UserFactory{}
	userA.UserA()

	userB := factories.UserFactory{}
	userB.UserB()

	userC := factories.UserFactory{}
	userC.UserC()

	newGroup := factories.GroupFactory{}
	newGroup.Init()

	id := userA.ID

	newGroup.MemberID = append(newGroup.MemberID, userB.ID)
	newGroup.MemberID = append(newGroup.MemberID, userC.ID)

	userCreateGroupRequest, _ := json.Marshal(requests.UserCreateGroupRequest{
		Name:      newGroup.Name,
		MemberID:  newGroup.MemberID,
		StartDate: newGroup.StartDate,
		EndDate:   newGroup.EndDate,
	})

	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(string(userCreateGroupRequest)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("id", id)

	group := responses.TestResponse[string]{}
	if assert.NoError(t, testGroupController.UserCreateGroup(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)

		if err := json.Unmarshal(rec.Body.Bytes(), &group); err != nil {
			t.Error(err.Error())
		}
	}

	// delete created group
	groupID, err := uuid.Parse(group.Data)
	if err != nil {
		t.Error(err.Error())
	}

	if err := db.Where(&entities.Group{
		GroupID: groupID,
	}).Delete(&entities.Group{}).Error; err != nil {
		t.Error(err.Error())
	}
}

func TestEditGroupInfo(t *testing.T) {
	db := database.DBTesting.GetConnection()
	e := echo.New()

	userA := factories.UserFactory{}
	id := userA.ID

	groupB := factories.GroupFactory{}
	groupB.GroupB()

	newStartDate := time.Now()
	newEndDate := time.Now().Add(time.Hour * 24)

	request := requests.EditGroupInfoRequest{
		GroupID:   groupB.GroupID,
		Name:      "new group b",
		StartDate: newStartDate,
		EndDate:   newEndDate,
	}

	userEditGroupInfoReq, _ := json.Marshal(request)

	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(string(userEditGroupInfoReq)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("id", id)

	group := responses.TestResponse[string]{}
	if assert.NoError(t, testGroupController.EditGroupInfo(c)) {
		assert.Equal(t, http.StatusAccepted, rec.Code)
		if err := json.Unmarshal(rec.Body.Bytes(), &group); err != nil {
			t.Error(err.Error())
		}
	}

	// check from db
	groupDb := entities.Group{}
	if err := db.Where(&entities.Group{GroupID: groupB.GroupID}).Find(&groupDb).Error; err != nil {
		t.Error(err.Error())
	}

	assert.Equal(t, request.Name, groupDb.Name)
	// assert.Equal(t, request.StartDate, groupDb.StartDate)
	// assert.Equal(t, request.EndDate, groupDb.EndDate)
}
