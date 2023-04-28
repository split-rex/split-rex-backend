package controllers

import (
	"fmt"
	"net/http"
	"split-rex-backend/entities"
	"split-rex-backend/entities/requests"
	"split-rex-backend/entities/responses"
	"split-rex-backend/types"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func (con *notificationController) InsertNotif(c echo.Context) error {
	db := con.db
	response := entities.Response[string]{}

	// get request body
	request := requests.NotifRequest{}
	if err := c.Bind(&request); err != nil {
		response.Message = types.ERROR_BAD_REQUEST
		fmt.Println(err)
		return c.JSON(http.StatusBadRequest, response)
	}

	// insert to notif table
	notif := entities.Notification{
		NotificationID: uuid.New(),
		GroupID:        request.GroupID,
		GroupName:      request.GroupName,
		Amount:         request.Amount,
		Name:           request.Name,
		Date:           request.Date,
	}
	if err := db.Create(&notif).Error; err != nil {
		response.Message = types.ERROR_INTERNAL_SERVER
		return c.JSON(http.StatusInternalServerError, response)
	}

	// insert notification id to user table with id = request.user_id
	user := entities.User{}
	condition := entities.User{ID: request.UserID}
	if err := db.Where(&condition).Find(&user).Error; err != nil {
		response.Message = types.ERROR_INTERNAL_SERVER
		return c.JSON(http.StatusInternalServerError, response)
	}
	newNotif := user.Notifications
	if newNotif == nil {
		newNotif = []uuid.UUID{}
	}
	newNotif = append(newNotif, notif.NotificationID)
	if err := db.Model(&user).Updates(entities.User{Notifications: newNotif}).Error; err != nil {
		response.Message = types.ERROR_INTERNAL_SERVER
		return c.JSON(http.StatusInternalServerError, response)
	}

	response.Message = types.SUCCESS
	return c.JSON(http.StatusOK, response)
}

func (con *notificationController) GetNotif(c echo.Context) error {
	db := con.db
	response := entities.Response[responses.NotificationResponse]{}

	id := c.Get("id").(uuid.UUID)
	user := entities.User{}
	if err := db.Find(&user, id).Error; err != nil {
		response.Message = types.ERROR_INTERNAL_SERVER
		return c.JSON(http.StatusInternalServerError, response)
	}

	notifications := []responses.NotificationDetail{}
	for _, notifID := range user.Notifications {
		notif := entities.Notification{}
		if err := db.Find(&notif, notifID).Error; err != nil {
			response.Message = types.ERROR_INTERNAL_SERVER
			return c.JSON(http.StatusInternalServerError, response)
		}
		notifDetail := responses.NotificationDetail{
			NotificationID: notif.NotificationID,
			GroupID:        notif.GroupID,
			GroupName:      notif.GroupName,
			Amount:         notif.Amount,
			Name:           notif.Name,
			Date:           notif.Date,
		}
		notifications = append(notifications, notifDetail)
	}

	response.Message = types.SUCCESS
	response.Data.Notifications = notifications
	return c.JSON(http.StatusOK, response)
}

func (con *notificationController) DeleteNotif(c echo.Context) error {
	db := con.db
	response := entities.Response[string]{}

	id := c.Get("id").(uuid.UUID)
	user := entities.User{}
	if err := db.Find(&user, id).Error; err != nil {
		response.Message = types.ERROR_INTERNAL_SERVER
		return c.JSON(http.StatusInternalServerError, response)
	}

	// delete all notification from notification table
	for _, notifID := range user.Notifications {
		if err := db.Delete(&entities.Notification{}, notifID).Error; err != nil {
			response.Message = types.ERROR_INTERNAL_SERVER
			return c.JSON(http.StatusInternalServerError, response)
		}
	}

	// delete all notification id from user table
	if err := db.Model(&user).Updates(entities.User{Notifications: []uuid.UUID{}}).Error; err != nil {
		response.Message = types.ERROR_INTERNAL_SERVER
		return c.JSON(http.StatusInternalServerError, response)
	}

	response.Message = types.SUCCESS
	return c.JSON(http.StatusOK, response)
}
