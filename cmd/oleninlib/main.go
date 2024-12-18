package main

import (
	"os"

	"github.com/rs/zerolog"
)

func main() {

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log := zerolog.New(os.Stdout).With().Timestamp().Logger()

	err := run(log)
	if err != nil {
		log.Fatal().Err(err).Send()
	}
}

func run(log zerolog.Logger) error {
	log.Info().Msg("OleninBot started...")

	log.Info().Msg("OleninBot finished successfully...")
	return nil
}
