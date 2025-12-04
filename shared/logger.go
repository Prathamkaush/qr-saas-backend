// shared/logger.go
package shared

import (
	"log"
	"os"
)

var Logger = newLogger()

func newLogger() *log.Logger {
	return log.New(os.Stdout, "[qr-saas] ", log.LstdFlags|log.Lshortfile)
}
