package controllers

import (
	"net/http"
	"split-rex-backend/entities"
	"split-rex-backend/entities/responses"
	"split-rex-backend/types"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func (con *groupController) GroupLent(c echo.Context) error {
	db := con.db
	response := entities.Response[responses.GroupLentResponse]{}

	id := c.Get("id").(uuid.UUID)
	user := entities.User{}
	if err := db.Find(&user, id).Error; err != nil {
		response.Message = types.ERROR_INTERNAL_SERVER
		return c.JSON(http.StatusInternalServerError, response)
	}

	// iterate through user's groups to get groups details
	totalLent := 0
	groups := []responses.GroupDetailResponse{}
	for _, groupID := range user.Groups {
		group := entities.Group{}
		condition := entities.Group{GroupID: groupID}
		if err := db.Where(condition).Find(&group).Error; err != nil {
			response.Message = types.ERROR_INTERNAL_SERVER
			return c.JSON(http.StatusInternalServerError, response)
		}

		// then for each group, search for transactions existed in group id
		transactions := []entities.Transaction{}
		conditionTransaction := entities.Transaction{GroupID: groupID}
		if err := db.Where(conditionTransaction).Find(&transactions).Error; err != nil {
			response.Message = types.ERROR_INTERNAL_SERVER
			return c.JSON(http.StatusInternalServerError, response)
		}

		// TODO: then for each transaction count for lent
		totalLent += 10000

		// then map to group detail
		groupDetail := responses.GroupDetailResponse{
			GroupID:   group.GroupID,
			Name:      group.Name,
			StartDate: group.StartDate,
			EndDate:   group.EndDate,
		}

		// then for each groups, iterate through members on user
		listMember := []responses.MemberDetail{}
		for _, memberId := range group.MemberID {
			member := entities.User{}
			if err := db.Find(&member, memberId).Error; err != nil {
				response.Message = types.ERROR_INTERNAL_SERVER
				return c.JSON(http.StatusInternalServerError, response)
			}

			listMember = append(listMember, responses.MemberDetail{
				ID:       member.ID,
				Name:     member.Name,
				Username: member.Username,
				Email:    member.Email,

				// TODO: calculate
				Type:        "HARDCODED",
				TotalUnpaid: 0,
			})
		}

		groupDetail.ListMember = listMember

		groups = append(groups, groupDetail)
	}

	// then return all
	response.Message = types.SUCCESS
	response.Data.TotalLent = totalLent
	response.Data.ListGroup = groups

	return c.JSON(http.StatusOK, response)
}
