package controllers

import (
	"net/http"
	"split-rex-backend/entities"
	"split-rex-backend/entities/responses"
	"split-rex-backend/types"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func (con *groupController) GroupOwed(c echo.Context) error {
	db := con.db
	response := entities.Response[responses.GroupOwedResponse]{}

	id := c.Get("id").(uuid.UUID)
	user := entities.User{}
	if err := db.Find(&user, id).Error; err != nil {
		response.Message = types.ERROR_INTERNAL_SERVER
		return c.JSON(http.StatusInternalServerError, response)
	}

	// iterate through user's groups to get groups details
	totalOwedGlobal := 0.0
	groups := []responses.UserGroupResponse{}
	for _, groupID := range user.Groups {
		totalOwed := 0.0
		totalExpense := 0.0
		group := entities.Group{}
		condition := entities.Group{GroupID: groupID}
		if err := db.Where(&condition).Find(&group).Error; err != nil {
			response.Message = types.ERROR_INTERNAL_SERVER
			return c.JSON(http.StatusInternalServerError, response)
		}

		// then for each group, search for transactions existed in group id
		transactions := []entities.Transaction{}
		conditionTransaction := entities.Transaction{GroupID: groupID}
		if err := db.Where(&conditionTransaction).Find(&transactions).Error; err != nil {
			response.Message = types.ERROR_INTERNAL_SERVER
			return c.JSON(http.StatusInternalServerError, response)
		}

		// compute totalExpense from transactions
		for _, transaction := range transactions {
			totalExpense = totalExpense + transaction.Total
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
			continue
		}
		totalOwedGlobal = totalOwedGlobal + totalOwed

		groups = append(groups,
			responses.UserGroupResponse{
				GroupID:      group.GroupID,
				Name:         group.Name,
				MemberID:     group.MemberID,
				StartDate:    group.StartDate,
				EndDate:      group.EndDate,
				Type:         types.TYPE_GROUP_OWED,
				TotalUnpaid:  totalOwed,
				TotalExpense: totalExpense})
	}

	// then return all
	response.Message = types.SUCCESS
	response.Data.TotalOwed = totalOwedGlobal
	response.Data.ListGroup = groups

	return c.JSON(http.StatusOK, response)
}
