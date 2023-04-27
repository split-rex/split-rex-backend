package configs

import (
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type Metadata struct {
	ApplicationName         string
	LoginExpirationDuration time.Duration
	JWTSigningMethod        jwt.SigningMethod
	JWTSignatureKey         []byte
	ResetPasswordKey        []byte
	EmailSenderName         string
	EmailSenderAddress      string
	EmailSenderPassword     string
}

type AuthConfig struct {
	metadata Metadata
	once     sync.Once
}

func (authConfig *AuthConfig) lazyInit() {
	authConfig.once.Do(func() {
		applicationName := os.Getenv("APPLICATION_NAME")
		numberOfSeconds, err := strconv.Atoi(os.Getenv("LOGIN_EXPIRATION_DURATION"))
		if err != nil {
			panic(err)
		}
		loginExpirationDuration := time.Duration(numberOfSeconds) * time.Second
		jwtSigningMethod := jwt.SigningMethodHS256
		jwtSignatureKey := []byte(os.Getenv("JWT_SIGNATURE_KEY"))

		authConfig.metadata.ApplicationName = applicationName
		authConfig.metadata.LoginExpirationDuration = loginExpirationDuration
		authConfig.metadata.JWTSigningMethod = jwtSigningMethod
		authConfig.metadata.JWTSignatureKey = jwtSignatureKey
		authConfig.metadata.ResetPasswordKey = []byte(os.Getenv("RESET_PASS_KEY"))
		authConfig.metadata.EmailSenderName = os.Getenv("EMAIL_SENDER_NAME")
		authConfig.metadata.EmailSenderAddress = os.Getenv("EMAIL_SENDER_ADDRESS")
		authConfig.metadata.EmailSenderPassword = os.Getenv("EMAIL_SENDER_PASSWORD")
	})
}

func (authConfig *AuthConfig) lazyTestInit() {
	authConfig.once.Do(func() {
		applicationName := "split-rex-backend-testing"
		numberOfSeconds, err := strconv.Atoi("86000")
		if err != nil {
			panic(err)
		}
		loginExpirationDuration := time.Duration(numberOfSeconds) * time.Second
		jwtSigningMethod := jwt.SigningMethodHS256
		jwtSignatureKey := []byte("1234")

		authConfig.metadata.ApplicationName = applicationName
		authConfig.metadata.LoginExpirationDuration = loginExpirationDuration
		authConfig.metadata.JWTSigningMethod = jwtSigningMethod
		authConfig.metadata.JWTSignatureKey = jwtSignatureKey
	})
}

func (authConfig *AuthConfig) GetTestMetadata() Metadata {
	authConfig.lazyTestInit()
	return authConfig.metadata
}

func (authConfig *AuthConfig) GetMetadata() Metadata {
	authConfig.lazyInit()
	return authConfig.metadata
}

var Config = &AuthConfig{}
