package utils

import (
	log "github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
	"os"
	"sync"
)

var instance *log.Logger
var once sync.Once

func GetLogger() *log.Logger {
	once.Do(func() {
		instance = &log.Logger{
			Out:   os.Stderr,
			Level: log.DebugLevel,
			Formatter: &prefixed.TextFormatter{
				TimestampFormat: "2006-01-02 15:04:05",
				FullTimestamp:   true,
				ForceFormatting: true,
			},
		}
	})
	return instance
}
