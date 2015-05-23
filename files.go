package main

import (
	"os"
	"path/filepath"
)

// Recursively outputs each file in the root directory
func walkFiles(root string) <-chan string {
	out := make(chan string)
	go func() {
		filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			abspath, _ := filepath.Abs(path)
			if _, elem := filepath.Split(abspath); elem != "" {
				// Skip various temporary or "hidden" files or directories.
				if elem[0] == '.' || elem[0] == '#' || elem[0] == '~' || elem[len(elem)-1] == '~' {
					if info.IsDir() {
						return filepath.SkipDir
					}
					return nil
				}
			}

			if info != nil && info.Mode()&os.ModeType == 0 {
				// out <- abspath
				out <- path
			}

			return nil
		})
		close(out)
	}()
	return out
}
