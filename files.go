package main

import (
	"os"
	"path/filepath"
)

// Recursively outputs each file in the root directory
func walkFiles(root string) <-chan string {
	out := make(chan string, 1000)

	go func() {
		filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
			// just skip and continue when folders fail
			if err != nil {
				return nil
			}

			abspath, _ := filepath.Abs(path)
			abspathclean := filepath.Clean(abspath)
			if _, elem := filepath.Split(abspathclean); elem != "" {
				// Skip various temporary or "hidden" files or directories.
				if elem[0] == '.' || elem[0] == '$' || elem[0] == '#' || elem[0] == '~' || elem[len(elem)-1] == '~' {
					if info.IsDir() {
						return filepath.SkipDir
					}
					return nil
				}
			}

			if info != nil && info.Mode()&os.ModeType == 0 {
				// out <- abspathclean
				out <- path
			}

			return nil
		})

		close(out)

	}()

	return out
}
