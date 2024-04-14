package main

import (
	"fmt"
	"os"
)

const (
	TEMP_DIR = ".temp"
)

func main() {
	if _, ok := os.Stat(TEMP_DIR); ok != nil {
		os.Mkdir(TEMP_DIR, os.ModePerm)
	}

	i := 0
	for {
		filename := fmt.Sprintf("%s/file_%d", TEMP_DIR, i)
		file, err := os.Create(filename)
		if err != nil {
			panic(err)
		}
		defer file.Close()
		file.WriteString("Hello, World!")
		i++
	}
}
