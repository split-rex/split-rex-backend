package factories

import (
	"split-rex-backend/entities/requests"
	"time"

	"github.com/bxcodec/faker/v4"
	"github.com/google/uuid"
)

type TransactionFactory struct {
	TransactionID uuid.UUID            `gorm:"not null;unique"`
	Name          string               `gorm:"not null"`
	Description   string               `gorm:"not null"`
	GroupID       uuid.UUID            `gorm:"not null"`
	Date          time.Time            `gorm:"not null"`
	Subtotal      float64              `gorm:"not null"`
	Tax           float64              `gorm:"not null"`
	Service       float64              `gorm:"not null"`
	Total         float64              `gorm:"not null"`
	BillOwner     uuid.UUID            `gorm:"not null"`
	Items         []requests.ItemRequest `gorm:"not null"`
}

func (tf *TransactionFactory) Init() {
	if tf.Name == "" {
		tf.Name = faker.Word()
	}
	if tf.Description == "" {
		tf.Description = faker.Sentence()
	}

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
}
