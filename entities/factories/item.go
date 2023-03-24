package factories

import (
	"split-rex-backend/types"

	"github.com/google/uuid"
)

type ItemFactory struct {
	ItemID        uuid.UUID         `gorm:"not null;unique"`
	Name          string            `gorm:"not null"`
	Quantity      int               `gorm:"not null"`
	Price         float64           `gorm:"not null"`
	Consumer      types.ArrayOfUUID `gorm:"not null"`
}

func (itf *ItemFactory) Init(){
	if itf.Name==""{
		itf.Name = "test item 1"
	}
	if itf.Quantity==0{
		itf.Quantity = 1
	}
	if itf.Price==0{
		itf.Price = 1000.0
	}
}