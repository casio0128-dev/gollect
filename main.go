package main

import (
	"gollect/utils"
	"gopkg.in/alecthomas/kingpin.v2"
	"os"
	"sync"
	"time"
)

var (
	target        = kingpin.Flag(ArgTargetLongFlag, ArgTargetHelpMsg).Short(ArgTargetShortFlag).Default(".").ExistingDir()
	length        = kingpin.Flag(ArgLengthLongFlag, ArgLengthHelpMsg).Short(ArgLengthShortFlag).Default("0").Uint()
	isHead        = kingpin.Flag(ArgIsHeadLongFlag, ArgIsHeadHelpMsg).Default("true").Bool()
	isTail        = kingpin.Flag(ArgIsTailLongFlag, ArgIsTailHelpMsg).Bool()
	isIncludeExt  = kingpin.Flag(ArgIsIncludeExtensionLongFlag, ArgIsIncludeExtensionHelpMsg).Default("false").Bool()
	isStdout      = kingpin.Flag(ArgIsStdoutLongFlag, ArgIsStdoutHelpMsg).Short(ArgIsStdoutShortFlag).Default("false").Bool()
	regexpPattern = kingpin.Flag(ArgRegexpLongFlag, ArgRegexPatternHelpMsg).Short(ArgRegexpShortFlag).Regexp()
)

func init() {
	kingpin.Parse()
}

func main() {
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

	outputPath, err := os.MkdirTemp(utils.MakePath("."), time.Now().Format("20060102150405_*"))
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
}
