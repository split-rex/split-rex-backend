package requests

type GenerateResetPassTokenRequest struct {
	Email string `json:"email" form:"email" query:"email" validate:"email"`
}
