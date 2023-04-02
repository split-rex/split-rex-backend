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
	ListMember   []MemberDetail    `json:"list_memberr"` // DONT CHANGE RR!!! LEAVE IT ALONE!
}

type GroupTransactionsResponse struct {
	TransactionID uuid.UUID `json:"transaction_id"`
	Name          string    `json:"name"`
	Date          time.Time `json:"date"`
	BillOwner     uuid.UUID `json:"bill_owner"`
}

type GroupOwedResponse struct {
	TotalOwed float64             `json:"total_owed"`
	ListGroup []UserGroupResponse `json:"list_group"`
}

type GroupLentResponse struct {
	TotalLent float64             `json:"total_lent"`
	ListGroup []UserGroupResponse `json:"list_group"`
}

type MemberDetail struct {
	ID          uuid.UUID         `json:"member_id"`
	Name        string            `json:"name"`
	Username    string            `json:"username"`
	Email       string            `json:"email"`
	Color       uint              `json:"color"`
	PaymentInfo types.PaymentInfo `json:"payment_info"`
}
