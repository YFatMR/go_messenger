package entity

type Program struct {
	ID               ProgramID
	Source           ProgramSource
	CodeRunnerOutput ProgramOutput
	LinterOutput     ProgramOutput
}
