package responses

type GenerateResetPassTokenResponse struct {
	EncryptedToken     string `json:"encrypted_token"`
	OnlyEncryptedToken string `json:"only_encrypted_token"`
	Token              string `json:"token"`
}
