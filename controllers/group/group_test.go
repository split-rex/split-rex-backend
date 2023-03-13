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
	testGroupController = NewGroupController(database.DBTesting.GetConnection())
)

func TestUserCreateGroup(t *testing.T) {
	db := database.DBTesting.GetConnection()

	e := echo.New()
	groupJson := `{
		"name": "New Group Yeay",
		"member_id": ["6251ac85-e43d-4b88-8779-588099df5008","183e04d7-c653-4c7d-aa66-3d751d4d7358"],
		"start_date": "2023-03-01T17:19:20.968831+07:00",
		"end_date" : "2023-03-01T19:19:20.968831+07:00"
	}`

	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(groupJson))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	groupID := responses.TestResponse[string]{}
	if assert.NoError(t, testGroupController.UserCreateGroup(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)

		if err := json.Unmarshal(rec.Body.Bytes(), &groupID); err != nil {
			panic(err)
		}
	}

	// delete created group
	group := entities.Group{
		GroupID: uuid.MustParse(groupID.Data),
	}
	db.Where(&group).Delete(&group)
}

// func TestUserCreateGroup(t *testing.T) {
// 	date := time.Date(2023, 3, 3, 0, 0, 0, 0, time.UTC)

// 	group := &requests.UserCreateGroupRequest{
// 		Name:      "Group Testing 12345",
// 		MemberID:  types.ArrayOfUUID{},
// 		StartDate: date,
// 		EndDate:   date,
// 	}
// 	body, _ := json.Marshal(group)

// 	res, err := http.Post("http://localhost:8080/userCreateGroup",
// 		"application/json", bytes.NewBuffer(body))
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	defer res.Body.Close()

// 	if res.StatusCode != http.StatusOK {
// 		t.Error("Expected status code 200, got ", res.StatusCode)
// 	}
// }

// func TestEditGroupInfo(t *testing.T) {
// 	date := time.Date(2023, 3, 3, 0, 0, 0, 0, time.UTC)

// 	group := &requests.UserCreateGroupRequest{
// 		Name:      "Group Testing 12345",
// 		MemberID:  types.ArrayOfUUID{},
// 		StartDate: date,
// 		EndDate:   date,
// 	}
// 	body, _ := json.Marshal(group)

// 	res, err := http.Post("http://localhost:8080/userCreateGroup",
// 		"application/json", bytes.NewBuffer(body))
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	defer res.Body.Close()

// 	// create new interface of response which consist of message and groupID

// 	if res.StatusCode != http.StatusOK {
// 		t.Error("Expected status code 200 while creating group, got ", res.StatusCode)
// 	}

// 	var response struct {
// 		Message string
// 		GroupID string
// 	}
// 	json.NewDecoder(res.Body).Decode(&response)

// 	groupID, _ := uuid.Parse(response.GroupID)

// 	editGroup := &requests.EditGroupInfoRequest{
// 		GroupID:   groupID,
// 		Name:      "Group Testing New",
// 		StartDate: date,
// 		EndDate:   date,
// 	}
// 	body, _ = json.Marshal(editGroup)

// 	res, err = http.Post("http://localhost:8080/editGroupInfo",
// 		"application/json", bytes.NewBuffer(body))
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	defer res.Body.Close()

// 	if res.StatusCode != http.StatusOK {
// 		t.Error("Expected status code 200, got ", res.StatusCode)
// 	}
// }
