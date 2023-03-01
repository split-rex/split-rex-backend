package requests

type FriendRequest struct {
	Friend_id string `json:"friend_id" form:"friend_id" query:"friend_id"`
}
