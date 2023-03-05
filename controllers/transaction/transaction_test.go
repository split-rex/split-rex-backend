package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"split-rex-backend/entities/requests"
	"split-rex-backend/types"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestUserCreateTransaction(t *testing.T) {
	groupID, _ := uuid.Parse("6251ac85-e43d-4b88-8779-588099df5008")
	billOwnerID, _ := uuid.Parse("6251ac85-e43d-4b88-8779-588099df5008")

	date := time.Date(2023, 3, 3, 0, 0, 0, 0, time.UTC)

	transaction := &requests.UserCreateTransactionRequest{
		Name:        "Transaction 1",
		Description: "Transaction 1 Description",
		GroupID:     groupID,
		Date:        date,
		Subtotal:    100,
		Tax:         10,
		Service:     10,
		Total:       120,
		BillOwner:   billOwnerID,
		Items:       types.ArrayOfUUID{},
	}
	body, _ := json.Marshal(transaction)

	res, err := http.Post("http://localhost:8080/userCreateTransaction",
		"application/json", bytes.NewBuffer(body))
	if err != nil {
		t.Error(err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Error("Expected status code 200, got ", res.StatusCode)
	}
}
