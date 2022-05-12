package main

import (
	"bufio"
	"bytes"
	"fmt"
	"gollect/log"
	"gollect/utils"
	"gopkg.in/alecthomas/kingpin.v2"
	"os"
	"sync"
)

var (
	target           = kingpin.Flag(ArgTargetLongFlag, ArgTargetHelpMsg).Short(ArgTargetShortFlag).Default(".").ExistingDir()
	length           = kingpin.Flag(ArgLengthLongFlag, ArgLengthHelpMsg).Short(ArgLengthShortFlag).Default("0").Uint()
	isHead           = kingpin.Flag(ArgIsHeadLongFlag, ArgIsHeadHelpMsg).Default("true").Bool()
	isTail           = kingpin.Flag(ArgIsTailLongFlag, ArgIsTailHelpMsg).Bool()
	isIncludeExt     = kingpin.Flag(ArgIsIncludeExtensionLongFlag, ArgIsIncludeExtensionHelpMsg).Default("false").Bool()
	isStdout         = kingpin.Flag(ArgIsStdoutLongFlag, ArgIsStdoutHelpMsg).Short(ArgIsStdoutShortFlag).Default("false").Bool()
	isShowCopiedTree = kingpin.Flag(ArgIsShowFileTreeLongFlag, ArgIsShowFileTreeHelpMsg).Short(ArgIsShowFileTreeShortFlag).Default("false").Bool()
	regexpPattern    = kingpin.Flag(ArgRegexpLongFlag, ArgRegexPatternHelpMsg).Short(ArgRegexpShortFlag).Regexp()
	isPrintLog       = kingpin.Flag(ArgIsPrintLog, ArgIsPrintLogHelpMsg).Default("false").Bool()
)

func init() {
	kingpin.Parse()
}

func main() {
	logChan := make(chan string)
	logger := log.NewLogger("", logChan, &sync.WaitGroup{})
	logger.Add(1)
	go logger.Do()

	if *isPrintLog {
		logger.Close()

		f, err := logger.OpenLogFile()
		if err != nil {
			panic(err)
		}

		logReader := bufio.NewScanner(f)
		for logReader.Scan() {
			fmt.Println(logReader.Text())
		}
		fmt.Println(logger.GetLogPath())
	}

	if *length <= 0 && *regexpPattern == nil {
		return
	}

	fileNames, err := utils.GetTargetFile(*target)
	if err != nil {
		panic(err)
	}

	patternMap := make(map[string][]string)
	if *regexpPattern != nil {
		patternMap = utils.CreatePatternMapByRegexp(fileNames, regexpPattern)
	} else {
		patternMap = utils.CreatePatternMapByStringLength(fileNames, isHead, isTail, isIncludeExt, length)
	}

	outputPath, err := os.MkdirTemp(utils.MakePath("."), utils.AppendPrefixTimeStamp("_*"))
	if err != nil {
		panic(err)
	}

	wg := &sync.WaitGroup{}
	for mathPattern, matchFileNames := range patternMap {
		outputPath := utils.MakePath(outputPath, mathPattern)
		if err := utils.Mkdir(outputPath); err != nil {
			panic(err)
		}

		for _, matchFileName := range matchFileNames {
			wg.Add(1)
			go utils.DoFileCopy(wg, target, matchFileName, outputPath, isStdout)
		}
	}
	wg.Wait()
	buf := &bytes.Buffer{}
	utils.PrintTree(outputPath, patternMap, buf)
	logger.Send(buf.String())
	logger.Close()

	if *isShowCopiedTree {
		fmt.Println(buf.String())
	}
	logger.Wait()
}
