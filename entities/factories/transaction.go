package factories

import (
	"split-rex-backend/types"
	"time"

	"github.com/google/uuid"
)

type TransactionFactory struct {
	TransactionID uuid.UUID         `gorm:"not null;unique"`
	Name          string            `gorm:"not null"`
	Description   string            `gorm:"not null"`
	GroupID       uuid.UUID         `gorm:"not null"`
	Date          time.Time         `gorm:"not null"`
	Subtotal      float64           `gorm:"not null"`
	Tax           float64           `gorm:"not null"`
	Service       float64           `gorm:"not null"`
	Total         float64           `gorm:"not null"`
	BillOwner     uuid.UUID         `gorm:"not null"`
	Items         types.ArrayOfUUID `gorm:"not null"`
}

func (tf *TransactionFactory) Init(){
	if tf.Name == "" {
		tf.Name = "New Transaction"
	}
	if tf.Description == "" {
		tf.Description = "New Transaction Description"
	}
	
	// if tf.GroupID == uuid.Nil {
	// 	tf.GroupID = "0b865d7f-e40e-4440-905e-eccf2caaa6ed"
	// }
	
	if tf.Date.IsZero() {
		tf.Date = time.Now()
	}
	if tf.Subtotal == 0 {
		tf.Subtotal = 1000.0
	}
	if tf.Tax == 0 {
		tf.Tax = 100.0
	}
	if tf.Service == 0 {
		tf.Service = 100.0
	}
	if tf.Total == 0 {
		tf.Total = 1200.0
	}

	// if tf.BillOwner == 0 {
	// 	tf.BillOwner = types.EncryptedString("testing")
	// }

	// if tf.Items.Count() == 0 {
	// 	item, _ := uuid.Parse("6251ac85-e43d-4b88-8779-588099df5008")
	// 	tf.Items = append(tf.Items,item)
	// }
}


