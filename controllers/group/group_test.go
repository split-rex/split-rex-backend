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
	"split-rex-backend/types"
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

	if err := db.Where(&entities.Group{GroupID: groupID}).Delete(&entities.Group{}).Error; err != nil {
		t.Error(err.Error())
	}
}

func TestEditGroupInfo(t *testing.T) {
	db := database.DBTesting.GetConnection()
	e := echo.New()

	userA := factories.UserFactory{}
	id := userA.ID

	group := factories.GroupFactory{}
	group.GroupA()

	newStartDate := time.Now()
	newEndDate := time.Now().Add(time.Hour * 24)

	request := requests.EditGroupInfoRequest{
		GroupID:   group.GroupID,
		Name:      "groupA",
		StartDate: newStartDate,
		EndDate:   newEndDate,
	}

	userEditGroupInfoReq, _ := json.Marshal(request)

	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(string(userEditGroupInfoReq)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("id", id)

	groupRes := responses.TestResponse[string]{}
	if assert.NoError(t, testGroupController.EditGroupInfo(c)) {
		assert.Equal(t, http.StatusAccepted, rec.Code)
		if err := json.Unmarshal(rec.Body.Bytes(), &groupRes); err != nil {
			t.Error(err.Error())
		}
	}

	// check from db
	groupDb := entities.Group{}
	if err := db.Where(&entities.Group{GroupID: group.GroupID}).Find(&groupDb).Error; err != nil {
		t.Error(err.Error())
	}

	assert.Equal(t, request.Name, groupDb.Name)
}

func TestAddNewMemberToGroup(t *testing.T) {
	db := database.DBTesting.GetConnection()
	e := echo.New()

	userAuth := factories.UserFactory{}
	userAuth.InitAuth()
	id := userAuth.ID

	// add user init to db
	userInit := factories.UserFactory{}
	userInit.Init(uuid.New())
	if err := db.Create(&entities.User{
		ID:       userInit.ID,
		Name:     userInit.Name,
		Email:    userInit.Email,
		Username: userInit.Username,
		Password: userInit.Password,
		Groups:   types.ArrayOfUUID{},
	}).Error; err != nil {
		t.Error(err.Error())
	}

	// add userAuth and userInit as a friend
	userAuthFriend := entities.Friend{
		ID:        id,
		Friend_id: types.ArrayOfUUID{},
	}
	userAuthFriend.Friend_id = append(userAuthFriend.Friend_id, userInit.ID)
	if err := db.Create(&userAuthFriend).Error; err != nil {
		t.Error(err.Error())
	}

	// then add userInit to group A
	group := factories.GroupFactory{}
	group.GroupA()
	request := requests.AddGroupMemberRequest{
		Group_id:   group.GroupID,
		Friends_id: types.ArrayOfUUID{},
	}
	request.Friends_id = append(request.Friends_id, userInit.ID)

	// create new request
	addGroupMemberReq, _ := json.Marshal(request)
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(string(addGroupMemberReq)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("id", id)

	groupRes := responses.TestResponse[string]{}
	if assert.NoError(t, testGroupController.AddGroupMember(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		if err := json.Unmarshal(rec.Body.Bytes(), &groupRes); err != nil {
			t.Error(err.Error())
		}
	}

	// check db if groupA already added in userInit
	userInitInDB := entities.User{}
	if err := db.Find(&userInitInDB, userInit.ID).Error; err != nil {
		t.Error(err.Error())
	}
	assert.Equal(t, userInitInDB.Groups.Contains(group.GroupID), true)

	// post condition update groupA, remove userAuth friend, remove user init
	group.MemberID = types.ArrayOfUUID{}
	group.MemberID = append(group.MemberID, userAuth.ID)
	if err := db.Model(&entities.Group{}).Where(&entities.Group{GroupID: group.GroupID}).Update("member_id", group.MemberID).Error; err != nil {
		t.Error(err.Error())
	}

	if err := db.Where(&entities.Friend{ID: id}).Delete(&entities.Friend{}).Error; err != nil {
		t.Error(err.Error())
	}

	if err := db.Where(&entities.User{ID: userInit.ID}).Delete(&entities.User{}).Error; err != nil {
		t.Error(err.Error())
	}
}
