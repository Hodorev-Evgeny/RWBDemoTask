package response

type ErrorResponse struct {
	Error   string `json:"error"`
	Massage string `json:"massage"`
}
