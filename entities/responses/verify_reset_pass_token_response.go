package responses

type VerifyResetPassTokenResponse struct {
	EncryptedToken     string `json:"encrypted_token"`
}
