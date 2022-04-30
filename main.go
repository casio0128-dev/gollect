package main

import (
	"fmt"
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
	target := kingpin.Flag(ArgTargetLongFlag, ArgTargetHelpMsg).Short(ArgTargetShortFlag).Default(".").ExistingDir()
	length := kingpin.Flag(ArgLengthLongFlag, ArgLengthHelpMsg).Short(ArgLengthShortFlag).Default("0").Uint()
	isHead := kingpin.Flag(ArgIsHeadLongFlag, ArgIsHeadHelpMsg).Default("true").Bool()
	isTail := kingpin.Flag(ArgIsTailLongFlag, ArgIsTailHelpMsg).Bool()
	isIncludeExt := kingpin.Flag(ArgIsIncludeExtensionLongFlag, ArgIsIncludeExtensionHelpMsg).Default("false").Bool()
	isStdout := kingpin.Flag(ArgIsStdoutLongFlag, ArgIsStdoutHelpMsg).Short(ArgIsStdoutShortFlag).Default("false").Bool()
	regexpPattern := kingpin.Flag(ArgRegexpLongFlag, ArgRegexPatternHelpMsg).Short(ArgRegexpShortFlag).Regexp()
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
	for mathPattern, matchFileNames := range patternMap {
		wg.Add(1)
		outputPath := makePath(outputPath, mathPattern)
		if err := mkdir(outputPath); err != nil {
			panic(err)
		}

		for _, matchFileName := range matchFileNames {
			wg.Add(1)
			go doFileCopy(wg, target, matchFileName, outputPath, isStdout)
		}
		wg.Done()
	}
	wg.Wait()
}

func doFileCopy(wg *sync.WaitGroup, target *string, matchFileName, outputPath string, isStdout *bool) {
	defer func() {
		switch recover().(type) {
		default:
			wg.Done()
		}
	}()

	srcPath := makePath(*target, matchFileName)
	dstPath := makePath(outputPath, matchFileName)

	if *isStdout {
		fmt.Println(filepath.Clean(srcPath), " -> ", filepath.Clean(dstPath))
	}
	if _, err := fileCopy(srcPath, dstPath); err != nil {
		panic(err)
	}
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
	return os.MkdirAll(path, 0777)
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
