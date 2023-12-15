package common

import (
	"log"
	"os"

	"go.uber.org/zap"
)

func touchFile(name string) error {
	file, err := os.OpenFile(name, os.O_RDONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	return file.Close()
}

func GetLogger(isProd bool) *zap.SugaredLogger {
	var config zap.Config
	if isProd {
		config = zap.NewProductionConfig()
	} else {
		config = zap.NewDevelopmentConfig()
	}

	err := os.MkdirAll("logs", os.ModePerm)
	if err != nil {
		log.Print(err)
	}

	err = touchFile("logs/log.txt")
	if err != nil {
		log.Print(err)
	}

	err = touchFile("logs/error-log.txt")
	if err != nil {
		log.Print(err)
	}

	config.OutputPaths = []string{"stdout", "logs/log.txt"}
	config.ErrorOutputPaths = []string{"stderr", "logs/error-log.txt"}

	logger, err := config.Build()
	if err != nil {
		log.Print(err)
	}

	return logger.Sugar()
}
