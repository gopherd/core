package mime64_test

import (
	"bytes"
	"encoding/base64"
	"image"
	"image/color"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/gopherd/core/encoding/mime64"
)

func TestEncodeType(t *testing.T) {
	tests := []struct {
		name     string
		mimeType string
		content  []byte
		want     string
	}{
		{
			name:     "Plain text",
			mimeType: "text/plain",
			content:  []byte("Hello, World!"),
			want:     "data:text/plain;base64,SGVsbG8sIFdvcmxkIQ==",
		},
		{
			name:     "Empty content",
			mimeType: "application/octet-stream",
			content:  []byte{},
			want:     "data:application/octet-stream;base64,",
		},
		{
			name:     "Binary data",
			mimeType: "application/octet-stream",
			content:  []byte{0xFF, 0x00, 0xAA, 0x55},
			want:     "data:application/octet-stream;base64,/wCqVQ==",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := mime64.EncodeType(tt.mimeType, tt.content)
			if string(got) != tt.want {
				t.Errorf("EncodeType() = %v, want %v", string(got), tt.want)
			}
		})
	}
}

func TestEncode(t *testing.T) {
	tests := []struct {
		name    string
		content []byte
		want    string
	}{
		{
			name:    "JPEG image",
			content: []byte{0xFF, 0xD8, 0xFF, 0xE0, 0x00, 0x10, 0x4A, 0x46, 0x49, 0x46},
			want:    "data:image/jpeg;base64,/9j/4AAQSkZJRg==",
		},
		{
			name:    "PNG image",
			content: []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A},
			want:    "data:image/png;base64,iVBORw0KGgo=",
		},
		{
			name:    "Plain text",
			content: []byte("Hello, World!"),
			want:    "data:text/plain; charset=utf-8;base64,SGVsbG8sIFdvcmxkIQ==",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := mime64.Encode(tt.content)
			if string(got) != tt.want {
				t.Errorf("Encode() = %v, want %v", string(got), tt.want)
			}
		})
	}
}

func TestEncodeToString(t *testing.T) {
	content := []byte("Hello, World!")
	want := "data:text/plain; charset=utf-8;base64,SGVsbG8sIFdvcmxkIQ=="

	got := mime64.EncodeToString(content)
	if got != want {
		t.Errorf("EncodeToString() = %v, want %v", got, want)
	}
}

func TestEncodeReader(t *testing.T) {
	content := "Hello, World!"
	reader := strings.NewReader(content)
	want := "data:text/plain; charset=utf-8;base64,SGVsbG8sIFdvcmxkIQ=="

	got, err := mime64.EncodeReader(reader)
	if err != nil {
		t.Fatalf("EncodeReader() error = %v", err)
	}
	if got != want {
		t.Errorf("EncodeReader() = %v, want %v", got, want)
	}

	// Test with a failing reader
	failingReader := &failingReader{err: io.ErrUnexpectedEOF}
	_, err = mime64.EncodeReader(failingReader)
	if err == nil {
		t.Error("EncodeReader() expected error, got nil")
	}
}

func TestEncodeFile(t *testing.T) {
	// Create a temporary file
	content := []byte("Hello, World!")
	tmpfile, err := os.CreateTemp("", "example")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write(content); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatalf("Failed to close temp file: %v", err)
	}

	want := "data:text/plain; charset=utf-8;base64,SGVsbG8sIFdvcmxkIQ=="

	got, err := mime64.EncodeFile(tmpfile.Name())
	if err != nil {
		t.Fatalf("EncodeFile() error = %v", err)
	}
	if got != want {
		t.Errorf("EncodeFile() = %v, want %v", got, want)
	}

	// Test with a non-existent file
	_, err = mime64.EncodeFile("non_existent_file.txt")
	if err == nil {
		t.Error("EncodeFile() expected error for non-existent file, got nil")
	}
}

func TestEncodeURL(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte("Hello, World!"))
	}))
	defer server.Close()

	got, err := mime64.EncodeURL(server.URL)
	if err != nil {
		t.Fatalf("EncodeURL() error = %v", err)
	}

	// Check if the result starts with the correct MIME type
	if !strings.HasPrefix(got, "data:text/plain") {
		t.Errorf("EncodeURL() result doesn't start with correct MIME type. Got: %v", got)
	}

	// Check if the base64 encoded part is correct
	expectedBase64 := "SGVsbG8sIFdvcmxkIQ=="
	if !strings.HasSuffix(got, expectedBase64) {
		t.Errorf("EncodeURL() base64 part is incorrect. Got: %v, Want suffix: %v", got, expectedBase64)
	}

	// Test with a non-existent URL
	_, err = mime64.EncodeURL("http://non-existent-url.com")
	if err == nil {
		t.Error("EncodeURL() expected error for non-existent URL, got nil")
	}
}

func TestEncodeImagePNG(t *testing.T) {
	img := createTestImage()

	got, err := mime64.EncodeImagePNG(img)
	if err != nil {
		t.Fatalf("EncodeImagePNG() error = %v", err)
	}

	if !strings.HasPrefix(got, "data:image/png;base64,") {
		t.Errorf("EncodeImagePNG() result doesn't have correct prefix")
	}

	// Decode the base64 part and check if it's a valid PNG
	parts := strings.SplitN(got, ",", 2)
	if len(parts) != 2 {
		t.Fatalf("EncodeImagePNG() result doesn't have expected format")
	}

	decoded, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		t.Fatalf("Failed to decode base64: %v", err)
	}

	_, err = png.Decode(bytes.NewReader(decoded))
	if err != nil {
		t.Errorf("EncodeImagePNG() didn't produce a valid PNG: %v", err)
	}
}

