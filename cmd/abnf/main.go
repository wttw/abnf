package main

import (
	"github.com/wttw/abnf/internal/cmd"
	"os"
)

func main() {
	os.Exit(cmd.Do(os.Args[1:], os.Stdin, os.Stdout, os.Stderr))
}
