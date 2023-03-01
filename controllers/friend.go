package controllers

import (
	"fmt"
	"net/http"
	"split-rex-backend/configs/database"
	"split-rex-backend/entities"
	"split-rex-backend/entities/requests"
	"split-rex-backend/entities/responses"
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
	response := entities.Response[[]responses.FriendResponse]{}

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
		users := []responses.FriendResponse{}
		for _, id := range userFriend.Req_sent {
			user := entities.User{}
			friend := responses.FriendResponse{}
			condition := entities.User{ID: id}
			db.Where(&condition).Select("id", "username", "name").Find(&user)
			friend.User_id = user.ID.String()
			friend.Username = user.Username
			friend.Fullname = user.Name
			users = append(users, friend)
		}
		response.Message = types.SUCCESS
		if len(users) > 0 {
			response.Data = users
		} else {
			response.Message = types.DATA_NOT_FOUND
		}
		return c.JSON(http.StatusOK, response)

	} else {
		response.Message = types.DATA_NOT_FOUND
		response.Data = []responses.FriendResponse{}
		return c.JSON(http.StatusOK, response)
	}
}

func FriendRequestReceived(c echo.Context) error {
	user_id := c.Get("id").(uuid.UUID)
	fmt.Println(user_id)
	db := database.DB.GetConnection()
	response := entities.Response[[]responses.FriendResponse]{}

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
		users := []responses.FriendResponse{}
		for _, id := range userFriend.Req_received {
			user := entities.User{}
			friend := responses.FriendResponse{}
			condition := entities.User{ID: id}
			db.Where(&condition).Select("id", "username", "name").Find(&user)
			friend.User_id = user.ID.String()
			friend.Username = user.Username
			friend.Fullname = user.Name
			users = append(users, friend)
		}
		response.Message = types.SUCCESS
		if len(users) > 0 {
			response.Data = users
		} else {
			response.Message = types.DATA_NOT_FOUND
		}
		return c.JSON(http.StatusOK, response)

	} else {
		response.Message = types.DATA_NOT_FOUND
		response.Data = []responses.FriendResponse{}
		return c.JSON(http.StatusOK, response)
	}
}

func SearchUser(c echo.Context) error {
	username := c.QueryParam("username")
	db := database.DB.GetConnection()
	response := entities.Response[responses.FriendResponse]{}

	//check if username exist in user table
	user := []entities.User{}
	conditionUser := entities.User{Username: username}
	if err := db.Where(&conditionUser).Find(&user).Error; err != nil {
		response.Message = types.ERROR_INTERNAL_SERVER
		return c.JSON(http.StatusInternalServerError, response)
	}

	if len(user) > 0 {
		response.Message = types.SUCCESS
		friend := responses.FriendResponse{}
		friend.User_id = user[0].ID.String()
		friend.Username = user[0].Username
		friend.Fullname = user[0].Name
		response.Data = friend
	} else {
		response.Message = types.DATA_NOT_FOUND
	}
	return c.JSON(http.StatusOK, response)
}
