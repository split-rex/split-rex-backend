package responses

import (
	"split-rex-backend/types"
	"time"

	"github.com/google/uuid"
)

type UserGroupResponse struct {
	GroupID      uuid.UUID         `json:"group_id"`
	Name         string            `json:"name"`
	MemberID     types.ArrayOfUUID `json:"member_id"`
	StartDate    time.Time         `json:"start_date"`
	EndDate      time.Time         `json:"end_date"`
	Type         string            `json:"type"`
	TotalUnpaid  float64           `json:"total_unpaid"`
	TotalExpense float64           `json:"total_expense"`
}

type GroupDetailResponse struct {
	GroupID    uuid.UUID         `json:"group_id"`
	Name       string            `json:"name"`
	MemberID   types.ArrayOfUUID `json:"member_id"`
	StartDate  time.Time         `json:"start_date"`
	EndDate    time.Time         `json:"end_date"`
	ListMember []MemberDetail    `json:"list_member"`
}

type GroupTransactionsResponse struct {
	TransactionID uuid.UUID         `json:"group_id"`
	Name          string            `json:"name"`
	Description   string            `json:"description"`
	Total         float64           `json:"total"`
	BillOwner     uuid.UUID         `json:"bill_owner"`
	ListMember    types.ArrayOfUUID `json:"list_member"`
}

type MemberDetail struct {
	ID          uuid.UUID `json:"member_id"`
	Type        string    `json:"type"`
	TotalUnpaid float64   `json:"total_unpaid"`
}
