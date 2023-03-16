package requests

type UpdateProfileRequest struct {
	Name     string `json:"name" form:"name" query:"name"`
	Password string `json:"password" form:"password" query:"password"`
	Color float32 `json:"color" form:"color" query:"color"`
}
