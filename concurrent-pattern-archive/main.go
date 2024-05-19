package main

import (
	seeder "concurrentarchive/fileSeeder"
	archive "concurrentarchive/secureArchive"
	"fmt"
	"log"
	"os"
	"time"
)

func main() {
	err := os.MkdirAll("files", 0755)
	if err != nil {
		log.Fatal(err)
	}

	now := time.Now()

	seeder.ExecSeeder(3000, "files")
	archive.ExecSecureArchive("./files")

	since := time.Since(now)
	fmt.Printf("Executed in %f second", since.Seconds())
}
