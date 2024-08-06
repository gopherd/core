// Package mime64 provides functionality to encode various types of content into base64 MIME format.
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

// EncodeType encodes content with a specified MIME type into a base64 MIME string.
func EncodeType(mimeType string, content []byte) []byte {
	var buf bytes.Buffer
	buf.WriteString("data:")
	buf.WriteString(mimeType)
	buf.WriteString(";base64,")
	buf.WriteString(base64.StdEncoding.EncodeToString(content))
	return buf.Bytes()
}

// Encode detects the MIME type of the content and encodes it into a base64 MIME string.
func Encode(content []byte) []byte {
	return EncodeType(http.DetectContentType(content), content)
}

// EncodeToString encodes content into a base64 MIME string.
func EncodeToString(content []byte) string {
	return string(Encode(content))
}

// EncodeReader reads content from an io.Reader and encodes it into a base64 MIME string.
func EncodeReader(r io.Reader) (string, error) {
	content, err := io.ReadAll(r)
	if err != nil {
		return "", err
	}
	return EncodeToString(content), nil
}

// EncodeFile reads content from a file and encodes it into a base64 MIME string.
func EncodeFile(filename string) (string, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return EncodeToString(content), nil
}

// EncodeURL fetches content from a URL and encodes it into a base64 MIME string.
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

// EncodeImagePNG encodes an image.Image as PNG into a base64 MIME string.
func EncodeImagePNG(img image.Image) (string, error) {
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return "", err
	}
	return EncodeToString(buf.Bytes()), nil
}

// EncodeImageJPEG encodes an image.Image as JPEG into a base64 MIME string.
func EncodeImageJPEG(img image.Image) (string, error) {
	var buf bytes.Buffer
	if err := jpeg.Encode(&buf, img, nil); err != nil {
		return "", err
	}
	return EncodeToString(buf.Bytes()), nil
}

// EncodeImageGIF encodes an image.Image as GIF into a base64 MIME string.
func EncodeImageGIF(img image.Image) (string, error) {
	var buf bytes.Buffer
	if err := gif.Encode(&buf, img, nil); err != nil {
		return "", err
	}
	return EncodeToString(buf.Bytes()), nil
}
