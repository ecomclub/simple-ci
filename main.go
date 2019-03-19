package main

import (
  "log"
  "net/http"
  "os"
)

func main() {
  // start logging
  logFile := os.Getenv("LOGS_FILE")
  // log to file
  f, err := os.OpenFile(logFile, os.O_RDWR | os.O_CREATE | os.O_APPEND, 0644)
  if err != nil {
    log.Fatalf("Error opening file: %v", err)
  }
  defer f.Close()
  log.SetOutput(f)

  log.Println("------")
  log.Println("Starting simple CI service")
}
