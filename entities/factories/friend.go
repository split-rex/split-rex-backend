package factories

import (
	"split-rex-backend/types"

	"github.com/google/uuid"
)

type FriendFactory struct {
	ID           uuid.UUID 
	Friend_id    types.ArrayOfUUID
	Req_received types.ArrayOfUUID
	Req_sent     types.ArrayOfUUID
}

func (ff *FriendFactory) Init(){
	
}
