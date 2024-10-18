package handlers

func BadRequestErrorResp(message string) *CommonErrorObject {
	return &CommonErrorObject{
		Message: message,
		Text:    "BAD_REQUEST",
	}
}

type CommonErrorObject struct {
	Text    string `json:"text" example:""`
	Message string `json:"message" example:""`
}
