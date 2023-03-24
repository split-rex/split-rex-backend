package requests

import (
	"split-rex-backend/types"

	"github.com/google/uuid"
)

type AddGroupMemberRequest struct {
	Group_id uuid.UUID`json:"group_id" form:"group_id" query:"group_id"`
	Friends_id types.ArrayOfUUID `json:"friends_id" form:"friends_id" query:"friends_id"`
}
