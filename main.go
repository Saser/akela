package main

import (
	"log"
	"os"
)

func main() {
	logger := log.New(os.Stderr, "[akela] ", log.Ldate|log.Ltime|log.Lshortfile)
	logger.Println("Hello, world!")
}
