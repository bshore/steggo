package process

import (
	"os"
	"path/filepath"
)

// WriteEmbeddedFile writes the new file to the output directory
func WriteEmbeddedFile(data []byte, out, ext string) error {
	err := os.MkdirAll(out, 0777)
	if err != nil {
		return err
	}
	newFile, err := os.Create(filepath.Join(out, "output."+ext))
	if err != nil {
		return err
	}
	defer newFile.Close()
	_, err = newFile.Write(data)
	return err
}
