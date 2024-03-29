package controllers

import (
	mathRand "math/rand"
	"net/http"
	"split-rex-backend/configs"
	"split-rex-backend/entities"
	"split-rex-backend/entities/requests"
	"split-rex-backend/types"

	"time"
	"unsafe"

	"github.com/labstack/echo/v4"
)

const (
	letterBytes   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ123456789"
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

func randStringBytesMaskImprSrcUnsafe(n int) string {
	var src = mathRand.NewSource(time.Now().UnixNano())
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return *(*string)(unsafe.Pointer(&b))
}

func (con *authController) GenerateResetPassTokenController(c echo.Context) error {
	config := configs.Config.GetMetadata()
	// check to db if email exist
	db := con.db
	response := entities.Response[string]{}

	generateTokenRequest := requests.GenerateResetPassTokenRequest{}
	if err := c.Bind(&generateTokenRequest); err != nil {
		response.Message = types.ERROR_BAD_REQUEST
		return c.JSON(http.StatusBadRequest, response)
	}

	if generateTokenRequest.Email == "" {
		response.Message = types.ERROR_BAD_REQUEST
		return c.JSON(http.StatusBadRequest, response)
	}

	user := entities.User{}
	condition := entities.User{Email: generateTokenRequest.Email}
	if err := db.Where(&condition).Find(&user).Error; err != nil {
		response.Message = types.ERROR_INTERNAL_SERVER
		return c.JSON(http.StatusInternalServerError, response)
	}
	// IF USER DOESNT EXIST, WE DO NOT INFORM THE USER DUE TO SECURITY
	if user.Email == "" {
		response.Message = types.SUCCESS
		return c.JSON(http.StatusOK, response)
	}

	// generate random for token
	token := randStringBytesMaskImprSrcUnsafe(12)
	// generate random for code
	code := randStringBytesMaskImprSrcUnsafe(8)

	passwordResetToken := entities.PasswordResetTokens{}
	passwordResetToken.Email = user.Email
	passwordResetToken.Token = token
	passwordResetToken.Code = code
	passwordResetToken.TokenExpiry = time.Now()

	userToken := entities.PasswordResetTokens{}
	condition2 := entities.PasswordResetTokens{Email: generateTokenRequest.Email}
	if err := db.Where(&condition2).Find(&userToken).Error; err != nil {
		response.Message = types.ERROR_INTERNAL_SERVER
		return c.JSON(http.StatusInternalServerError, response)
	}
	if userToken.Token == "" {
		if err := db.Create(&passwordResetToken).Error; err != nil {
			response.Message = types.ERROR_INTERNAL_SERVER
			return c.JSON(http.StatusInternalServerError, response)
		}
	} else {
		if err := db.Model(&userToken).Updates(passwordResetToken).Error; err != nil {
			response.Message = types.ERROR_INTERNAL_SERVER
			return c.JSON(http.StatusInternalServerError, response)
		}
	}
	// send to email the code
	emailSenderName := config.EmailSenderName
	emailSenderAddress := config.EmailSenderAddress
	emailSenderPassword := config.EmailSenderPassword
	sender := NewGmailSender(emailSenderName, emailSenderAddress, emailSenderPassword)

	subject := "Password Reset"
	content := `
	<h1>Request to Reset Your Password</h1>
	<p>Someone requested a password reset at this email address for Split-rex Mobile. To complete the reset password, enter the verification code below: </p>
	<h2>` + code + `</h2>
	<p>Code will expire in 2 minutes. If you did not request a password reset, you can safely ignore this email.</p>
	<p>© 2023 Splitrex</p>
	`
	to := []string{user.Email}

	err := sender.SendEmail(subject, content, to, nil, nil, nil)
	if err != nil {
		response.Message = types.ERROR_INTERNAL_SERVER
		return c.JSON(http.StatusInternalServerError, response)
	}

	response.Message = types.SUCCESS
	return c.JSON(http.StatusOK, response)
}
