package token

import (
	"crypto/md5"
	"fmt"
	"time"

	"github.com/satori/go.uuid"
)

func Generate(seed string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(fmt.Sprintf("%s.%d.%s", uuid.NewV4().String(), time.Now().Unix(), seed))))
}
