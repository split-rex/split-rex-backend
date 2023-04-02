package responses

import "split-rex-backend/types"

type ProfileResponse struct {
	User_id      string            `json:"user_id" form:"user_id" query:"user_id"`
	Email        string            `json:"email" form:"email" query:"email"`
	Username     string            `json:"username" form:"username" query:"username"`
	Fullname     string            `json:"fullname" form:"fullname" query:"fullname"`
	Color        uint              `json:"color" form:"color" query:"color"`
	Payment_info types.PaymentInfo `json:"payment_info" form:"payment_info" query:"payment_info"`
}
