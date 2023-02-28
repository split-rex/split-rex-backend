package controllers

import (
	"fmt"
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

	//check if friend_id exist in user table
	friend := entities.User{}
	friend_id, _ := uuid.Parse(friendRequest.Friend_id)
	condition := entities.User{ID: friend_id}
	if err := db.Where(&condition).Find(&friend).Error; err != nil {
		response.Message = types.ERROR_INTERNAL_SERVER
		return c.JSON(http.StatusInternalServerError, response)
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

	//if userExist true
	if userExist {
		// insert friend_id to friend table

		userFriend.Req_sent = append(userFriend.Req_sent, friend_id)
		if err := db.Save(&userFriend).Error; err != nil {
			response.Message = types.ERROR_INTERNAL_SERVER
			return c.JSON(http.StatusInternalServerError, response)
		}
	} else {
		//insert user_id and friend_id to friend table
		userFriend.ID = user_id
		userFriend.Req_sent = append(userFriend.Req_sent, friend_id)
		if err := db.Create(&userFriend).Error; err != nil {
			response.Message = types.ERROR_INTERNAL_SERVER
			return c.JSON(http.StatusInternalServerError, response)
		}
	}

	//if friendExist true
	if friendExist {
		// insert user_id to friend table
		friendFriend.Req_received = append(friendFriend.Req_received, user_id)
		if err := db.Save(&friendFriend).Error; err != nil {
			response.Message = types.ERROR_INTERNAL_SERVER
			return c.JSON(http.StatusInternalServerError, response)
		}
	} else {
		//insert friend_id and user_id to friend table
		friendFriend.ID = friend_id
		friendFriend.Req_received = append(friendFriend.Req_received, user_id)
		if err := db.Create(&friendFriend).Error; err != nil {
			response.Message = types.ERROR_INTERNAL_SERVER
			return c.JSON(http.StatusInternalServerError, response)
		}
	}

	response.Message = types.SUCCESS
	return c.JSON(http.StatusOK, response)
}

func FriendRequestSent(c echo.Context) error {
	user_id := c.Get("id").(uuid.UUID)
	fmt.Println(user_id)
	db := database.DB.GetConnection()
	response := entities.Response[[]entities.User]{}

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

	if userExist {
		// get username and full name where user_id in Req_sent
		// TODO: ini kayanya harus loop
		users := []entities.User{}
		if err := db.Where("id IN ?", userFriend.Req_sent).Select("id", "username", "name").Find(&users).Error; err != nil {
			response.Message = types.ERROR_INTERNAL_SERVER
			return c.JSON(http.StatusInternalServerError, response)
		}

		response.Message = types.SUCCESS
		response.Data = users
		return c.JSON(http.StatusOK, response)

		// return json of id, username, and full name
		// newResponse := entities.Response[[]byte]{}
		// usersJson, _ := json.MarshalIndent(users, "", " ")
		// fmt.Println(users)
		// fmt.Println(usersJson)

		// newResponse.Message = types.SUCCESS
		// newResponse.Data = usersJson

	} else {
		response.Message = types.SUCCESS
		response.Data = []entities.User{}
		return c.JSON(http.StatusOK, response)

		// newResponse := entities.Response[[]byte]{}
		// newResponse.Message = types.SUCCESS
		// newResponse.Data = []byte("[]")
		// fmt.Println("kosong")
		// return c.JSON(http.StatusOK, newResponse)
	}
}
