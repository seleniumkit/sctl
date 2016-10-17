package main

import (
	"flag"
	"github.com/seleniumkit/sctl/cmd"
)

func init() {
	flag.Parse()
}

func main() {
	cmd.Execute()
}