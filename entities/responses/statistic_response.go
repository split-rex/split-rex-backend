package responses

type PercentageResponse struct {
	OwedPercentage int `json:"owed_percentage"`
	LentPercentage int `json:"lent_percentage"`
}
