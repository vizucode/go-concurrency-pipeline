package securearchive

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
	"sync"
	"time"
)

// create an object file
type File struct {
	*os.File
}

/*
This function ReadFile demonstrates the implementation of a fan-out pattern using goroutines in Go.
The purpose is to read files from a source directory and pass File objects into a chanFileOut channel,
enabling streaming of file objects for further processing.
*/
func readFile(sourceDir string) <-chan File {
	chanFileOut := make(chan File)

	/*
		It will be used as a fan-out pattern.
		in loop `filepath.Walkdir` will append to channel `chanFileOut`.
		then will send by streaming through `chanFileOut`.
	*/
	go func(sourceDir string) {
		err := filepath.WalkDir(sourceDir, func(path string, data fs.DirEntry, err error) error {
			if !data.IsDir() {
				file, err := os.OpenFile(path, os.O_RDWR, 0644)
				if err != nil {
					log.Println(err)
				}

				chanFileOut <- File{
					File: file,
				}
			}
			return nil
		})

		if err != nil {
			log.Println(err)
		}

		close(chanFileOut)
	}(sourceDir)

	return chanFileOut
}

/*
The `encryptContent` function is designed to encrypt string data using an RSA 32-bit key.
It reads plain-text data from a channel containing a `File` struct, encrypts it, and
passes the encrypted data to another channel for consumption by other workers.
*/
func encryptContent(chanFile <-chan File) <-chan File {

	chanFileOut := make(chan File)

	/*
		This loop will be executed concurrently by multiple goroutines.
		If the `chanFile` is closed, the loop will ended.
	*/
	go func() {
		defer close(chanFileOut)
		for file := range chanFile {

			block, err := aes.NewCipher([]byte("P2FE8Zdc7AB9BgKfBclKNYFMx4NzSJVV"))
			if err != nil {
				log.Println(err)
			}

			gcm, err := cipher.NewGCM(block)
			if err != nil {
				log.Println(err)
			}

			nonce := make([]byte, gcm.NonceSize())
			if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
				log.Println(err)
			}

			contentFile, _ := io.ReadAll(file.File)
			encrypted := gcm.Seal(nonce, nonce, contentFile, nil)

			_, err = file.File.Write(encrypted)
			if err != nil {
				log.Println(err)
				continue
			}

			_, err = file.File.Seek(0, 0)
			if err != nil {
				log.Println(err)
			}

			chanFileOut <- File{
				File: file.File,
			}
		}
	}()

	return chanFileOut
}

/*
Will merge all channel from the `encryptContent` worker.
`MultiPlexerEncrypt` or “fan-in“ will be returning a single channel, that have been merged.

Why using multiplexer? Because I want to read from multiple worker at the same time.
For simple to consume in next function, just a single channel is enough.
*/
func multiPlexerEncrypt(chanFile ...<-chan File) <-chan File {
	chanFileOut := make(chan File)
	wg := new(sync.WaitGroup)
	wg.Add(len(chanFile))

	for _, channel := range chanFile {
		go func(chFile <-chan File) {
			defer wg.Done()
			for file := range chFile {
				chanFileOut <- file
			}
		}(channel)
	}

	go func() {
		wg.Wait()
		close(chanFileOut)
	}()

	return chanFileOut
}

/*
`archive` will execute the archive process.
it will create `output` folder, used to store zip file.
*/
func archive(chanFile <-chan File) {
	os.Mkdir("output", 0755)
	unix := time.Now().Unix()

	zipFile, err := os.Create(fmt.Sprintf("output/archive-%d.zip", unix))
	if err != nil {
		log.Println(err)
	}
	defer zipFile.Close()
	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	/*
		this loop doens't in goroutine, because it's not needed in this case.
		if these needed will be in goroutine.
	*/

	for file := range chanFile {
		stat, err := file.File.Stat()
		if err != nil {
			log.Println(err)
		}

		zipHeader, err := zip.FileInfoHeader(stat)
		if err != nil {
			log.Println(err)
		}

		writer, err := zipWriter.CreateHeader(zipHeader)
		if err != nil {
			log.Println(err)
		}

		_, err = file.File.Seek(0, 0)
		if err != nil {
			log.Println(err)
		}

		_, err = io.Copy(writer, file.File)
		if err != nil {
			log.Println(err)
		}

		file.Close()
	}
}

// will execute secrure archive process
func ExecSecureArchive(srcDir string) {
	chanFile := readFile(srcDir)

	chanEnc1 := encryptContent(chanFile)
	chanEnc2 := encryptContent(chanFile)
	chanEnc3 := encryptContent(chanFile)

	chanEncOut := multiPlexerEncrypt(chanEnc1, chanEnc2, chanEnc3)

	archive(chanEncOut)
}
