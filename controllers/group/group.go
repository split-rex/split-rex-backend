package controllers

import (
	"net/http"
	"split-rex-backend/entities"
	"split-rex-backend/entities/requests"
	"split-rex-backend/entities/responses"
	"split-rex-backend/types"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func (con *groupController) UserCreateGroup(c echo.Context) error {
	db := con.db
	response := entities.Response[string]{}

	request := requests.UserCreateGroupRequest{}
	if err := c.Bind(&request); err != nil {
		response.Message = types.ERROR_BAD_REQUEST
		return c.JSON(http.StatusBadRequest, response)
	}

	user_id := c.Get("id").(uuid.UUID)
	request.MemberID = append(request.MemberID, user_id)

	// check if all corresponding member exist in user table
	for _, member := range request.MemberID {
		user := entities.User{}
		if err := db.Find(&user, member).Error; err != nil {
			response.Message = types.ERROR_BAD_REQUEST
			return c.JSON(http.StatusBadRequest, response)
		}
		if user.Name == "" {
			response.Message = types.ERROR_BAD_REQUEST
			return c.JSON(http.StatusBadRequest, response)
		}
	}

	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	newUUID := uuid.New()
	group := &entities.Group{
		GroupID:   newUUID,
		Name:      request.Name,
		MemberID:  request.MemberID,
		StartDate: request.StartDate,
		EndDate:   request.EndDate,
	}

	if err := tx.Create(group).Error; err != nil {
		tx.Rollback()
		response.Message = types.ERROR_INTERNAL_SERVER
		return c.JSON(http.StatusInternalServerError, response)
	}

	// adding group uuid to user groups
	for _, memberID := range request.MemberID {
		user := entities.User{}
		condition := entities.User{ID: memberID}

		if err := tx.Find(&user, &condition).Error; err != nil {
			tx.Rollback()
			response.Message = err.Error()
			return c.JSON(http.StatusInternalServerError, response)
		}
		if err := tx.Model(&user).Where(&condition).Update("groups", append(user.Groups, group.GroupID)).Error; err != nil {
			tx.Rollback()
			response.Message = err.Error()
			return c.JSON(http.StatusInternalServerError, response)
		}
	}

	tx.Commit()
	response.Message = types.SUCCESS
	response.Data = newUUID.String()
	return c.JSON(http.StatusCreated, response)
}

func (con *groupController) EditGroupInfo(c echo.Context) error {
	db := con.db
	response := entities.Response[string]{}

	request := requests.EditGroupInfoRequest{}
	if err := c.Bind(&request); err != nil {
		response.Message = err.Error()
		return c.JSON(http.StatusBadRequest, response)
	}

	group := entities.Group{}
	condition := entities.Group{GroupID: request.GroupID}

	if err := db.Model(&group).Where(&condition).Updates(entities.Group{
		Name:      request.Name,
		StartDate: request.StartDate,
		EndDate:   request.EndDate,
	}).Error; err != nil {
		response.Message = err.Error()
		return c.JSON(http.StatusInternalServerError, response)
	}

	response.Message = types.SUCCESS
	return c.JSON(http.StatusAccepted, response)
}

func (con *groupController) UserGroups(c echo.Context) error {
	db := con.db
	response := entities.Response[[]responses.UserGroupResponse]{}

	userID := c.Get("id").(uuid.UUID)

	user := entities.User{}
	condition := entities.User{ID: userID}
	if err := db.Model(&entities.User{}).Where(&condition).Find(&user).Error; err != nil {
		response.Message = err.Error()
		return c.JSON(http.StatusInternalServerError, response)
	}

	data := []responses.UserGroupResponse{}

	for _, groupID := range user.Groups {
		totalUnpaid := 0.0
		totalExpense := 0.0
		groupType := types.TYPE_GROUP_EQUAL

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
		conditionPayment := entities.Payment{GroupID: groupID, UserID1: userID}
		if err := db.Where(&conditionPayment).Find(&payments).Error; err != nil {
			response.Message = types.ERROR_INTERNAL_SERVER
			return c.JSON(http.StatusInternalServerError, response)
		}

		// compute totalOwed from payments
		for _, payment := range payments {
			totalUnpaid = totalUnpaid + payment.TotalUnpaid
		}

		// if totalOwed is negative then not in groupOwed
		if totalUnpaid > 0 {
			groupType = types.TYPE_GROUP_OWED
		} else if totalUnpaid < 0 {
			groupType = types.TYPE_GROUP_LENT
		} else {
			groupType = types.TYPE_GROUP_EQUAL
		}

		groupResponse := responses.UserGroupResponse{
			GroupID:      group.GroupID,
			Name:         group.Name,
			MemberID:     group.MemberID,
			StartDate:    group.StartDate,
			EndDate:      group.EndDate,
			Type:         groupType,
			TotalUnpaid:  totalUnpaid,
			TotalExpense: totalExpense}

		for _, memberID := range group.MemberID {
			// get member detail
			user := entities.User{}
			condition := entities.User{ID: memberID}
			if err := db.Where(&condition).Find(&user).Error; err != nil {
				response.Message = types.ERROR_INTERNAL_SERVER
				return c.JSON(http.StatusInternalServerError, response)
			}
			groupResponse.ListMember = append(groupResponse.ListMember,
				responses.MemberDetail{
					ID:       memberID,
					Name:     user.Name,
					Username: user.Username,
					Email:    user.Email,
					Color:    user.Color,
				})
		}

		data = append(data, groupResponse)
	}
	response.Message = types.SUCCESS
	response.Data = data

	return c.JSON(http.StatusAccepted, response)
}

func (h *groupController) GroupDetail(c echo.Context) error {
	db := h.db
	response := entities.Response[responses.UserGroupResponse]{}

	totalUnpaid := 0.0
	totalExpense := 0.0
	groupType := types.TYPE_GROUP_EQUAL

	group := entities.Group{}
	groupID, _ := uuid.Parse(c.QueryParam("id"))
	userID := c.Get("id").(uuid.UUID)

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
	conditionPayment := entities.Payment{GroupID: groupID, UserID1: userID}
	if err := db.Where(&conditionPayment).Find(&payments).Error; err != nil {
		response.Message = types.ERROR_INTERNAL_SERVER
		return c.JSON(http.StatusInternalServerError, response)
	}

	// compute totalOwed from payments
	for _, payment := range payments {
		totalUnpaid = totalUnpaid + payment.TotalUnpaid
	}

	// if totalOwed is negative then not in groupOwed
	if totalUnpaid > 0 {
		groupType = types.TYPE_GROUP_OWED
	} else if totalUnpaid < 0 {
		groupType = types.TYPE_GROUP_LENT
	} else {
		groupType = types.TYPE_GROUP_EQUAL
	}

	data := responses.UserGroupResponse{
		GroupID:      group.GroupID,
		Name:         group.Name,
		MemberID:     group.MemberID,
		StartDate:    group.StartDate,
		EndDate:      group.EndDate,
		ListMember:   []responses.MemberDetail{},
		Type:         groupType,
		TotalUnpaid:  totalUnpaid,
		TotalExpense: totalExpense,
	}

	for _, memberID := range group.MemberID {
		// get member detail
		user := entities.User{}
		condition := entities.User{ID: memberID}
		if err := db.Where(&condition).Find(&user).Error; err != nil {
			response.Message = types.ERROR_INTERNAL_SERVER
			return c.JSON(http.StatusInternalServerError, response)
		}
		data.ListMember = append(data.ListMember,
			responses.MemberDetail{
				ID:       memberID,
				Name:     user.Name,
				Username: user.Username,
				Email:    user.Email,
				Color:    user.Color,
			})
	}

	response.Message = types.SUCCESS
	response.Data = data

	return c.JSON(http.StatusAccepted, response)
}

func (h *groupController) GroupTransactions(c echo.Context) error {
	db := h.db
	response := entities.Response[[]responses.GroupTransactionsResponse]{}

	groupID, _ := uuid.Parse(c.QueryParam("id"))

	transactions := []entities.Transaction{}
	condition := entities.Transaction{GroupID: groupID}
	if err := db.Where(&condition).Find(&transactions).Error; err != nil {
		response.Message = err.Error()
		return c.JSON(http.StatusAccepted, response)
	}

	data := []responses.GroupTransactionsResponse{}

	// TODO: get all members from consumer field in Item for ListMember
	listMember := types.ArrayOfUUID{}

	for _, transaction := range transactions {
		data = append(data, responses.GroupTransactionsResponse{
			TransactionID: transaction.TransactionID,
			Name:          transaction.Name,
			Description:   transaction.Description,
			Total:         transaction.Total,
			BillOwner:     transaction.BillOwner,
			ListMember:    listMember,
		})
	}

	response.Message = types.SUCCESS
	response.Data = data

	return c.JSON(http.StatusAccepted, response)
}
