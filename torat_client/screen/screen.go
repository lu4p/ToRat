// +build !android

package screen

import (
	"bytes"
	"image/png"

	"github.com/vova616/screenshot"
)

// Take takes a screenshot and returns it as byte slice
func Take() []byte {
	img, err := screenshot.CaptureScreen()
	if err != nil {
		panic(err)
	}
	buf := new(bytes.Buffer)
	if err != nil {
		panic(err)
	}
	err = png.Encode(buf, img)
	if err != nil {
		panic(err)
	}
	return buf.Bytes()
}
