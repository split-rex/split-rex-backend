package requests

type LoginRequest struct {
	Email    string `json:"email" form:"email" query:"email" validate:"email"`
	Password string `json:"password" form:"password" query:"password"`
}
