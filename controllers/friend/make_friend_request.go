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

func MakeFriendRequest(c echo.Context) error {
	db := database.DB.GetConnection()
	// config := configs.Config.GetMetadata()
	response := entities.Response[string]{}

	friendRequest := requests.FriendRequest{}
	if err := c.Bind(&friendRequest); err != nil {
		response.Message = types.ERROR_BAD_REQUEST
		return c.JSON(http.StatusBadRequest, response)
	}

	user_id := c.Get("id").(uuid.UUID)
	friend_id, _ := uuid.Parse(friendRequest.Friend_id)

	if user_id == friend_id {
		response.Message = types.ERROR_BAD_REQUEST + ": user_id and friend_id cannot be the same"
		return c.JSON(http.StatusBadRequest, response)
	}

	//check if friend_id exist in user table
	friend := entities.User{}
	condition := entities.User{ID: friend_id}
	if err := db.Where(&condition).Find(&friend).Error; err != nil {
		response.Message = types.ERROR_INTERNAL_SERVER
		return c.JSON(http.StatusInternalServerError, response)
	}

	//check if friend empty
	if friend.ID == uuid.Nil {
		response.Message = types.ERROR_USER_NOT_FOUND
		return c.JSON(http.StatusNotFound, response)
	}

	//check if user_id exist in friend table
	userFriend := entities.Friend{}
	userExist := true
	conditionFriend := entities.Friend{ID: user_id}
	if err := db.Where(&conditionFriend).Find(&userFriend).Error; err != nil {
		response.Message = types.ERROR_INTERNAL_SERVER
		return c.JSON(http.StatusInternalServerError, response)
	}
	//check if userFriend empty
	if userFriend.ID == uuid.Nil {
		userExist = false
	}

	//check if friend_id exist in friend table
	friendFriend := entities.Friend{}
	friendExist := true
	conditionFriend = entities.Friend{ID: friend_id}
	if err := db.Where(&conditionFriend).Find(&friendFriend).Error; err != nil {
		response.Message = types.ERROR_INTERNAL_SERVER
		return c.JSON(http.StatusInternalServerError, response)
	}
	if friendFriend.ID == uuid.Nil {
		friendExist = false
	}

	// TRANSACTION
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	//if userExist true
	if userExist {
		// check if friend_id already in Friend_id
		for _, friend := range userFriend.Friend_id {
			if friend == friend_id {
				response.Message = types.ERROR_ALREADY_FRIEND
				return c.JSON(http.StatusConflict, response)
			}
		}
		// check if friend_id already in Req_received
		for _, friend := range userFriend.Req_received {
			if friend == friend_id {
				response.Message = types.ERROR_ALREADY_REQUESTED
				return c.JSON(http.StatusConflict, response)
			}
		}
		// check if friend_id already in Req_sent
		for _, friend := range userFriend.Req_sent {
			if friend == friend_id {
				response.Message = types.ERROR_ALREADY_REQUESTED
				return c.JSON(http.StatusConflict, response)
			}
		}
		// insert friend_id to friend table
		userFriend.Req_sent = append(userFriend.Req_sent, friend_id)
		if err := tx.Save(&userFriend).Error; err != nil {
			tx.Rollback()
			response.Message = types.ERROR_INTERNAL_SERVER
			return c.JSON(http.StatusInternalServerError, response)
		}
	} else {
		//insert user_id and friend_id to friend table
		userFriend.ID = user_id
		userFriend.Req_sent = append(userFriend.Req_sent, friend_id)
		if err := tx.Create(&userFriend).Error; err != nil {
			tx.Rollback()
			response.Message = types.ERROR_INTERNAL_SERVER
			return c.JSON(http.StatusInternalServerError, response)
		}
	}

	//if friendExist true
	if friendExist {
		// insert user_id to friend table
		friendFriend.Req_received = append(friendFriend.Req_received, user_id)
		if err := tx.Save(&friendFriend).Error; err != nil {
			tx.Rollback()
			response.Message = types.ERROR_INTERNAL_SERVER
			return c.JSON(http.StatusInternalServerError, response)
		}
	} else {
		//insert friend_id and user_id to friend table
		friendFriend.ID = friend_id
		friendFriend.Req_received = append(friendFriend.Req_received, user_id)
		if err := tx.Create(&friendFriend).Error; err != nil {
			tx.Rollback()
			response.Message = types.ERROR_INTERNAL_SERVER
			return c.JSON(http.StatusInternalServerError, response)
		}
	}

	// commit transaction
	if err := tx.Commit().Error; err != nil {
		response.Message = types.ERROR_INTERNAL_SERVER
		return c.JSON(http.StatusInternalServerError, response)
	}
	response.Message = types.SUCCESS
	return c.JSON(http.StatusOK, response)
}
