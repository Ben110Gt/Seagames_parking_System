package util

import (
	"encoding/base64"

	"github.com/skip2/go-qrcode"
)

func GenerateQRCodeBase64(content string) (string, error) {

	var png []byte
	png, err := qrcode.Encode(content, qrcode.Medium, 256)
	if err != nil {
		return "", err
	}
	base64Image := base64.StdEncoding.EncodeToString(png)

	return "data:image/png;base64 format," + base64Image, nil
}
