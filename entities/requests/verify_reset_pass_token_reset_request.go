package requests

type VerifyResetPassTokenRequest struct {
	Email          string `json:"email" form:"email" query:"email" validate:"email"`
	Code           string `json:"code" form:"code" query:"code" validate:"code"`
}
