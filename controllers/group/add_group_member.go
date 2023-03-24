package controllers

import (
	"net/http"
	"split-rex-backend/entities"
	"split-rex-backend/entities/requests"
	"split-rex-backend/types"

	"github.com/labstack/echo/v4"
)

func (con *groupController) AddGroupMember(c echo.Context) error {
	db := con.db
	response := entities.Response[string]{}

	request := requests.AddGroupMemberRequest{}
	if err := c.Bind(&request); err != nil {
		response.Message = types.ERROR_BAD_REQUEST
		return c.JSON(http.StatusBadRequest, response)
	}

	for _, member := range request.Friends_id {
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

	// add group member uuid to group
	group := entities.Group{}
	condition := entities.Group{GroupID: request.Group_id}

	if err := db.Model(&group).Where(&condition).Error; err != nil {
		tx.Rollback()
		response.Message = err.Error()
		return c.JSON(http.StatusInternalServerError, response)
	}
	newGroupMembers := append(group.MemberID, request.Friends_id...)
	if err := db.Model(&group).Where(&condition).Updates(entities.Group{
		MemberID: newGroupMembers,
	}).Error; err != nil {
		tx.Rollback()
		response.Message = err.Error()
		return c.JSON(http.StatusInternalServerError, response)
	}

	// adding group uuid to user groups
	for _, memberID := range request.Friends_id {
		user := entities.User{}
		condition := entities.User{ID: memberID}

		if err := tx.Find(&user, &condition).Error; err != nil {
			tx.Rollback()
			response.Message = err.Error()
			return c.JSON(http.StatusInternalServerError, response)
		}
		if err := tx.Model(&user).Where(&condition).Update("groups", append(user.Groups, request.Group_id)).Error; err != nil {
			tx.Rollback()
			response.Message = err.Error()
			return c.JSON(http.StatusInternalServerError, response)
		}
	}

	tx.Commit()
	response.Message = types.SUCCESS

	return c.JSON(http.StatusOK, response)

}
