package services

import (
	"bufio"
	"encoding/base64"
	"io/ioutil"

	"os"
)

func ImagesToBase64(str_images string) string {
	f, _ := os.Open(str_images)
	defer f.Close()
	reader := bufio.NewReader(f)
	content, _ := ioutil.ReadAll(reader)
	encoded := base64.StdEncoding.EncodeToString(content)
	return encoded
}
