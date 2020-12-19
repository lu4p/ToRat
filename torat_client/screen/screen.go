// +build !android

package screen

import (
	"bytes"
	"image/png"

	"github.com/vova616/screenshot"
)

// Take takes a screenshot and returns it as byte slice
func Take() ([]byte, error) {
	img, err := screenshot.CaptureScreen()
	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	if err != nil {
		return nil, err
	}

	err = png.Encode(buf, img)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
