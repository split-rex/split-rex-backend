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
		response.Message = types.ERROR_BAD_REQUEST
		return c.JSON(http.StatusInternalServerError, response)
	}
	response.Message = types.SUCCESS
	return c.JSON(http.StatusOK, response)
}

func FriendRequestSent(c echo.Context) error {
	user_id := c.Get("id").(uuid.UUID)
	// fmt.Println(user_id)
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
			if err := db.Where(&condition).Select("id", "username", "name").Find(&user).Error; err != nil {
				response.Message = types.ERROR_INTERNAL_SERVER
				return c.JSON(http.StatusInternalServerError, response)
			}
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
			if err := db.Where(&condition).Select("id", "username", "name").Find(&user).Error; err != nil {
				response.Message = types.ERROR_INTERNAL_SERVER
				return c.JSON(http.StatusInternalServerError, response)
			}
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

func UserFriendList(c echo.Context) error {
	db := database.DB.GetConnection()

	// get user id from context, cast to uuid
	user_id := c.Get("id").(uuid.UUID)

	// returning array of friendResponse struct
	response := entities.Response[[]responses.FriendResponse]{}

	//check if user_id exist in friend table
	userInFriend := entities.Friend{}
	conditionSearchUser := entities.Friend{ID: user_id}
	if err := db.Where(&conditionSearchUser).Find(&userInFriend).Error; err != nil {
		response.Message = types.ERROR_INTERNAL_SERVER
		return c.JSON(http.StatusInternalServerError, response)
	}

	//check if userFriend empty
	if userInFriend.ID == uuid.Nil {
		response.Message = types.DATA_NOT_FOUND
		response.Data = []responses.FriendResponse{}
		return c.JSON(http.StatusOK, response)
	}

	// get username and full name where friend_id (Friend table) exist for user
	friends := []responses.FriendResponse{}
	for _, id := range userInFriend.Friend_id {
		user := entities.User{}
		friend := responses.FriendResponse{}
		condition := entities.User{ID: id}
		// get id, username, and name from user table
		if err := db.Where(&condition).Select("id", "username", "name").Find(&user).Error; err != nil {
			response.Message = types.ERROR_INTERNAL_SERVER
			return c.JSON(http.StatusInternalServerError, response)
		}
		friend.User_id = user.ID.String()
		friend.Username = user.Username
		friend.Fullname = user.Name
		friends = append(friends, friend)
	}

	response.Data = friends

	// if no friends, data not found but status success
	if len(friends) <= 0 {
		response.Message = types.DATA_NOT_FOUND
		return c.JSON(http.StatusOK, response)
	}

	response.Message = types.SUCCESS
	return c.JSON(http.StatusOK, response)

}

func AcceptRequest(c echo.Context) error {
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
	// user_new_req_received := []uuid.UUID{}
	fmt.Println(userInFriend.Req_received)
	for i, id := range userInFriend.Req_received {
		fmt.Println(id.String())
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
	// update userInFriend value
	req_id, _ := uuid.Parse(requester_id)
	userInFriend.Friend_id = append(userInFriend.Friend_id, req_id)

	// update user on database
	if err := tx.Save(&userInFriend).Error; err != nil {
		tx.Rollback()
		response.Message = types.ERROR_BAD_REQUEST
		return c.JSON(http.StatusInternalServerError, response)
	}

	// requester
	requesterInFriend := entities.Friend{}
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
	requesterInFriend.Friend_id = append(requesterInFriend.Friend_id, user_id)

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
