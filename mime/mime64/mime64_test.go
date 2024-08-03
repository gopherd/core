package mime64_test

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"testing"

	"github.com/gopherd/core/mime/mime64"
)

var x = ""

func xor(x, y int) uint8 {
	var v = x ^ y
	var b0 = uint8((v & 0xff))
	var b1 = uint8((v >> 8) & 0xff)
	var b2 = uint8((v >> 16) & 0xff)
	var b3 = uint8((v >> 24) & 0xff)
	return b0 ^ b1 ^ b2 ^ b3
}

func TestEncode(t *testing.T) {
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
