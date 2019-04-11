package types

// AddUserCredentialRequestPayload defines the payload for adding a user
type AddUserCredentialRequestPayload struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// AddUserCredentialResponsePayload defines the response payload after adding a user
type AddUserCredentialResponsePayload struct {
	UserCredentialID string `json:"userCredentialId"`
}

// UpdateUserCredentialRequestPayload defines the payload for updating a user
type UpdateUserCredentialRequestPayload struct {
	NewUsername string `json:"newUsername"`
	NewPassword string `json:"newPassword"`
	OldPassword string `json:"oldPassword"`
}
