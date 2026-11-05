package main

import (
	"fmt"
	"io/fs"
	"log"
	"path/filepath"
)

func main() {
	root := ".."

	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			fmt.Printf("dir: %s\n", path)
		} else {
			fmt.Printf("file: %s\n", path)
		}

		return nil
	})
	if err != nil {
		log.Fatalf("walkdir err: %v", err)
	}
}