func TestEncodeImageJPEG(t *testing.T) {
	img := createTestImage()

	got, err := mime64.EncodeImageJPEG(img)
	if err != nil {
		t.Fatalf("EncodeImageJPEG() error = %v", err)
	}

	if !strings.HasPrefix(got, "data:image/jpeg;base64,") {
		t.Errorf("EncodeImageJPEG() result doesn't have correct prefix")
	}

	// Decode the base64 part and check if it's a valid JPEG
	parts := strings.SplitN(got, ",", 2)
	if len(parts) != 2 {
		t.Fatalf("EncodeImageJPEG() result doesn't have expected format")
	}

	decoded, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		t.Fatalf("Failed to decode base64: %v", err)
	}

	_, err = jpeg.Decode(bytes.NewReader(decoded))
	if err != nil {
		t.Errorf("EncodeImageJPEG() didn't produce a valid JPEG: %v", err)
	}
}

func TestEncodeImageGIF(t *testing.T) {
	img := createTestImage()

	got, err := mime64.EncodeImageGIF(img)
	if err != nil {
		t.Fatalf("EncodeImageGIF() error = %v", err)
	}

	if !strings.HasPrefix(got, "data:image/gif;base64,") {
		t.Errorf("EncodeImageGIF() result doesn't have correct prefix")
	}

	// Decode the base64 part and check if it's a valid GIF
	parts := strings.SplitN(got, ",", 2)
	if len(parts) != 2 {
		t.Fatalf("EncodeImageGIF() result doesn't have expected format")
	}

	decoded, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		t.Fatalf("Failed to decode base64: %v", err)
	}

	_, err = gif.Decode(bytes.NewReader(decoded))
	if err != nil {
		t.Errorf("EncodeImageGIF() didn't produce a valid GIF: %v", err)
	}
}

// Helper function to create a test image
func createTestImage() image.Image {
	img := image.NewRGBA(image.Rect(0, 0, 100, 100))
	for y := 0; y < 100; y++ {
		for x := 0; x < 100; x++ {
			img.Set(x, y, color.RGBA{uint8(x % 256), uint8(y % 256), 100, 255})
		}
	}
	return img
}

// failingReader is a helper type for testing reader failures
type failingReader struct {
	err error
}

func (r *failingReader) Read(p []byte) (n int, err error) {
	return 0, r.err
}

func xor(x, y int) uint8 {
	var v = x ^ y
	var b0 = uint8((v & 0xff))
	var b1 = uint8((v >> 8) & 0xff)
	var b2 = uint8((v >> 16) & 0xff)
	var b3 = uint8((v >> 24) & 0xff)
	return b0 ^ b1 ^ b2 ^ b3
}

func TestEncodeXOR(t *testing.T) {
	const width = 128
	const height = 128
	var img = image.NewRGBA(image.Rect(0, 0, width, height))
	for i := 0; i < width; i++ {
		for j := 0; j < height; j++ {
			img.Set(i, j, color.Gray{Y: xor(i, j)})
		}
	}
	var out bytes.Buffer
	if err := png.Encode(&out, img); err != nil {
		t.Fatalf("encode png error: %v", err)
	}
	const want = "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAIAAAACACAIAAABMXPacAAAClElEQVR4nOzdUYoaURCFYUvcuEtzKS7BRx/EECoI4ZLQMNPTf3H+epBDcWkGKl/NTSPM5XQ6lXVcnavq/X77edTn5fcQzufX6/VnIOafzQpgCOixPJ/Pv9aT/f37CmAI4OzEtKwAhoBa6vF4rM3P9Dz/XecVwBDA2YlpWQEMAT0Wwr04ra8AhgDOTkzLCmAIqM11v9+3H/5M2+f/8+GEfwXJn/4O8HdA9qf/Dzi4rwCGAM5OTMsKYAiopaa8T59+XgEMAZydmJYVwBDQYyHci9P6CmAI4OzEtKwAhoDarW63234P7xr98yuAIYCzE9OyAhgCeiyEe3FaXwEMAZydmJYVwBBQS015nz79vAIYAjg7MS0rgCGgx0K4F6f1FcAQwNmJaVkBDAG1uWjfr5/+fAUwBHB2YlpWAENAj4VwL07rK4AhgLMT07ICGAJqqSnv06efVwBDAGcnpmUFMAT0WAj34rS+AhgCODsxLSuAIaDG1vV6PfpH+FIpgCGAsxPTsgIYAnoshHtxWl8BDAGcnZiWFcAQUEtNeZ8+/bwCGAI4OzEtK4AhoMdCuBen9RXAEMDZiWlZAQwBtblo36+f/nwFMARwdmJaVgBDQI+FcC9O6yuAIYCzE9OyAhgCaqkp79Onn1cAQwBnJ6ZlBTAE9FgI9+K0vgIYAjg7MS0rgCGgdiv/fsD/SwEMAZydmJYVwBDQYyHci9P6CmAI4OzEtKwAhoBaasr79OnnFcAQwNmJaVkBDAE9FsK9OK2vAIYAzk5MywpgCKjNRft+/fTnK4AhgLMT07ICGAJ6LIR7cVpfAQwBnJ2YlhXAEFBLTXmfPv28AhgCODsxLSuAIaDHQrgXp/UVwBDA2YlpWQEHf/4KAAD//xgbFb2maABzAAAAAElFTkSuQmCC"
	var got = mime64.EncodeToString(out.Bytes())
	if want != got {
		t.Fatalf("want %v, but got %v", want, got)
	}
	t.Logf("got: %v", got)
}
