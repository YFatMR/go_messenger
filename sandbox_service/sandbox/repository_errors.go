package sandbox

import "errors"

var (
	ErrProgramCreation           = errors.New("can't create program")
	ErrUpdateProgramRunnerOutput = errors.New("can't update program runner output")
	ErrUpdateProgramSource       = errors.New("can't update program source")
	ErrGetProgramByID            = errors.New("can't get program by id")
	ErrProgramNotFount           = errors.New("program with provided id not exist")
	ErrUpdateLinterOutput        = errors.New("can not update linter output")
)
