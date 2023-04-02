package requests

type PaymentInfo struct {
	Payment_method string `json:"payment_method" form:"payment_method" query:"payment_method"`
	Account_number uint   `json:"account_number" form:"account_number" query:"account_number"`
	Account_name   string `json:"account_name" form:"account_name" query:"account_name"`
}
