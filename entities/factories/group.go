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
	id, _ := uuid.Parse("c2d77a48-4f0c-402e-91f1-746eb8c99482")
	gf.GroupID = id
	gf.Name = "groupA"
	gf.MemberID = append(gf.MemberID, uuid.MustParse("9368048d-7e77-4fc1-b81e-a5abfef844e1"))
	gf.StartDate, _ = time.Parse("2023-03-18 00:27:42.437 +0700", "2023-03-18 00:27:42.437 +0700")
	gf.EndDate, _ = time.Parse("2023-03-19 00:27:42.437 +0700", "2023-03-19 00:27:42.437 +0700")
}

// for edit group
func (gf *GroupFactory) GroupB() {
	id, _ := uuid.Parse("f76f6d07-08d1-4d75-91ad-05fe6d003990")
	gf.GroupID = id
	gf.Name = "groupB"
	gf.MemberID = append(gf.MemberID, uuid.MustParse("e7d56a3d-930f-45aa-9fcf-a154f2e2db8c"))
	gf.MemberID = append(gf.MemberID, uuid.MustParse("af902382-60d8-4dc6-bd76-dc3f1f061e7a"))
	gf.MemberID = append(gf.MemberID, uuid.MustParse("9368048d-7e77-4fc1-b81e-a5abfef844e1"))
	gf.StartDate, _ = time.Parse("2023-03-18 00:33:55.770 +0700", "2023-03-18 00:33:55.770 +0700")
	gf.EndDate, _ = time.Parse("2023-03-19 00:33:55.770 +0700", "2023-03-19 00:33:55.770 +0700")
}
