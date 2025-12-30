package infr

import (
	"log/slog"
	"md/internal/infr/logger"
	"os"

	"github.com/joho/godotenv"
)

func init() {
	godotenv.Load(".env.local")
	godotenv.Load(".env")
}

func init() {
	slog.SetDefault(slog.New(logger.NewHandler(os.Stderr)))
}
