package utils

import (
	"io"
	"os"
	"strings"
)

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
	for _, dir := range dirs {
		if skipCondition(dir) {
			continue
		}
		result = append(result, dir.Name())
	}
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
