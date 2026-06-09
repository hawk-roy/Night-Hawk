package test

import (
	"io"
	"log"
	"os"
	"path/filepath"
)

func Trace_Debug_Path(traceInfo any) {
	if err := os.MkdirAll("logs", 0755); err != nil {
		log.Fatal(err)
	}

	logPath := filepath.Join("logs", "app.log")

	file, err := os.OpenFile(logPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	writer := io.MultiWriter(os.Stdout, file)

	logger := log.New(writer, "[MYAPP] ", log.Ldate|log.Ltime|log.Lshortfile)

	logger.Println(traceInfo)
}
