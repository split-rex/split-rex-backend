package responses

type PercentageResponse struct {
	OwedPercentage int `json:"owed_percentage"`
	LentPercentage int `json:"lent_percentage"`
}

type MutationResponse struct {
	TotalPaid     float64          `json:"total_paid"`
	TotalReceived float64          `json:"total_received"`
	ListMutation  []MutationDetail `json:"list_mutation"`
}

type MutationDetail struct {
	Name         string  `json:"name"`
	Color        uint    `json:"color"`
	MutationType string  `json:"mutation_type"`
	Amount       float64 `json:"amount"`
}
