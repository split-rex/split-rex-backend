package factories

import (
	"split-rex-backend/types"
	"time"

	"github.com/bxcodec/faker/v4"
	"github.com/google/uuid"
)

type GroupFactory struct {
	GroupID   uuid.UUID         `gorm:"primaryKey"`
	Name      string            `gorm:"not null"`
	MemberID  types.ArrayOfUUID `gorm:"not null"`
	StartDate time.Time         `gorm:"not null"`
	EndDate   time.Time         `gorm:"not null"`
}

// random group, will be deleted after created
func (gf *GroupFactory) Init() {
	if gf.Name == "" {
		gf.Name = faker.Word()
	}
	// member id make from people
	if gf.StartDate.IsZero() {
		gf.StartDate = time.Now()
	}
	if gf.EndDate.IsZero() {
		gf.EndDate = time.Now().Add(time.Hour * 24)
	}
}

func (gf *GroupFactory) GroupA() {
	id, _ := uuid.Parse("29295e79-5281-4453-8375-beb104c86beb")
	gf.GroupID = id
	gf.Name = "groupA"
	gf.StartDate, _ = time.Parse("2023-03-18 00:27:42.437 +0700", "2023-03-18 00:27:42.437 +0700")
	gf.EndDate, _ = time.Parse("2023-03-19 00:27:42.437 +0700", "2023-03-19 00:27:42.437 +0700")
}

// for edit group
func (gf *GroupFactory) GroupB() {
	id, _ := uuid.Parse("657ff578-36df-4904-a39d-3c0007ea8a4a")
	gf.GroupID = id
	gf.Name = "groupB"
	gf.StartDate, _ = time.Parse("2023-03-18 00:33:55.770 +0700", "2023-03-18 00:33:55.770 +0700")
	gf.EndDate, _ = time.Parse("2023-03-19 00:33:55.770 +0700", "2023-03-19 00:33:55.770 +0700")
}


