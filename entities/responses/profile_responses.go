package responses

type ProfileResponse struct {
	User_id  string `json:"user_id" form:"user_id" query:"user_id"`
	Username string `json:"username" form:"username" query:"username"`
	Fullname string `json:"fullname" form:"fullname" query:"fullname"`
}
