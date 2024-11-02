package repl

import (
	"bufio"
	"fmt"
	"io"
	"log"

	"github.com/vincentlabelle/monkey/evaluator"
	"github.com/vincentlabelle/monkey/lexer"
	"github.com/vincentlabelle/monkey/object"
	"github.com/vincentlabelle/monkey/parser"
)

const PROMPT = ">> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	env := object.NewEnvironment()

	for {
		fmt.Fprint(out, PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		lex := lexer.New(line)
		p := parser.New(lex)
		program := p.ParseProgram()
		evaluated := evaluator.Eval(program, env)

		if evaluated != nil {
			_, err := io.WriteString(out, evaluated.Inspect()+"\n")
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}
