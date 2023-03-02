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

	if err := db.Create(
		&entities.Group{
			GroupID:   uuid.New(),
			Name:      request.Name,
			MemberID:  request.MemberID,
			StartDate: request.StartDate,
			EndDate:   request.EndDate,
		}).Error; err != nil {
		response.Message = types.ERROR_INTERNAL_SERVER
		return c.JSON(http.StatusInternalServerError, response)
	}

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
	groups := []entities.Group{}

	// PRINT member_id column from groups table

	if err := db.Select("group_id").Where(`? = ANY(string_to_array(member_id, ','))`, userID).Table("groups").Find(&groups).Error; err != nil {
		response.Message = err.Error() + " PERTAMA"
		return c.JSON(http.StatusInternalServerError, response)
	}

	data := []responses.UserGroupResponse{}

	for _, group := range groups {
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
	groupID, _ := uuid.Parse("fb3313e2-492c-4689-9bd7-f9dfd5fcfb4a")

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
