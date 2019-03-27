package start

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

// logger is zap.logger instance.
var logger *zap.Logger

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("File .env not found, reading configuration from ENV")
	}
	switch os.Getenv("LOGGER_TYPE") {
	case "PRODUCTION":
		zapLogger, err := zap.NewProduction()
		if err != nil {
			log.Fatalf("can't init logger")
		}
		logger = zapLogger
	default:
		zapLogger, err := zap.NewDevelopment()
		if err != nil {
			log.Fatalf("can't init logger")
		}
		logger = zapLogger
	}
}

type startServer struct {
}

func CreateServer() (*startServer, error) {

	return &startServer{}, nil
}

func (s *startServer) Run() error {
	return nil
}
