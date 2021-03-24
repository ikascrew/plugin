package terminal_test

import (
	"fmt"

	"github.com/ikascrew/plugin/terminal"
)

func TestTerminal() {
	v, err := terminal.New()
	if err != nil {
		panic(err)
	}

	fmt.Println(v)
}
