package fileseeder

import (
	"fmt"
	"os"
	"strings"
	"time"
)

type File struct {
	*os.File
}

func createFile(nFile int, targetFolder string) <-chan File {
	chanFileOut := make(chan File)

	go func(out chan File) {
		defer close(chanFileOut)
		for i := 0; i < nFile; i++ {
			file, _ := os.Create(fmt.Sprintf("%s/text-%d.txt", targetFolder, i))
			out <- File{File: file}
		}
	}(chanFileOut)

	return chanFileOut
}

func writeFile(chanFileIn <-chan File) {
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

	for file := range chanFileIn {
		go func(file File) {
			var sb strings.Builder
			for j := 0; j < 5000; j++ {
				word := words[int(time.Now().Unix())%len(words)]
				sb.WriteString(word)
			}
			file.Write([]byte(sb.String()))
			defer file.Close()
		}(file)
	}
}

func ExecSeeder(nFile int, targetFolder string) {
	chanFile := createFile(nFile, targetFolder)
	writeFile(chanFile)
	writeFile(chanFile)
	writeFile(chanFile)
}
