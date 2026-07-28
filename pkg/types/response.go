package types

type SuccessResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Errors  interface{} `json:"errors"`
}

type ErrorResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Errors  interface{} `json:"errors"`
}
