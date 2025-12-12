package request

// DataKeyBackupPayload instructs the backend what encrypted key material should be stored.
type DataKeyBackupPayload struct {
	EncryptedDataKey string `json:"encryptedDataKey" validate:"required"`
	Salt             string `json:"salt" validate:"required"`
}
