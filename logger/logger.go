package logger

import (
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

var Log *logrus.Logger

func InitLogger() {
	_ = godotenv.Load()

	Log = logrus.New()
	Log.Out = os.Stdout
	Log.SetFormatter(&logrus.JSONFormatter{})

	logLevelStr := os.Getenv("LOG_LEVEL")
	if logLevelStr == "" {
		logLevelStr = "info"
	}

	logLevel, err := logrus.ParseLevel(strings.ToLower(logLevelStr))
	if err != nil {
		log.Fatalf("Invalid log level: %v", err)
	}
	Log.SetLevel(logLevel)
}
