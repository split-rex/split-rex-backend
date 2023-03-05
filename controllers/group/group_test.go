package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"split-rex-backend/entities/requests"
	"split-rex-backend/types"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestUserCreateGroup(t *testing.T) {
	date := time.Date(2023, 3, 3, 0, 0, 0, 0, time.UTC)

	group := &requests.UserCreateGroupRequest{
		Name:      "Group Testing 12345",
		MemberID:  types.ArrayOfUUID{},
		StartDate: date,
		EndDate:   date,
	}
	body, _ := json.Marshal(group)

	res, err := http.Post("http://localhost:8080/userCreateGroup",
		"application/json", bytes.NewBuffer(body))
	if err != nil {
		t.Error(err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Error("Expected status code 200, got ", res.StatusCode)
	}
}

func TestEditGroupInfo(t *testing.T) {
	date := time.Date(2023, 3, 3, 0, 0, 0, 0, time.UTC)

	group := &requests.UserCreateGroupRequest{
		Name:      "Group Testing 12345",
		MemberID:  types.ArrayOfUUID{},
		StartDate: date,
		EndDate:   date,
	}
	body, _ := json.Marshal(group)

	res, err := http.Post("http://localhost:8080/userCreateGroup",
		"application/json", bytes.NewBuffer(body))
	if err != nil {
		t.Error(err)
	}
	defer res.Body.Close()

	// create new interface of response which consist of message and groupID

	if res.StatusCode != http.StatusOK {
		t.Error("Expected status code 200 while creating group, got ", res.StatusCode)
	}

	var response struct {
		Message string
		GroupID string
	}
	json.NewDecoder(res.Body).Decode(&response)

	groupID, _ := uuid.Parse(response.GroupID)

	editGroup := &requests.EditGroupInfoRequest{
		GroupID:   groupID,
		Name:      "Group Testing New",
		StartDate: date,
		EndDate:   date,
	}
	body, _ = json.Marshal(editGroup)

	res, err = http.Post("http://localhost:8080/editGroupInfo",
		"application/json", bytes.NewBuffer(body))
	if err != nil {
		t.Error(err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Error("Expected status code 200, got ", res.StatusCode)
	}
}
