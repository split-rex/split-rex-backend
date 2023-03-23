package controllers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
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
	testMetadata         = configs.Config.GetTestMetadata()
	testFriendController = NewFriendController(database.DBTesting.GetConnection())
)

func searchAndMakeRequest(t *testing.T, user_id uuid.UUID, friendUsername string) {
	e := echo.New()

	q := make(url.Values)
	q.Set("username", friendUsername)
	req := httptest.NewRequest(http.MethodGet, "/?"+q.Encode(), nil)

	// search user to add
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("id", user_id)

	searchUserToAddResp := responses.TestResponse[responses.ProfileResponse]{}
	if assert.NoError(t, testFriendController.SearchUserToAdd(c)) {
		if err := json.Unmarshal(rec.Body.Bytes(), &searchUserToAddResp); err != nil {
			t.Error(err.Error())
		}
		if searchUserToAddResp.Message == types.ERROR_ALREADY_FRIEND ||
			searchUserToAddResp.Message == types.ERROR_ALREADY_REQUESTED_RECEIVED ||
			searchUserToAddResp.Message == types.ERROR_ALREADY_REQUESTED_SENT ||
			searchUserToAddResp.Message == types.ERROR_CANNOT_ADD_SELF {

			assert.Equal(t, http.StatusConflict, rec.Code)
		} else {
			assert.Equal(t, http.StatusOK, rec.Code)
			if err := json.Unmarshal(rec.Body.Bytes(), &searchUserToAddResp); err != nil {
				t.Error(err.Error())
			}
		}
	}

	friend_id := searchUserToAddResp.Data.User_id

	makeFriendRequest, _ := json.Marshal(requests.FriendRequest{
		Friend_id: friend_id,
	})

	req = httptest.NewRequest(http.MethodPost, "/", strings.NewReader(string(makeFriendRequest)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	c.Set("id", user_id)

	response := responses.TestResponse[string]{}

	if assert.NoError(t, testFriendController.MakeFriendRequest(c)) {
		if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
			t.Error(err.Error())
		}
		if searchUserToAddResp.Message == types.ERROR_ALREADY_FRIEND ||
			searchUserToAddResp.Message == types.ERROR_ALREADY_REQUESTED_RECEIVED ||
			searchUserToAddResp.Message == types.ERROR_ALREADY_REQUESTED_SENT ||
			searchUserToAddResp.Message == types.ERROR_CANNOT_ADD_SELF {
			assert.Equal(t, http.StatusConflict, rec.Code)
		} else {
			assert.Equal(t, http.StatusOK, rec.Code)
		}
	}
}

func getFriendRequestSent(t *testing.T, user_id uuid.UUID) types.ArrayOfUUID {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("id", user_id)

	userFriendReqSentResp := responses.TestResponse[[]responses.ProfileResponse]{}
	if assert.NoError(t, testFriendController.FriendRequestSent(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)

		if err := json.Unmarshal(rec.Body.Bytes(), &userFriendReqSentResp); err != nil {
			t.Error(err.Error())
		}
	}

	requestSent := types.ArrayOfUUID{}

	for _, friend := range userFriendReqSentResp.Data {
		requestSent = append(requestSent, uuid.MustParse(friend.User_id))
	}

	return requestSent
}

func getFriendRequestReceived(t *testing.T, user_id uuid.UUID) types.ArrayOfUUID {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("id", user_id)

	userFriendReqRecvResp := responses.TestResponse[[]responses.ProfileResponse]{}
	if assert.NoError(t, testFriendController.FriendRequestReceived(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)

		if err := json.Unmarshal(rec.Body.Bytes(), &userFriendReqRecvResp); err != nil {
			t.Error(err.Error())
		}
	}

	requestRecv := types.ArrayOfUUID{}

	for _, friend := range userFriendReqRecvResp.Data {
		requestRecv = append(requestRecv, uuid.MustParse(friend.User_id))
	}

	return requestRecv
}

func rejectRequest(t *testing.T, user_id uuid.UUID, friend_id uuid.UUID) {
	e := echo.New()

	requestBody, _ := json.Marshal(requests.FriendRequest{
		Friend_id: friend_id.String(),
	})

	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(string(requestBody)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("id", user_id)

	if assert.NoError(t, testFriendController.RejectRequest(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}

}

func acceptRequest(t *testing.T, user_id uuid.UUID, friend_id uuid.UUID) {
	e := echo.New()

	requestBody, _ := json.Marshal(requests.FriendRequest{
		Friend_id: friend_id.String(),
	})

	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(string(requestBody)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("id", user_id)

	if assert.NoError(t, testFriendController.AcceptRequest(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}
}

func getUserFriendList(t *testing.T, user_id uuid.UUID) types.ArrayOfUUID {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("id", user_id)

	userFriendListResp := responses.TestResponse[[]responses.ProfileResponse]{}
	if assert.NoError(t, testFriendController.UserFriendList(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)

		if err := json.Unmarshal(rec.Body.Bytes(), &userFriendListResp); err != nil {
			t.Error(err.Error())
		}
	}

	friendList := types.ArrayOfUUID{}

	for _, friend := range userFriendListResp.Data {
		friendList = append(friendList, uuid.MustParse(friend.User_id))
	}

	return friendList
}

func addUserToDb(t *testing.T, uf factories.UserFactory) {
	db := database.DBTesting.GetConnection()
	user := entities.User{
		ID:       uf.ID,
		Name:     uf.Name,
		Username: uf.Username,
		Email:    uf.Email,
		Password: uf.Password,
		Groups:   types.ArrayOfUUID{},
	}
	if err := db.Create(&user).Error; err != nil {
		t.Error(err.Error())
	}
}

func TestFriendMain(t *testing.T) {
	/*
		TEST FLOW:
		1. A add D
		2. B add A
		3. C add A
		4. A reject B
		5. A approve C
		6. D approve A
	*/
	db := database.DBTesting.GetConnection()

	usersDB := []entities.User{}

	userA := factories.UserFactory{}
	userA.Init(uuid.New())
	usersDB = append(usersDB, entities.User(userA))
	addUserToDb(t, userA)

	userB := factories.UserFactory{}
	userB.Init(uuid.New())
	usersDB = append(usersDB, entities.User(userB))
	addUserToDb(t, userB)

	userC := factories.UserFactory{}
	userC.Init(uuid.New())
	usersDB = append(usersDB, entities.User(userC))
	addUserToDb(t, userC)

	userD := factories.UserFactory{}
	userD.Init(uuid.New())
	usersDB = append(usersDB, entities.User(userD))
	addUserToDb(t, userD)

	// a -> d
	searchAndMakeRequest(t, userA.ID, userD.Username)
	// b -> a
	searchAndMakeRequest(t, userB.ID, userA.Username)
	// c -> a
	searchAndMakeRequest(t, userC.ID, userA.Username)

	friendRequestSentA := getFriendRequestSent(t, userA.ID)
	friendRequestReceivedA := getFriendRequestReceived(t, userA.ID)

	// check req sent for a -> d
	if !friendRequestSentA.Contains(userD.ID) {
		t.Error()
	}

	// check should be 2 (from b and c)
	assert.Equal(t, 2, len(friendRequestReceivedA))

	// check req recv for b -> a
	if !friendRequestReceivedA.Contains(userB.ID) {
		t.Error()
	}

	// check req recv for c -> a
	if !friendRequestReceivedA.Contains(userC.ID) {
		t.Error()
	}

	// add a to d ERROR_ALREADY_REQUESTED_SENT
	searchAndMakeRequest(t, userA.ID, userD.Username)

	// add a to b ERROR_ALREADY_REQUESTED_RECEIVED
	searchAndMakeRequest(t, userA.ID, userB.Username)

	// a reject request for b -> a
	rejectRequest(t, userA.ID, userB.ID)

	// check friend req recv a should be only c left
	friendRequestReceivedA = getFriendRequestReceived(t, userA.ID)
	assert.Equal(t, 1, len(friendRequestReceivedA))
	if !friendRequestReceivedA.Contains(userC.ID) {
		t.Error()
	}

	// a accept request for c -> a
	acceptRequest(t, userA.ID, userC.ID)

	// check friend req recv a should be null
	friendRequestReceivedA = getFriendRequestReceived(t, userA.ID)
	assert.Equal(t, 0, len(friendRequestReceivedA))

	// d acc a
	acceptRequest(t, userD.ID, userA.ID)
	friendRequestSentA = getFriendRequestSent(t, userA.ID)
	assert.Equal(t, 0, len(friendRequestSentA))

	// check friend a
	userAFriendList := getUserFriendList(t, userA.ID)
	// check should be 2 (c, d)
	assert.Equal(t, 2, len(userAFriendList))
	// check c
	if !userAFriendList.Contains(userC.ID) {
		t.Error()
	}
	// check d
	if !userAFriendList.Contains(userD.ID) {
		t.Error()
	}
	// add a to c should be error ERROR_ALREADY_FRIEND
	searchAndMakeRequest(t, userA.ID, userC.Username)

	// add self should be error ERROR_CANNOT_ADD_SELF
	searchAndMakeRequest(t, userA.ID, userA.Username)

	// delete user a b c d in user and friends
	for _, user := range usersDB {
		// delete from user
		db.Where(&entities.User{
			Username: user.Username,
			Email:    user.Email,
			Name:     user.Name,
		}).Delete(&entities.User{})
		// delete from friend
		db.Where(&entities.Friend{
			ID: user.ID,
		}).Delete(&entities.Friend{})
	}
}
