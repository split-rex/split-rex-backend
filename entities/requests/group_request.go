package requests

import (
	"split-rex-backend/types"
	"time"

	"github.com/google/uuid"
)

// @patrickamadeus this time.Time & []uuid.UUID can't be parsed, types must be string
// then after string, need to be parsed
type UserCreateGroupRequest struct {
	Name      string            `json:"name" form:"name" query:"name"`
	MemberID  types.ArrayOfUUID `json:"member_id" form:"member_id" query:"member_id"`
	StartDate time.Time         `json:"start_date" form:"start_date" query:"start_date" default:"time.Now()"`
	EndDate   time.Time         `json:"end_date" form:"end_date" query:"end_date" default:"time.Now()"`
}

type EditGroupInfoRequest struct {
	GroupID   uuid.UUID `json:"group_id" form:"group_id" query:"group_id"`
	Name      string    `json:"name" form:"name" query:"name"`
	StartDate time.Time `json:"start_date" form:"start_date" query:"start_date"`
	EndDate   time.Time `json:"end_date" form:"end_date" query:"end_date"`
}
