package main

import (
	"archive/zip"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func encryptAES(data []byte, key []byte) (resp []byte, err error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return resp, err
	}

	resp = gcm.Seal(nonce, nonce, data, nil)

	return resp, nil
}

func createZip(sourceDir string) (err error) {
	// zip process
	nowUnix := time.Now().Unix()

	// create output folder
	os.Mkdir("output", 0755)

	zippedFile, err := os.Create(fmt.Sprintf("output/result-%d.zip", nowUnix))
	if err != nil {
		return err
	}
	defer zippedFile.Close()

	zipWriter := zip.NewWriter(zippedFile)
	defer zipWriter.Close()

	// archive file on temp folder to output as file.zip
	err = filepath.WalkDir(sourceDir, func(path string, data fs.DirEntry, err error) error {
		// add file to zip writer
		if !data.IsDir() {
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()

			fileInfo, err := file.Stat()
			if err != nil {
				return err
			}

			zipHeader, err := zip.FileInfoHeader(fileInfo)
			if err != nil {
				return err
			}

			writer, err := zipWriter.CreateHeader(zipHeader)
			if err != nil {
				return err
			}

			_, err = io.Copy(writer, file)
			if err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func EncryptArchive(filePath string) {
	// creating temp folder in root path
	tempDirPath, err := os.MkdirTemp(".", "*")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(tempDirPath)

	// read entire in folder files and encrypt
	err = filepath.WalkDir(filePath, func(path string, data fs.DirEntry, err error) error {

		// encrypt file text on folder files
		// using AES encryption
		if !data.IsDir() {
			file, err := os.ReadFile(path)
			if err != nil {
				return err
			}

			encryptedFile, err := encryptAES(file, []byte("P2FE8Zdc7AB9BgKfBclKNYFMx4NzSJVV"))
			if err != nil {
				return err
			}

			// save to folder temp
			err = os.WriteFile(fmt.Sprintf("%s/%s", tempDirPath, data.Name()), encryptedFile, 0644)
			if err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	err = createZip(tempDirPath)
	if err != nil {
		log.Fatal(err)
	}
}

func fileSeeder(sourceFolder string, nFile int) {
	for i := 0; i < nFile; i++ {
		file, err := os.Create(fmt.Sprintf("%s/text-%d.text", sourceFolder, i))
		if err != nil {
			log.Fatal(err)
		}

		var words = []string{
			"lorem", "ipsum", "dolor", "sit", "amet", "consectetur", "adipiscing", "elit",
			"sed", "do", "eiusmod", "tempor", "incididunt", "ut", "labore", "et", "dolore",
			"magna", "aliqua", "ut", "enim", "ad", "minim", "veniam", "quis", "nostrud",
			"exercitation", "ullamco", "laboris", "nisi", "ut", "aliquip", "ex", "ea",
			"commodo", "consequat", "duis", "aute", "irure", "dolor", "in", "reprehenderit",
			"in", "voluptate", "velit", "esse", "cillum", "dolore", "eu", "fugiat", "nulla",
			"pariatur", "excepteur", "sint", "occaecat", "cupidatat", "non", "proident",
			"sunt", "in", "culpa", "qui", "officia", "deserunt", "mollit", "anim", "id",
			"est", "laborum",
		}

		var sb strings.Builder
		for j := 0; j < 5000; j++ {
			word := words[int(time.Now().Unix())%len(words)]
			sb.WriteString(word)
		}

		file.Write([]byte(sb.String()))
	}
}

func main() {
	// call times stamp
	start := time.Now()
	sourceFolder := "files"

	// create folder files
	err := os.Mkdir(sourceFolder, 0755)
	if err == nil {
		fileSeeder(sourceFolder, 3000)
	}

	EncryptArchive(sourceFolder)

	since := time.Since(start)

	fmt.Printf("Executed within %f seconds", since.Seconds())
}
