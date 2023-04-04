package middlewares

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"split-rex-backend/configs"
	"split-rex-backend/entities"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type AuthClaims struct {
	jwt.RegisteredClaims
	ID uuid.UUID `json:"id"`
}

// GoogleClaims -
type GoogleClaims struct {
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	FirstName     string `json:"given_name"`
	LastName      string `json:"family_name"`
	jwt.StandardClaims
}

func AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		config := configs.Config.GetMetadata()
		response := entities.Response[string]{}

		authHeader := c.Request().Header.Get("Authorization")
		if !strings.Contains(authHeader, "Bearer") {
			response.Message = "ERROR: NO TOKEN PROVIDED"
			return c.JSON(http.StatusUnauthorized, response)
		}

		authString := strings.Replace(authHeader, "Bearer ", "", -1)
		authClaim := AuthClaims{}
		authToken, err := jwt.ParseWithClaims(authString, &authClaim, func(authToken *jwt.Token) (interface{}, error) {
			if method, ok := authToken.Method.(*jwt.SigningMethodHMAC); !ok || method != config.JWTSigningMethod {
				return nil, fmt.Errorf("ERROR: SIGNING METHOD INVALID")
			}
			return config.JWTSignatureKey, nil
		})
		if err != nil {
			response.Message = "ERROR: TOKEN CANNOT BE PARSED"
			return c.JSON(http.StatusInternalServerError, response)
		}
		if !authToken.Valid {
			response.Message = "ERROR: CLAIMS INVALID"
			return c.JSON(http.StatusBadRequest, response)
		}

		c.Set("id", authClaim.ID)

		return next(c)
	}
}

// ValidateGoogleJWT
func ValidateGoogleJWT(tokenString string) (GoogleClaims, error) {
	claimsStruct := GoogleClaims{}

	token, err := jwt.ParseWithClaims(
		tokenString,
		&claimsStruct,
		func(token *jwt.Token) (interface{}, error) {
			pem, err := getGooglePublicKey(fmt.Sprintf("%s", token.Header["kid"]))
			if err != nil {
				return nil, err
			}
			key, err := jwt.ParseRSAPublicKeyFromPEM([]byte(pem))
			if err != nil {
				return nil, err
			}
			return key, nil
		},
	)
	if err != nil {
		return GoogleClaims{}, err
	}

	claims, ok := token.Claims.(*GoogleClaims)
	if !ok {
		return GoogleClaims{}, errors.New("invalid Google JWT token")
	}

	if claims.Issuer != "accounts.google.com" && claims.Issuer != "https://accounts.google.com" {
		return GoogleClaims{}, errors.New("iss is invalid")
	}

	if claims.Audience != "YOUR_CLIENT_ID_HERE" {
		return GoogleClaims{}, errors.New("aud is invalid")
	}

	if claims.ExpiresAt < time.Now().UTC().Unix() {
		return GoogleClaims{}, errors.New("JWT is expired")
	}

	return *claims, nil
}

func getGooglePublicKey(keyID string) (string, error) {
	resp, err := http.Get("https://www.googleapis.com/oauth2/v1/certs")
	if err != nil {
		return "", err
	}
	dat, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	response := map[string]string{}
	err = json.Unmarshal(dat, &response)
	if err != nil {
		return "", err
	}
	key, ok := response[keyID]
	if !ok {
		return "", errors.New("key not found")
	}
	return key, nil
}
