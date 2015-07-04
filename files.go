package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func isVolumeRoot(path string) bool {
	return os.IsPathSeparator(path[len(path)-1])
}

func isGitRoot(path string) bool {
	gitpath := filepath.Join(path, ".git")

	fi, err := os.Stat(gitpath)
	if err != nil {
		return false
	}

	return fi.IsDir()
}

func findGitRoot(path string) (bool, string) {
	path = filepath.Clean(path)

	for isVolumeRoot(path) == false {
		if isGitRoot(path) {
			return true, path
		} else {
			path = filepath.Dir(path)
		}
	}
	return false, ""
}

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
