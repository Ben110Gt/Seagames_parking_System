package util

import "fmt"

// BarcodeURL returns the relative URL path for fetching the QR code image.
func BarcodeURL(code string) string {
	return fmt.Sprintf("/api/v1/barcode/%s", code)
}

// APIResponse is the standard JSON envelope for all responses.
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

func OK(data interface{}, msg string) APIResponse {
	return APIResponse{Success: true, Message: msg, Data: data}
}

func Fail(msg string) APIResponse {
	return APIResponse{Success: false, Error: msg}
}
