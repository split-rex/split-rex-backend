package controllers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func makeFriendRequest(t *testing.T, user_id uuid.UUID, friend_id uuid.UUID) {

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("id", user_id)
	c.SetParamNames("username")
	c.SetParamValues("lalallaa")

	
}

// buat request dari user a ke d --> /searchUserToAdd, /makeFriendRequest

// buat request dari user b ke a --> /searchUserToAdd, /makeFriendRequest

// buat request dari user c ke a --> /searchUserToAdd, /makeFriendRequest

// dari user a cek friend req sent  --> /friendRequestSent

// dari user a cek friend req received --> /friendRequestReceived

// user a reject user b --> /rejectRequest

// cek friend req received harusnya tinggal c --> /friendRequestReceived

// user a acc user c -->  /acceptRequest

// cek friend req received a harusnya kosong--> /friendRequestReceived

// d acc a --> /acceptRequest

// cek friend req sent di a hrusnya kosong --> /friendRequestSent

// cek friend a hrusnya udh nambah c + d --> /userFriendList

// nambah a ke c  --> /searchUserToAdd lagi harusnya gbs --> sekarang masih

// nambah a ke a -->

// BUAT YANG search user to add, kalo lagi friend req dikasi tau
