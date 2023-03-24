package controllers

import (
	"split-rex-backend/entities"
	"split-rex-backend/entities/requests"
	"split-rex-backend/types"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func createItems(db *gorm.DB, items []requests.ItemRequest) (types.ArrayOfUUID, error) {
	arrayOfItemUUID := types.ArrayOfUUID{}
	for _, item := range items {
		newItem := entities.Item{
			ItemID:   uuid.New(),
			Name:     item.Name,
			Quantity: item.Quantity,
			Price:    item.Price,
			Consumer: item.Consumer,
		}

		if err := db.Save(&newItem).Error; err != nil {
			return arrayOfItemUUID, err
		}

		arrayOfItemUUID = append(arrayOfItemUUID, newItem.ItemID)
	}

	return arrayOfItemUUID, nil
}
