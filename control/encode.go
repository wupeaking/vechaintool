package control

import (
	"encoding/base64"
	"fmt"
	"github.com/andlabs/ui"
)

func Base64Encode(content string, log *ui.MultilineEntry) {
	log.SetText("")
	result := base64.StdEncoding.EncodeToString([]byte(content))
	log.Append(fmt.Sprintf("base64 encode result: %s", result))
}