package controllers

import (
	"net/http"
	"split-rex-backend/entities"
	"split-rex-backend/entities/requests"
	"split-rex-backend/types"

	"github.com/google/uuid"
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

	// get all new members
	newMembers := []entities.User{}
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

		newMembers = append(newMembers, user)
	}

	// get user friend to check later
	id := c.Get("id").(uuid.UUID)
	currUserFriends := entities.Friend{}
	if err := db.Find(&currUserFriends, id).Error; err != nil {
		response.Message = err.Error()
		return c.JSON(http.StatusInternalServerError, response)
	}

	// if new member friend is not friend with current user, then error bad request
	for _, userFriendId := range request.Friends_id {
		if !currUserFriends.Friend_id.Contains(userFriendId) {
			response.Message = types.ERROR_BAD_REQUEST
			return c.JSON(http.StatusBadRequest, response)
		}
	}

	// begin transaction
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// get group
	group := entities.Group{}
	if err := db.Find(&group, request.Group_id).Error; err != nil {
		tx.Rollback()
		response.Message = err.Error()
		return c.JSON(http.StatusInternalServerError, response)
	}

	// check if member already exist in group, then exclude
	newMemberIDs := types.ArrayOfUUID{}
	for _, newMemberId := range request.Friends_id {
		if !group.MemberID.Contains(newMemberId) {
			newMemberIDs = append(newMemberIDs, newMemberId)
		}
	}
	// add group member uuid to group
	newGroupMembers := append(group.MemberID, newMemberIDs...)

	// then update groups memberid
	condition := entities.Group{GroupID: request.Group_id}
	if err := db.Model(&group).Where(&condition).Updates(entities.Group{
		MemberID: newGroupMembers,
	}).Error; err != nil {
		tx.Rollback()
		response.Message = err.Error()
		return c.JSON(http.StatusInternalServerError, response)
	}

	// adding group uuid to new member's groups
	for _, member := range newMembers {
		user := entities.User{}
		condition := entities.User{ID: member.ID}

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
