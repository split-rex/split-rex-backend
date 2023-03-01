package controllers

import (
	"net/http"
	"split-rex-backend/configs/database"
	"split-rex-backend/entities"
	"split-rex-backend/entities/requests"
	"split-rex-backend/types"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func RejectRequest(c echo.Context) error {
	db := database.DB.GetConnection()

	response := entities.Response[string]{}

	// get user id from context, cast to uuid
	user_id := c.Get("id").(uuid.UUID)

	// get request_id from body
	friendRequest := requests.FriendRequest{}
	if err := c.Bind(&friendRequest); err != nil {
		response.Message = types.ERROR_BAD_REQUEST
		return c.JSON(http.StatusBadRequest, response)
	}
	requester_id := friendRequest.Friend_id

	// check if user in friend table
	userInFriend := entities.Friend{}
	conditionFriend := entities.Friend{ID: user_id}
	if err := db.Where(&conditionFriend).Find(&userInFriend).Error; err != nil {
		response.Message = types.ERROR_BAD_REQUEST
		return c.JSON(http.StatusInternalServerError, response)
	}

	// error if no user in friend table
	if userInFriend.ID == uuid.Nil {
		response.Message = types.ERROR_BAD_REQUEST
		return c.JSON(http.StatusInternalServerError, response)
	}

	// check if requester_id exist in req received[] user's friend table
	found := false
	for i, id := range userInFriend.Req_received {
		if id.String() == requester_id {
			found = true
			userInFriend.Req_received = append(userInFriend.Req_received[:i], userInFriend.Req_received[i+1:]...)
			break
		}
	}

	// error if no friend request from userFriend.ID
	if !found {
		response.Message = types.ERROR_BAD_REQUEST
		return c.JSON(http.StatusInternalServerError, response)
	}

	// TRANSACTION
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// user
	// update user
	if err := tx.Save(&userInFriend).Error; err != nil {
		tx.Rollback()
		response.Message = types.ERROR_BAD_REQUEST
		return c.JSON(http.StatusInternalServerError, response)
	}

	// requester
	requesterInFriend := entities.Friend{}
	req_id, _ := uuid.Parse(requester_id)
	conditionRequester := entities.Friend{ID: req_id}
	// check requester in table friends
	if err := tx.Where(&conditionRequester).Find(&requesterInFriend).Error; err != nil {
		tx.Rollback()
		response.Message = types.ERROR_BAD_REQUEST
		return c.JSON(http.StatusInternalServerError, response)
	}

	// update requesterInFriend value
	for i, id := range requesterInFriend.Req_sent {
		if id.String() == user_id.String() {
			//delete element
			requesterInFriend.Req_sent = append(requesterInFriend.Req_sent[:i], requesterInFriend.Req_sent[i+1:]...)
			break
		}
	}

	// update requester on database
	if err := tx.Save(&requesterInFriend).Error; err != nil {
		tx.Rollback()
		response.Message = types.ERROR_BAD_REQUEST
		return c.JSON(http.StatusInternalServerError, response)
	}

	// commit transaction
	if err := tx.Commit().Error; err != nil {
		response.Message = types.ERROR_BAD_REQUEST
		return c.JSON(http.StatusInternalServerError, response)
	}

	response.Message = types.SUCCESS
	return c.JSON(http.StatusOK, response)
}
