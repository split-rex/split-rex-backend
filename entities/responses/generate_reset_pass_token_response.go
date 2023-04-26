package responses

type GenerateResetPassTokenResponse struct {
	EncryptedToken     string `json:"encrypted_token"`
}
