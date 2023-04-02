package requests

type EditPaymentInfoRequest struct {
	Old_payment_method string `json:"old_payment_method" form:"old_payment_method" query:"old_payment_method"`
	Old_account_number uint   `json:"old_account_number" form:"old_account_number" query:"old_account_number"`
	Old_account_name   string `json:"old_account_name" form:"old_account_name" query:"old_account_name"`
	New_payment_method string `json:"new_payment_method" form:"new_payment_method" query:"new_payment_method"`
	New_account_number uint   `json:"new_account_number" form:"new_account_number" query:"new_account_number"`
	New_account_name   string `json:"new_account_name" form:"new_account_name" query:"new_account_name"`
}
