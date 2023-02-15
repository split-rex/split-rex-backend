package tests

import (
	"split-rex-backend/configs/database"
	"split-rex-backend/entities"
	"split-rex-backend/types"
	"testing"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

func TestRegister(t *testing.T) {
	err := godotenv.Load()
	if err != nil {
		t.Fatalf("err loading: %v", err)
	}

	db := database.DB.GetConnection()

	if err := db.Create(&entities.User{
		ID:       uuid.New(),
		Name:     "unit_test",
		Username: "unit_test",
		Password: types.EncryptedString("unit_test"),
		Email:    "unit_test@gmail.com",
	}).Error; err != nil {
		t.Error("Failed: Register user")
	}
}

func TestLogin(t *testing.T) {
	err := godotenv.Load()
	if err != nil {
		t.Fatalf("err loading: %v", err)
	}

	db := database.DB.GetConnection()

	user := entities.User{}
	if err := db.Where(&entities.User{
		Username: "unit_test",
	}).Find(&user).Error; err != nil {
		t.Error("Failed: Login user")
	}

	db.Delete(&user)
}
