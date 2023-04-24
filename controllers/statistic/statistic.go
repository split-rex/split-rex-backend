package controllers

import (
	"math"
	"net/http"
	"split-rex-backend/entities"
	"split-rex-backend/entities/responses"
	"split-rex-backend/types"

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
