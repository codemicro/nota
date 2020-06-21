package helpers

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"github.com/disintegration/imaging"
	"os"
	"strings"
)

func RandomHex(n int) (string, error) {
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func GetFileExtension(s string) (string, error) {
	if !strings.ContainsAny(s, ".") {
		return "", errors.New("no points found in input (impossible to parse)")
	}

	splitString := strings.Split(s, ".")
	return splitString[len(splitString)-1], nil

}

func IsStringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func SaveBytesToDisk(path string, content []byte) error {
	f, err := os.Create(path)
	defer f.Close()
	if err != nil {
		return err
	}
	_, err = f.Write(content)
	if err != nil {
		return err
	}

	return nil
}

func MimeTypeToFileExt(mimeType string) (string, error) {
	mapping := map[string]string{
		"image/jpeg": "jpg",
		"image/png":  "png",
	}[mimeType]

	if mapping == "" {
		return "", errors.New("MIME type not found")
	}

	return mapping, nil
}

func RotateImage(imgPath string) error {
	img, err := imaging.Open(imgPath)
	if err != nil {
		return err
	}
	img = imaging.Rotate90(img)
	err = imaging.Save(img, imgPath)
	if err != nil {
		return err
	}
	return nil
}
