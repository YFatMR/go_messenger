package programresult

type Entity struct {
	stdout string
	stderr string
}

func New(stdout string, stderr string) *Entity {
	return &Entity{
		stdout: stdout,
		stderr: stderr,
	}
}

func (e *Entity) GetStdOut() string {
	return e.stdout
}

func (e *Entity) GetStdErr() string {
	return e.stderr
}
