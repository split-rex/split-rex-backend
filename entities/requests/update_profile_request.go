package requests

type UpdateProfileRequest struct {
	Name  string `json:"name" form:"name" query:"name"`
	Color uint   `json:"color" form:"color" query:"color"`
}
