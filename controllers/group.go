package controllers

import (
	"net/http"
	"split-rex-backend/configs/database"
	"split-rex-backend/entities"
	"split-rex-backend/entities/requests"
	"split-rex-backend/entities/responses"
	"split-rex-backend/types"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func UserCreateGroupController(c echo.Context) error {
	db := database.DB.GetConnection()
	response := entities.Response[string]{}

	request := requests.UserCreateGroupRequest{}
	if err := c.Bind(&request); err != nil {
		response.Message = types.ERROR_BAD_REQUEST
		return c.JSON(http.StatusBadRequest, response)
	}

	// check if all corresponding member exist in user table
	for _, member := range request.MemberID {
		user := entities.User{}
		if err := db.Find(&user, member).Error; err != nil {
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

	group := &entities.Group{
		GroupID:   uuid.New(),
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
	user := entities.User{}
	for _, memberID := range request.MemberID {
		condition := entities.User{ID: memberID}
		if err := tx.Model(&user).Where(&condition).Update("groups", append(user.Groups, group.GroupID)).Error; err != nil {
			tx.Rollback()
			response.Message = err.Error()
			return c.JSON(http.StatusInternalServerError, response)
		}
	}

	tx.Commit()
	response.Message = types.SUCCESS
	return c.JSON(http.StatusAccepted, response)
}

func EditGroupInfoController(c echo.Context) error {
	db := database.DB.GetConnection()
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

func UserGroupsController(c echo.Context) error {
	db := database.DB.GetConnection()
	response := entities.Response[[]responses.UserGroupResponse]{}

	userID := c.Get("id").(uuid.UUID)

	groupIDs := []uuid.UUID{}
	condition := entities.User{ID: userID}
	if err := db.Model(&entities.User{}).Where(&condition).Pluck("groups", &groupIDs).Error; err != nil {
		response.Message = err.Error()
		return c.JSON(http.StatusInternalServerError, response)
	}

	data := []responses.UserGroupResponse{}

	for _, groupID := range groupIDs {
		group := entities.Group{}
		condition := entities.Group{GroupID: groupID}

		if err := db.Where(&condition).Find(&group).Error; err != nil {
			response.Message = err.Error()
			return c.JSON(http.StatusInternalServerError, response)
		}

		data = append(data, responses.UserGroupResponse{
			GroupID:   group.GroupID,
			Name:      group.Name,
			MemberID:  group.MemberID,
			StartDate: group.StartDate,
			EndDate:   group.EndDate,

			// TODO: Calculate Things & Determine Types
			Type:         "HARDCODED",
			TotalUnpaid:  0,
			TotalExpense: 0,
		})
	}
	response.Message = types.SUCCESS
	response.Data = data

	return c.JSON(http.StatusAccepted, response)
}

func GroupDetailController(c echo.Context) error {
	db := database.DB.GetConnection()
	response := entities.Response[responses.GroupDetailResponse]{}

	group := entities.Group{}
	groupID, _ := uuid.Parse(c.Param("id"))

	condition := entities.Group{GroupID: groupID}
	if err := db.Where(&condition).Find(&group).Error; err != nil {
		response.Message = err.Error()
		return c.JSON(http.StatusInternalServerError, response)
	}

	data := responses.GroupDetailResponse{
		GroupID:    group.GroupID,
		Name:       group.Name,
		MemberID:   group.MemberID,
		StartDate:  group.StartDate,
		EndDate:    group.EndDate,
		ListMember: []responses.MemberDetail{},
	}

	for _, memberID := range group.MemberID {
		data.ListMember = append(data.ListMember,
			responses.MemberDetail{
				ID: memberID,

				// TODO: Calculate Things & Determine Types
				Type:        "hardcoded",
				TotalUnpaid: 0,
			})
	}

	response.Message = types.SUCCESS
	response.Data = data

	return c.JSON(http.StatusAccepted, response)
}

func GroupTransactionsController(c echo.Context) error {
	// db := database.DB.GetConnection()
	// config := configs.Config.GetMetadata()
	response := entities.Response[string]{}
	return c.JSON(http.StatusAccepted, response)
}
