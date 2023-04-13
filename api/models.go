package api

/*
Model for the api response body, with the informative message.
*/
type ErroMsg struct {
	Message string `json:"message" validate:"required"` // Message with information about the error
}
