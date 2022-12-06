package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/ieee0824/zeroentd"
	"github.com/rs/zerolog"
)

func main() {
	writer, err := zeroentd.New("zeroentd.test.log")
	if err != nil {
		if e, ok := err.(fmt.Formatter); ok {
			log.Fatalf("%+v\n", e)
		}
		log.Fatalln(err)
	}

	defer writer.Close()

	logger := zerolog.New(io.MultiWriter(writer, os.Stdout)).With().Timestamp().Caller().Logger()

	logger.Info().Msg("test info msg")
	logger.Error().Err(errors.New("any error")).Msg("test error")
}
