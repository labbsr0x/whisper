package types

// RegistrationPage defines the information needed to load a registration page
type RegistrationPage struct {
	Page
}

// UpdatePage defins the information needed to load a update user credentials page
type UpdatePage struct {
	Page
	Username string
	Email    string
}

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
