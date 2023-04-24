package controllers

import (
	"math"
	"net/http"
	"split-rex-backend/entities"
	"split-rex-backend/entities/responses"
	"split-rex-backend/types"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func (con *statisticController) OwedLentPercentage(c echo.Context) error {
	db := con.db
	response := entities.Response[responses.PercentageResponse]{}

	id := c.Get("id").(uuid.UUID)
	user := entities.User{}
	if err := db.Find(&user, id).Error; err != nil {
		response.Message = types.ERROR_INTERNAL_SERVER
		return c.JSON(http.StatusInternalServerError, response)
	}

	// iterate through user's groups to get groups details
	totalOwedGlobal := 0.0
	totalLentGlobal := 0.0
	for _, groupID := range user.Groups {
		totalOwed := 0.0
		group := entities.Group{}
		condition := entities.Group{GroupID: groupID}
		if err := db.Where(&condition).Find(&group).Error; err != nil {
			response.Message = types.ERROR_INTERNAL_SERVER
			return c.JSON(http.StatusInternalServerError, response)
		}

		// then for each group, search for payments existed in group id
		payments := []entities.Payment{}
		conditionPayment := entities.Payment{GroupID: groupID, UserID1: id}
		if err := db.Where(&conditionPayment).Find(&payments).Error; err != nil {
			response.Message = types.ERROR_INTERNAL_SERVER
			return c.JSON(http.StatusInternalServerError, response)
		}

		// compute totalOwed from payments
		for _, payment := range payments {
			totalOwed = totalOwed + payment.TotalUnpaid
		}

		// if totalOwed is negative then not in groupOwed
		if totalOwed <= 0 {
			totalLentGlobal = totalLentGlobal - totalOwed
		} else {
			totalOwedGlobal = totalOwedGlobal + totalOwed
		}
	}

	// then return all
	response.Message = types.SUCCESS
	if totalOwedGlobal == 0 && totalLentGlobal == 0 {
		response.Data.OwedPercentage = 50
		response.Data.LentPercentage = 50
	} else {
		owedPercentage := int(math.Round(totalOwedGlobal * 100 / (totalOwedGlobal + totalLentGlobal)))
		lentPercentage := 100 - owedPercentage
		response.Data.OwedPercentage = owedPercentage
		response.Data.LentPercentage = lentPercentage
	}

	return c.JSON(http.StatusOK, response)
}

func (con *statisticController) PaymentMutation(c echo.Context) error {
	db := con.db
	response := entities.Response[responses.MutationResponse]{}
	totalPaid := 0.0
	totalReceived := 0.0

	id := c.Get("id").(uuid.UUID)
	user := entities.User{}
	if err := db.Find(&user, id).Error; err != nil {
		response.Message = types.ERROR_INTERNAL_SERVER
		return c.JSON(http.StatusInternalServerError, response)
	}

	startDate, _ := time.Parse("2006-01-02", c.QueryParam("start_date"))
	endDate, _ := time.Parse("2006-01-02", c.QueryParam("end_date"))

	//iterate through activity table where user_id = id and type = payment
	activities := []entities.Activity{}
	condition := entities.Activity{UserID: id, ActivityType: "PAYMENT"}
	if err := db.Where(&condition).Find(&activities).Error; err != nil {
		response.Message = types.ERROR_INTERNAL_SERVER
		return c.JSON(http.StatusInternalServerError, response)
	}

	// filter activities by date
	filteredActivities := []entities.Activity{}
	for _, activity := range activities {
		if activity.Date.After(startDate) && activity.Date.Before(endDate) {
			filteredActivities = append(filteredActivities, activity)
		}
	}

	// check if payment activites status = confirmed
	mutationDetail := []responses.MutationDetail{}
	for _, activity := range filteredActivities {
		paymentActivity := entities.PaymentActivity{}
		if err := db.Find(&paymentActivity, activity.DetailID).Error; err != nil {
			response.Message = types.ERROR_INTERNAL_SERVER
			return c.JSON(http.StatusInternalServerError, response)
		}
		if paymentActivity.Status == "CONFIRMED" {
			// get user color
			userDetail := entities.User{}
			condition := entities.User{Name: paymentActivity.Name}
			if err := db.Where(&condition).Find(&userDetail).Error; err != nil {
				response.Message = types.ERROR_INTERNAL_SERVER
				return c.JSON(http.StatusInternalServerError, response)
			}
			mutationDetail = append(mutationDetail, responses.MutationDetail{
				Name:         paymentActivity.Name,
				Color:        userDetail.Color,
				MutationType: "PAID",
				Amount:       paymentActivity.Amount,
			})
			totalPaid = totalPaid + paymentActivity.Amount
		}
	}

	// iterate through payment activity where name = user.name and status = confirmed
	paymentActivities := []entities.PaymentActivity{}
	paymentCondition := entities.PaymentActivity{Name: user.Name, Status: "CONFIRMED"}
	if err := db.Where(&paymentCondition).Find(&paymentActivities).Error; err != nil {
		response.Message = types.ERROR_INTERNAL_SERVER
		return c.JSON(http.StatusInternalServerError, response)
	}

	// iterate through activity where detail_id = payment_activity.id
	for _, paymentActivity := range paymentActivities {
		activity := entities.Activity{}
		condition := entities.Activity{DetailID: paymentActivity.PaymentActivityID}
		if err := db.Where(&condition).Find(&activity).Error; err != nil {
			response.Message = types.ERROR_INTERNAL_SERVER
			return c.JSON(http.StatusInternalServerError, response)
		}
		// check if activity date is between start_date and end_date
		if activity.Date.After(startDate) && activity.Date.Before(endDate) {
			// get user name and color
			userDetail := entities.User{}
			condition := entities.User{ID: activity.UserID}
			if err := db.Where(&condition).Find(&userDetail).Error; err != nil {
				response.Message = types.ERROR_INTERNAL_SERVER
				return c.JSON(http.StatusInternalServerError, response)
			}
			mutationDetail = append(mutationDetail, responses.MutationDetail{
				Name:         userDetail.Name,
				Color:        userDetail.Color,
				MutationType: "RECEIVED",
				Amount:       paymentActivity.Amount,
			})
			totalReceived = totalReceived + paymentActivity.Amount
		}
	}

	// return
	response.Message = types.SUCCESS
	response.Data.ListMutation = mutationDetail
	response.Data.TotalPaid = totalPaid
	response.Data.TotalReceived = totalReceived
	return c.JSON(http.StatusOK, response)
}
