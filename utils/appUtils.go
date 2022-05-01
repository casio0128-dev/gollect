package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
)

func CreatePatternMapByRegexp(fileNames []string, regexpPattern **regexp.Regexp) map[string][]string {
	patternMap := make(map[string][]string)
	for _, fileName := range fileNames {
		var key string
		if *regexpPattern != nil {
			if (*regexpPattern).Match([]byte(fileName)) {
				key = (*regexpPattern).FindString(fileName)
			} else {
				continue
			}
		}
		patternMap[key] = append(patternMap[key], fileName)
	}
	return patternMap
}

func CreatePatternMapByStringLength(fileNames []string, isHead, isTail, isIncludeExt *bool, length *uint) map[string][]string {
	patternMap := make(map[string][]string)
	for _, fileName := range fileNames {
		var key string
		if *isHead {
			key = fileName[:*length]
		}
		if *isTail {
			tailStartIndex := len(fileName) - int(*length)
			if *isIncludeExt {
				key = fileName[tailStartIndex:]
			} else {
				ext := filepath.Ext(fileName)
				if !strings.EqualFold(ext, "") {
					fn := strings.Split(fileName, ext)[0]
					key = fn[tailStartIndex-len(ext):]
				} else {
					key = fileName[tailStartIndex:]
				}
			}
		}
		patternMap[key] = append(patternMap[key], fileName)
	}
	return patternMap
}

func DoFileCopy(wg *sync.WaitGroup, target *string, matchFileName, outputPath string, isStdout *bool) {
	defer func() {
		switch recover().(type) {
		default:
			wg.Done()
		}
	}()

	srcPath := MakePath(*target, matchFileName)
	dstPath := MakePath(outputPath, matchFileName)

	if *isStdout {
		fmt.Println(filepath.Clean(srcPath), " -> ", filepath.Clean(dstPath))
	}
	if _, err := fileCopy(srcPath, dstPath); err != nil {
		panic(err)
	}
}

func GetTargetFile(target string) ([]string, error) {
	return getDirNames(strings.Join([]string{target}, string(os.PathSeparator)), func(de os.DirEntry) bool {
		return de.IsDir()
	})
}