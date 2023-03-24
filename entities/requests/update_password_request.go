package requests

type UpdatePasswordRequest struct {
	OldPassword string `json:"old_password" form:"old_password" query:"old_password"`
	NewPassword string `json:"new_password" form:"new_password" query:"new_password"`
}
