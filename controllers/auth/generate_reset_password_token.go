package controllers

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"io"
	mathRand "math/rand"
	"net/http"
	"os"
	"split-rex-backend/entities"
	"split-rex-backend/entities/requests"
	"split-rex-backend/entities/responses"
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

	// check to db if email exist
	db := con.db
	response := entities.Response[responses.GenerateResetPassTokenResponse]{}

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

	// generate random for passwordToken
	key := []byte(os.Getenv("RESET_PASS_KEY"))

	// generate a new aes cipher using our 32 byte long key
	cip, err := aes.NewCipher(key)
	// if there are any errors, handle them
	if err != nil {
		response.Message = types.ERROR_INTERNAL_SERVER
		return c.JSON(http.StatusInternalServerError, response)
	}

	gcm, err := cipher.NewGCM(cip)
	// if any error generating new GCM
	// handle them
	if err != nil {
		response.Message = types.ERROR_INTERNAL_SERVER
		return c.JSON(http.StatusInternalServerError, response)
	}

	// creates a new byte array the size of the nonce
	// which must be passed to Seal
	nonce := make([]byte, gcm.NonceSize())
	// populates our nonce with a cryptographically secure
	// random sequence
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		response.Message = types.ERROR_INTERNAL_SERVER
		return c.JSON(http.StatusInternalServerError, response)
	}
	token := randStringBytesMaskImprSrcUnsafe(10)
	// encryptedToken := gcm.Seal(nonce, nonce, []byte(token), nil)
	encryptedToken := base64.StdEncoding.EncodeToString(gcm.Seal(nonce, nonce, []byte(token), nil))

	// generate random for code
	code := randStringBytesMaskImprSrcUnsafe(6)
	// encryptedCode := base64.StdEncoding.EncodeToString(gcm.Seal(nonce, nonce, code, nil))

	// save to db raw code
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

	// return the passwordResetToken
	response.Data.EncryptedToken = encryptedToken
	response.Data.Token = token
	response.Data.Code = code

	response.Message = types.SUCCESS
	return c.JSON(http.StatusOK, response)
}
