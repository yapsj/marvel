package util

import (
	"io"
	"log"
	"os"
)

func NewLogger(w io.Writer, prefix string) *log.Logger {
	logger := log.New(w, prefix, log.Ldate|log.Ltime|log.Lshortfile)
	multiWriter := io.MultiWriter(os.Stdout, w)
	logger.SetOutput(multiWriter)
	return logger
}

func NewInfoLogger(w io.Writer) *log.Logger {
	return NewLogger(w, "INFO: ")
}

func NewErrorLogger(w io.Writer) *log.Logger {
	return NewLogger(w, "ERROR: ")
}
