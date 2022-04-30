package main

import (
	"gopkg.in/alecthomas/kingpin.v2"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"
)

func main() {
	target := kingpin.Flag("target", "target directory path").Short('t').Default(".").ExistingDir()
	length := kingpin.Flag("length", "compare to filename length").Short('l').Default("0").Uint()
	isHead := kingpin.Flag("head", "compare to head").Default("true").Bool()
	isTail := kingpin.Flag("tail", "compare to tail").Bool()
	isIncludeExt := kingpin.Flag("include-ext", "include extension name").Default("false").Bool()
	regexpPattern := kingpin.Flag("regexp", "regex pattern").Short('r').Regexp()
	kingpin.Parse()

	if *length <= 0 && *regexpPattern == nil {
		return
	}

	fileNames, err := getFileNames(*target)
	if err != nil {
		panic(err)
	}

	patternMap := make(map[string][]string)
	if *regexpPattern != nil {
		patternMap = createPatternMapByRegexp(fileNames, regexpPattern)
	} else {
		patternMap = createPatternMapByStringLength(fileNames, isHead, isTail, isIncludeExt, length)
	}

	outputPath, err := os.MkdirTemp(makePath("."), time.Now().Format("20060102150405_*"))
	if err != nil {
		panic(err)
	}

	wg := &sync.WaitGroup{}
	for key, values := range patternMap {
		wg.Add(1)

		go func(k string, v []string) {
			defer func() {
				switch recover().(type) {
				default:
					wg.Done()
				}
			}()

			p := makePath(outputPath, k)
			if err := mkdir(p); err != nil {
				panic(err)
			}

			wg2 := &sync.WaitGroup{}
			for _, val := range v {
				wg2.Add(1)
				go func(val2 string) {
					defer func() {
						switch recover().(type) {
						default:
							wg2.Done()
						}
					}()
					srcPath := makePath(*target, val2)
					dstPath := makePath(outputPath, k, val2)
					if _, err := fileCopy(srcPath, dstPath); err != nil {
						panic(err)
					}
				}(val)
			}
			wg2.Wait()
		}(key, values)
	}

	wg.Wait()
}

func createPatternMapByRegexp(fileNames []string, regexpPattern **regexp.Regexp) map[string][]string {
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

func createPatternMapByStringLength(fileNames []string, isHead, isTail, isIncludeExt *bool, length *uint) map[string][]string {
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

func mkdir(path string) error {
	return os.MkdirAll(path, 0766)
}

func makePath(paths ...string) string {
	return strings.Join(paths, string(os.PathSeparator))
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

func getFileNames(target string) ([]string, error) {
	return getDirNames(strings.Join([]string{target}, string(os.PathSeparator)), func(de os.DirEntry) bool {
		return de.IsDir()
	})
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
