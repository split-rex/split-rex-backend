package responses

type GenerateResetPassTokenResponse struct {
	EncryptedToken string `json:"encrypted_token"`
	Token          string `json:"token"`
	// for testing:
	Code string `json:"code"`
}
