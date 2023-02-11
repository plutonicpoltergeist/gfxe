package main

var (
	errCreatePatternFile    = "Cannot create pattern file: %w"
	errGetPattern           = "Error getting patterns: %s"
	errGetPatternDir        = "Error while get pattern directory: %w"
	errNoPattern            = "No such pattern for '%s'"
	errOperatorCmdNotFound  = "Operator '%s' command could not be found"
	errPatternDirNotFound   = "Pattern directory not found in either %s or %s"
	errPatternFileMalformed = "Pattern file '%s' is malformed: %s"
	errPatternFileNoPattern = "Pattern file '%s' contains no pattern(s)"
	errSavePattern          = "Error saving pattern: %s"
	errWritePatternFile     = "Cannot write pattern file: %w"
)

const (
	errOpenUserPatternDir = "Error opening user's pattern directory"
	errNoPatternInput     = "No pattern input provided"
	errPatternFlagsEmpty  = "Pattern flags cannot be empty"
	errPatternNameEmpty   = "Pattern name cannot be empty"
	errPatternValueEmpty  = "Pattern value cannot be empty"
)
