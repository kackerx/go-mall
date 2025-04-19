package resources

import (
	"bytes"
	"embed"
	"io"
)

//go:embed *
var f embed.FS

func LoadResourceFile(path string) (io.Reader, error) {
	bs, err := f.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(bs), nil
}
