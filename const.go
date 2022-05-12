package main

const (
	ArgTargetLongFlag  = "target"
	ArgTargetShortFlag = 't'
	ArgTargetHelpMsg   = "Specify the directory containing the files to be classification."

	ArgLengthLongFlag  = "length"
	ArgLengthShortFlag = 'l'
	ArgLengthHelpMsg   = "Specify the number of characters that will be the classification condition."

	ArgRegexpLongFlag      = "regexp"
	ArgRegexpShortFlag     = 'r'
	ArgRegexPatternHelpMsg = "Specifies the regular expression pattern to be used for classification."

	ArgIsHeadLongFlag = "head"
	ArgIsHeadHelpMsg  = "Count the number of characters to be classification from the beginning."

	ArgIsTailLongFlag = "tail"
	ArgIsTailHelpMsg  = "Count the number of characters to be classification from the end."

	ArgIsIncludeExtensionLongFlag = "include-ext"
	ArgIsIncludeExtensionHelpMsg  = "The extension is included in the string to be classified."

	ArgIsStdoutLongFlag  = "print"
	ArgIsStdoutShortFlag = 'p'
	ArgIsStdoutHelpMsg   = "Displays the source and destination of the copy on standard output."

	ArgIsShowFileTreeLongFlag  = "show-tree"
	ArgIsShowFileTreeShortFlag = 's'
	ArgIsShowFileTreeHelpMsg   = "Display the files to be copied in a tree view."

	ArgIsInteractiveLongFlag  = "interactive"
	ArgIsInteractiveShortFlag = 'i'
	ArgIsInteractiveHelpMsg   = "Confirm whether to create it."

	ArgIsPrintLog        = "print-log"
	ArgIsPrintLogHelpMsg = "Prints the contents of the log."
)
