package controllers

import (
	"math"
	"net/http"
	"sort"
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
	response.Data.TotalLent = totalLentGlobal
	response.Data.TotalOwed = totalOwedGlobal

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

func (con *statisticController) SpendingBuddies(c echo.Context) error {
	db := con.db
	response := entities.Response[responses.BuddyResponse]{}

	id := c.Get("id").(uuid.UUID)
	user := entities.User{}
	if err := db.Find(&user, id).Error; err != nil {
		response.Message = types.ERROR_INTERNAL_SERVER
		return c.JSON(http.StatusInternalServerError, response)
	}

	// get all friends of user
	friends := entities.Friend{}
	if err := db.Find(&friends, id).Error; err != nil {
		response.Message = types.ERROR_INTERNAL_SERVER
		return c.JSON(http.StatusInternalServerError, response)
	}

	// create map from friend_id to number of common transaction
	commonTransaction := map[uuid.UUID]int{}

	// iterate through friends
	for _, friend := range friends.Friend_id {
		commonTransaction[friend] = 0
		// get friend group
		friendGroup := entities.User{}
		if err := db.Find(&friendGroup, friend).Error; err != nil {
			response.Message = types.ERROR_INTERNAL_SERVER
			return c.JSON(http.StatusInternalServerError, response)
		}

		// get mutual group with user
		mutualGroup := []uuid.UUID{}
		for _, group := range user.Groups {
			for _, friendGroup := range friendGroup.Groups {
				if group == friendGroup {
					mutualGroup = append(mutualGroup, group)
				}
			}
		}

		// get all transaction dengan group_id = mutual group
		transactions := []entities.Transaction{}
		for _, group := range mutualGroup {
			transaction := []entities.Transaction{}
			condition := entities.Transaction{GroupID: group}
			if err := db.Where(&condition).Find(&transaction).Error; err != nil {
				response.Message = types.ERROR_INTERNAL_SERVER
				return c.JSON(http.StatusInternalServerError, response)
			}
			transactions = append(transactions, transaction...)
		}

		// iterate through all transaction
		for _, transaction := range transactions {
			userFound := false
			friendFound := false
			if transaction.BillOwner == id {
				userFound = true
			}
			if transaction.BillOwner == friend {
				friendFound = true
			}
			// iterate through all item in transaction to check if user or friend is consumer
			for _, item := range transaction.Items {
				// get item detail from item id
				itemDetail := entities.Item{}
				if err := db.Find(&itemDetail, item).Error; err != nil {
					response.Message = types.ERROR_INTERNAL_SERVER
					return c.JSON(http.StatusInternalServerError, response)
				}
				// iterate through all consumer in item
				for _, consumer := range itemDetail.Consumer {
					if consumer == id {
						userFound = true
					}
					if consumer == friend {
						friendFound = true
					}
					if userFound && friendFound {
						break
					}
				}
				if userFound && friendFound {
					break
				}
			}
			if userFound && friendFound {
				commonTransaction[friend] = commonTransaction[friend] + 1
			}
		}
	}

	// get 3 highest common transaction
	highestCommonTransaction := []responses.BuddyDetail{}
	for friend, common := range commonTransaction {
		// get friend name and color
		friendDetail := entities.User{}
		condition := entities.User{ID: friend}
		if err := db.Where(&condition).Find(&friendDetail).Error; err != nil {
			response.Message = types.ERROR_INTERNAL_SERVER
			return c.JSON(http.StatusInternalServerError, response)
		}
		highestCommonTransaction = append(highestCommonTransaction, responses.BuddyDetail{
			Name:  friendDetail.Name,
			Color: friendDetail.Color,
			Count: common,
		})
	}

	// sort highestCommonTransaction
	sort.Slice(highestCommonTransaction, func(i, j int) bool {
		return highestCommonTransaction[i].Count > highestCommonTransaction[j].Count
	})

	// return
	response.Message = types.SUCCESS
	for i := 0; i < len(highestCommonTransaction); i++ {
		if i == 0 {
			response.Data.Buddy1 = highestCommonTransaction[i]
		} else if i == 1 {
			response.Data.Buddy2 = highestCommonTransaction[i]
		} else if i == 2 {
			response.Data.Buddy3 = highestCommonTransaction[i]
		}
	}

	return c.JSON(http.StatusOK, response)
}

func (con *statisticController) ExpenseChart(c echo.Context) error {
	db := con.db
	response := entities.Response[responses.ChartResponse]{}

	id := c.Get("id").(uuid.UUID)

	dateNow := time.Now()
	month := dateNow.Month()

	// get expense from user within this month
	expenses := []entities.Expense{}
	condition := entities.Expense{UserID: id}
	if err := db.Where(&condition).Find(&expenses).Error; err != nil {
		response.Message = types.ERROR_INTERNAL_SERVER
		return c.JSON(http.StatusInternalServerError, response)
	}

	// filter expense within this month
	expensesThisMonth := []entities.Expense{}
	for _, expense := range expenses {
		if expense.Date.Month() == month {
			expensesThisMonth = append(expensesThisMonth, expense)
		}
	}

	// initialize daily expense
	dailyExpense := []float64{}
	if month == time.January || month == time.March || month == time.May || month == time.July || month == time.August || month == time.October || month == time.December {
		for i := 0; i < 31; i++ {
			dailyExpense = append(dailyExpense, 0)
		}
	} else if month == time.April || month == time.June || month == time.September || month == time.November {
		for i := 0; i < 30; i++ {
			dailyExpense = append(dailyExpense, 0)
		}
	} else if month == time.February {
		for i := 0; i < 28; i++ {
			dailyExpense = append(dailyExpense, 0)
		}
	}

	// get daily expense
	totalExpense := 0.0
	for _, expense := range expensesThisMonth {
		dailyExpense[expense.Date.Day()-1] = dailyExpense[expense.Date.Day()-1] + expense.Amount
		totalExpense = totalExpense + expense.Amount
	}

	// return
	response.Message = types.SUCCESS
	response.Data.Month = month.String()
	response.Data.DailyExpense = dailyExpense
	response.Data.TotalExpense = totalExpense

	return c.JSON(http.StatusOK, response)
}
