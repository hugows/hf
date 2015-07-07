package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
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
				if elem[0] == '.' ||
					elem[0] == '$' ||
					elem[0] == '#' ||
					elem[0] == '~' ||
					elem[len(elem)-1] == '~' ||
					strings.HasSuffix(elem, ".app") {
					if info.IsDir() {
						return filepath.SkipDir
					}
					return nil
				}
			}
			if info != nil && info.Mode()&os.ModeType == 0 {
				if strings.HasPrefix(path, root) {
					path = path[len(root)+1:]
				}
				out <- path
			}

			return nil
		}) // walk fn

		close(out)

	}()

	return out
}

func walkFilesFake(count int) <-chan string {
	out := make(chan string, 1000)

	go func() {
		for i := 0; i < count; i++ {
			out <- fmt.Sprintf("brasil%d", i)
		}

		close(out)

	}()

	return out
}
