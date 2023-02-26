package requests

type FriendRequest struct {
	User_id   string `json:"user_id" form:"user_id" query:"user_id""`
	Friend_id string `json:"friend_id" form:"friend_id" query:"friend_id"`
}
