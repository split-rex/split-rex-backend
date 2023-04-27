package requests

type ChangePasswordRequest struct {
	Email          string `json:"email" form:"email" query:"email" validate:"email"`
	Code           string `json:"code" form:"code" query:"code" validate:"code"`
	EncryptedToken string `json:"encrypted_token" form:"encrypted_token" query:"encrypted_token" validate:"encrypted_token"`
	NewPassword    string `json:"new_password" form:"new_password" query:"new_password"`
}
