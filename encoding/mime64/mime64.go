package mime64

import (
	"bytes"
	"encoding/base64"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"net/http"
	"os"
)

// EncodeType encodes mime content with specified mime type
func EncodeType(mimeType string, content []byte) []byte {
	var buf bytes.Buffer
	buf.WriteString("data:")
	buf.WriteString(mimeType)
	buf.WriteString(";base64,")
	buf.WriteString(base64.StdEncoding.EncodeToString(content))
	return buf.Bytes()
}

// Encode encodes mime content
func Encode(mime []byte) []byte {
	return EncodeType(http.DetectContentType(mime), mime)
}

//----------------------------------------------------------
// Helper functions

// EncodeToString encodes mime content to string
func EncodeToString(content []byte) string {
	return string(Encode(content))
}

// EncodeReader encodes mime content from reader to string
func EncodeReader(r io.Reader) (string, error) {
	content, err := io.ReadAll(r)
	if err != nil {
		return "", err
	}
	return EncodeToString(content), nil
}

// EncodeFile encodes mime content from file to string
func EncodeFile(filename string) (string, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return EncodeToString(content), nil
}

// EncodeURL encodes mime content from url to string
func EncodeURL(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return EncodeToString(content), nil
}

// EncodeImagePNG encodes image/png to string
func EncodeImagePNG(img image.Image) (string, error) {
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return "", err
	}
	return EncodeToString(buf.Bytes()), nil
}

// EncodeImageJPEG encodes image/jpeg to string
func EncodeImageJPEG(img image.Image) (string, error) {
	var buf bytes.Buffer
	if err := jpeg.Encode(&buf, img, nil); err != nil {
		return "", err
	}
	return EncodeToString(buf.Bytes()), nil
}

// EncodeImageGIF encodes image/gif to string
func EncodeImageGIF(img image.Image) (string, error) {
	var buf bytes.Buffer
	if err := gif.Encode(&buf, img, nil); err != nil {
		return "", err
	}
	return EncodeToString(buf.Bytes()), nil
}
