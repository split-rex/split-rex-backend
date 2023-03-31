package controllers

import (
	"net/http"
	"split-rex-backend/entities"
	"split-rex-backend/entities/responses"
	"split-rex-backend/types"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func (h *activityController) GetUserActivity(c echo.Context) error {
	db := h.db
	response := entities.Response[[]interface{}]{}

	//get user_id from token
	userID := c.Get("id").(uuid.UUID)

	// get activity from activity table where UserID = userID
	activities := []entities.Activity{}
	conditionActivity := entities.Activity{UserID: userID}
	if err := db.Where(&conditionActivity).Find(&activities).Error; err != nil {
		response.Message = types.ERROR_INTERNAL_SERVER
		return c.JSON(http.StatusInternalServerError, response)
	}

	// move activities to response data
	activityResponse := []interface{}{}
	for _, activity := range activities {
		if activity.ActivityType == "REMINDER" {
			reminder := entities.ReminderActivity{}
			if err := db.Find(&reminder, activity.DetailID).Error; err != nil {
				response.Message = types.ERROR_INTERNAL_SERVER
				return c.JSON(http.StatusInternalServerError, response)
			}
			reminderResponse := responses.ActivityResponse[responses.ReminderActivityResponse]{
				ActivityID:   activity.ActivityID,
				ActivityType: activity.ActivityType,
				Date:         activity.Date,
				RedirectID:   activity.RedirectID,
				Detail: responses.ReminderActivityResponse{
					ReminderActivityID: reminder.ReminderActivityID,
					Name:               reminder.Name,
					GroupName:          reminder.GroupName,
				},
			}
			activityResponse = append(activityResponse, reminderResponse)

		} else if activity.ActivityType == "PAYMENT" {
			payment := entities.PaymentActivity{}
			if err := db.Find(&payment, activity.DetailID).Error; err != nil {
				response.Message = types.ERROR_INTERNAL_SERVER
				return c.JSON(http.StatusInternalServerError, response)
			}
			paymentResponse := responses.ActivityResponse[responses.PaymentActivityResponse]{
				ActivityID:   activity.ActivityID,
				ActivityType: activity.ActivityType,
				Date:         activity.Date,
				RedirectID:   activity.RedirectID,
				Detail: responses.PaymentActivityResponse{
					PaymentActivityID: payment.PaymentActivityID,
					Name:              payment.Name,
					Status:            payment.Status,
					Amount:            payment.Amount,
					GroupName:         payment.GroupName,
				},
			}
			activityResponse = append(activityResponse, paymentResponse)

		} else if activity.ActivityType == "TRANSACTION" {
			transaction := entities.TransactionActivity{}
			if err := db.Find(&transaction, activity.DetailID).Error; err != nil {
				response.Message = types.ERROR_INTERNAL_SERVER
				return c.JSON(http.StatusInternalServerError, response)
			}
			transactionResponse := responses.ActivityResponse[responses.TransactionActivityResponse]{
				ActivityID:   activity.ActivityID,
				ActivityType: activity.ActivityType,
				Date:         activity.Date,
				RedirectID:   activity.RedirectID,
				Detail: responses.TransactionActivityResponse{
					TransactionActivityID: transaction.TransactionActivityID,
					Name:                  transaction.Name,
					GroupName:             transaction.GroupName,
				},
			}
			activityResponse = append(activityResponse, transactionResponse)
		} else {
			continue
		}
	}

	response.Message = types.SUCCESS
	response.Data = activityResponse
	return c.JSON(http.StatusOK, response)
}

func (h *activityController) GetGroupActivity(c echo.Context) error {
	db := h.db
	response := entities.Response[[]responses.GroupActivityResponse]{}

	//get user_id from token and group_id from param
	userID := c.Get("id").(uuid.UUID)
	groupID, _ := uuid.Parse(c.QueryParam("group_id"))

	// get groupActivity where groupID = groupID and userID1 = userID
	groupActivities := []entities.GroupActivity{}
	conditionGroupActivity := entities.GroupActivity{GroupID: groupID, UserID1: userID}
	if err := db.Where(&conditionGroupActivity).Find(&groupActivities).Error; err != nil {
		response.Message = types.ERROR_INTERNAL_SERVER
		return c.JSON(http.StatusInternalServerError, response)
	}

	// get groupActivity where groupID = groupID and userID2 = userID
	groupActivities2 := []entities.GroupActivity{}
	conditionGroupActivity2 := entities.GroupActivity{GroupID: groupID, UserID2: userID}
	if err := db.Where(&conditionGroupActivity2).Find(&groupActivities2).Error; err != nil {
		response.Message = types.ERROR_INTERNAL_SERVER
		return c.JSON(http.StatusInternalServerError, response)
	}

	// concat groupActivities and groupActivities2
	groupActivities = append(groupActivities, groupActivities2...)

	// get name of user1 and user2
	activityResponse := []responses.GroupActivityResponse{}
	for _, groupActivity := range groupActivities {
		user1 := entities.User{}
		if err := db.Find(&user1, groupActivity.UserID1).Error; err != nil {
			response.Message = types.ERROR_INTERNAL_SERVER
			return c.JSON(http.StatusInternalServerError, response)
		}
		user2 := entities.User{}
		if err := db.Find(&user2, groupActivity.UserID2).Error; err != nil {
			response.Message = types.ERROR_INTERNAL_SERVER
			return c.JSON(http.StatusInternalServerError, response)
		}

		activityResponse = append(activityResponse, responses.GroupActivityResponse{
			ActivityID: groupActivity.ActivityID,
			Date:       groupActivity.Date,
			Name1:      user1.Name,
			Name2:      user2.Name,
			Amount:     groupActivity.Amount,
		})
	}

	response.Message = types.SUCCESS
	response.Data = activityResponse
	return c.JSON(http.StatusOK, response)
}
