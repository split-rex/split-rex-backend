package responses

type PercentageResponse struct {
	TotalOwed      float64 `json:"total_owed"`
	TotalLent      float64 `json:"total_lent"`
	OwedPercentage int     `json:"owed_percentage"`
	LentPercentage int     `json:"lent_percentage"`
}

type MutationResponse struct {
	TotalPaid     float64          `json:"total_paid"`
	TotalReceived float64          `json:"total_received"`
	ListMutation  []MutationDetail `json:"list_mutation"`
}

type BuddyResponse struct {
	Buddy1 BuddyDetail `json:"buddy1"`
	Buddy2 BuddyDetail `json:"buddy2"`
	Buddy3 BuddyDetail `json:"buddy3"`
}

type ChartResponse struct {
	Month        string    `json:"month"`
	TotalExpense float64   `json:"total_expense"`
	DailyExpense []float64 `json:"daily_expense"`
}

type MutationDetail struct {
	Name         string  `json:"name"`
	Color        uint    `json:"color"`
	MutationType string  `json:"mutation_type"`
	Amount       float64 `json:"amount"`
}

type BuddyDetail struct {
	Name  string `json:"name"`
	Color uint   `json:"color"`
	Count int    `json:"count"`
}
