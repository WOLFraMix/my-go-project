package main

import (
	"fmt"
	"io/fs"
	"log"
	"path/filepath"
)

func main() {
	root := "."

	size, err := dirSize(root)
	if err != nil {
		log.Fatalf("Ошибка при подсчёте размера: %v", err)
	}

	fmt.Printf("Общий размер: %d байт\n", size)
}

func dirSize(path string) (int64, error) {
	var size int64

	err := filepath.WalkDir(path, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			log.Printf("access denied %q: %v\n", path, err)
			return err
		}
		if d.IsDir() {
			return nil
		}
		info, err := d.Info()
		if err != nil {
			return err
		}
		size += info.Size()

		return nil
	})
	if err != nil {
		return 0, fmt.Errorf("walk dir err: %w", err)
	}

	return size, nil
}
