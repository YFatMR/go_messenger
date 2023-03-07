package test

import (
	"fmt"
)

type Program struct {
	SourceCode string
	Stdout     string
	Stderr     string
}

func NewHelloWorldProgram() Program {
	sourse := `
package main

import (
	"fmt"
)

func main() {
	fmt.Println("Hello world!")
}
`
	return Program{
		SourceCode: sourse,
		Stdout:     "Hello world!\n",
		Stderr:     "",
	}
}

func NewHugeHelloWorldWithComment(comment string) Program {
	sourse := fmt.Sprintf(`
package main

import (
	"fmt"
)

func main() {
	// %s
	fmt.Println("Hello world!")
}
`, comment)
	return Program{
		SourceCode: sourse,
		Stdout:     "Hello world!\n",
		Stderr:     "",
	}
}

func NewInfinityLoopProgram(comment string) Program {
	sourse := fmt.Sprintf(`
package main

import (
	"fmt"
)

func main() {
	// %s
	for {
		fmt.Println("Hello world!")
}
`, comment)
	return Program{
		SourceCode: sourse,
		Stdout:     "Hello world!\n",
		Stderr:     "",
	}
}
