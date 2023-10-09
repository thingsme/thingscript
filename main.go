package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/thingsme/thingscript/eval"
	"github.com/thingsme/thingscript/lexer"
	"github.com/thingsme/thingscript/object"
	"github.com/thingsme/thingscript/parser"
	"github.com/thingsme/thingscript/stdlib"
)

func main() {
	var verbose = false
	var content string

	flag.BoolVar(&verbose, "verbose", false, "verbose")
	flag.Parse()

	args := flag.Args()
	if len(args) == 1 {
		b, err := os.ReadFile(args[0])
		if err != nil {
			fmt.Println("File not found", err.Error())
			os.Exit(2)
		}
		content = string(b)
	} else if len(args) != 1 {
		fmt.Println("Usage: thingscript <flags> filename")
		os.Exit(1)
	} else {
		reader := bufio.NewReader(os.Stdin)
		buff := []string{}
		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				if err == io.EOF {
					break
				}
				fmt.Println("Input stream", err.Error())
				os.Exit(2)
			}
			buff = append(buff, line)
		}
		content = strings.Join(buff, "\n")
	}

	l := lexer.New(content)
	p := parser.New(l)
	program := p.ParseProgram()
	if len(p.Errors()) > 0 {
		for _, e := range p.Errors() {
			fmt.Println("ERR", e)
		}
		os.Exit(3)
	}
	env := object.NewEnvironment()
	env.RegisterPackages(stdlib.Packages()...)
	eval.Eval(program, env)
}
