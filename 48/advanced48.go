package main

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
)

func main() {
	err := atomicWriteFile("./mylog.txt", []byte("Important info!"), 0755)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Done")
}

func atomicWriteFile(filename string, data []byte, perm fs.FileMode) error {
	dir := filepath.Dir(filename)
	tmpFile, err := os.CreateTemp(dir, filepath.Base(filename)+".*")
	if err != nil {
		return fmt.Errorf("create temp file: %w", err)
	}

	tmpName := tmpFile.Name()

	defer os.Remove(tmpName)
	defer tmpFile.Close()

	if _, err := tmpFile.Write(data); err != nil {
		return fmt.Errorf("write: %w", err)
	}

	if err := tmpFile.Sync(); err != nil {
		return fmt.Errorf("sync: %w", err)
	}

	if err := tmpFile.Chmod(perm); err != nil {
		return fmt.Errorf("chmod: %w", err)
	}

	if err := tmpFile.Close(); err != nil {
		return fmt.Errorf("close: %w", err)
	}

	if err := os.Rename(tmpName, filename); err != nil {
		return fmt.Errorf("close: %w", err)
	}

	return nil
}
