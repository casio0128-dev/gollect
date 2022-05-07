package utils

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

func PrintTree(target string, tree map[string][]string, out io.Writer) {
	var targetIndent int
	currentTarget := strings.Split(filepath.Clean(target), string(os.PathSeparator))
	for i, t := range currentTarget {
		if i == 0 {
			fmt.Fprintf(out, "%s\n", t)
		} else {
			fmt.Fprintf(out, "%s└%s\n", createBlank(uint(targetIndent), " "), t)
		}
		targetIndent += i + 1
	}

	var currentIndent = uint(targetIndent + 1)
	var dirCount = len(tree)
	var dirIndex int
	for dirName, fileNames := range tree {
		dirIndex++
		var lineStr string
		if dirIndex == dirCount {
			lineStr += fmt.Sprintf("%s└%s\n", createBlank(currentIndent, " "), dirName)
		} else {
			lineStr += fmt.Sprintf("%s├%s\n", createBlank(currentIndent, " "), dirName)
		}

		for fileIndex, fileName := range fileNames {
			if dirIndex != dirCount {
				lineStr += fmt.Sprintf("%s│", createBlank(currentIndent, " "))
			}

			if fileIndex < len(fileNames)-1 {
				lineStr += fmt.Sprintf("%s├%s\n", createBlank(currentIndent+1, " "), fileName)
			} else {
				lineStr += fmt.Sprintf("%s└%s", createBlank(currentIndent+1, " "), fileName)
			}
		}
		if _, err := fmt.Fprintln(out, lineStr); err != nil {
			panic(err)
		}
	}
}

func createBlank(count uint, blankStr string) string {
	if count <= 0 {
		return ""
	}
	var result string
	for i := 0; i < int(count); i++ {
		result += blankStr
	}
	return result
}

func Mkdir(path string) error {
	return os.MkdirAll(path, 0777)
}

func MakePath(paths ...string) string {
	return strings.Join(paths, string(os.PathSeparator))
}

func getDirNames(path string, skipCondition func(entry os.DirEntry) bool) ([]string, error) {
	dirs, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	var result []string
	wg := &sync.WaitGroup{}
	for _, dir := range dirs {
		wg.Add(1)
		go func(file os.DirEntry) {
			defer func() {
				switch recover().(type) {
				default:
					wg.Done()
				}
			}()

			if skipCondition(file) {
				return
			}
			result = append(result, file.Name())
		}(dir)
	}
	wg.Wait()
	return result, nil
}

func fileCopy(src, dst string) (int64, error) {
	srcFile, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer func() {
		if err := srcFile.Close(); err != nil {
			panic(err)
		}
	}()

	dstFile, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer func() {
		if err := dstFile.Close(); err != nil {
			panic(err)
		}
	}()

	return io.Copy(dstFile, srcFile)
}
